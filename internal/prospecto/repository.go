package prospecto

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository persiste y consulta prospectos en PostgreSQL.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository crea un repositorio de prospectos.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// CursoIDs devuelve el conjunto de ids de cursos existentes.
func (r *Repository) CursoIDs(ctx context.Context) (map[int]struct{}, error) {
	rows, err := r.pool.Query(ctx, `SELECT id FROM cursos`)
	if err != nil {
		return nil, fmt.Errorf("consultar cursos: %w", err)
	}
	ids, err := pgx.CollectRows(rows, pgx.RowTo[int])
	if err != nil {
		return nil, fmt.Errorf("leer cursos: %w", err)
	}

	set := make(map[int]struct{}, len(ids))
	for _, id := range ids {
		set[id] = struct{}{}
	}
	return set, nil
}

// InsertBatch inserta los prospectos en una sola transacción.
// Los ids que ya existen se ignoran (ON CONFLICT DO NOTHING), por lo que reenviar es idempotente.
func (r *Repository) InsertBatch(ctx context.Context, ps []ProspectoInput) error {
	if len(ps) == 0 {
		return nil
	}

	const q = `
		INSERT INTO prospectos (id, nombre, email, telefono, nivel_educativo, curso_id,
		                        modalidad, canal_contacto, observaciones, capturado_en)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) DO NOTHING`

	return pgx.BeginFunc(ctx, r.pool, func(tx pgx.Tx) error {
		batch := &pgx.Batch{}
		for _, p := range ps {
			batch.Queue(q, p.ID, p.Nombre, p.Email, p.Telefono, p.NivelEducativo, p.CursoID,
				p.Modalidad, p.CanalContacto, p.Observaciones, p.CapturadoEn.UTC())
		}
		if err := tx.SendBatch(ctx, batch).Close(); err != nil {
			return fmt.Errorf("insertar prospectos: %w", err)
		}
		return nil
	})
}

// List devuelve los prospectos consolidados, del más reciente al más antiguo.
func (r *Repository) List(ctx context.Context) ([]Prospecto, error) {
	const q = `
		SELECT p.id::text, p.nombre, p.email, p.telefono, p.nivel_educativo, p.curso_id,
		       c.nombre, a.nombre, p.modalidad, p.canal_contacto, p.observaciones,
		       p.capturado_en, p.recibido_en
		FROM prospectos p
		JOIN cursos c ON c.id = p.curso_id
		JOIN areas a ON a.id = c.area_id
		ORDER BY p.capturado_en DESC`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return nil, fmt.Errorf("consultar prospectos: %w", err)
	}
	ps, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (Prospecto, error) {
		var p Prospecto
		err := row.Scan(&p.ID, &p.Nombre, &p.Email, &p.Telefono, &p.NivelEducativo, &p.CursoID,
			&p.CursoNombre, &p.AreaNombre, &p.Modalidad, &p.CanalContacto, &p.Observaciones,
			&p.CapturadoEn, &p.RecibidoEn)
		p.CapturadoEn = p.CapturadoEn.UTC()
		p.RecibidoEn = p.RecibidoEn.UTC()
		return p, err
	})
	if err != nil {
		return nil, fmt.Errorf("leer prospectos: %w", err)
	}
	return ps, nil
}
