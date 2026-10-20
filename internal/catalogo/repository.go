package catalogo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository lee el catálogo de áreas y cursos desde PostgreSQL.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository crea un repositorio de catálogo.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Get devuelve todas las áreas con sus cursos anidados, ordenados por id.
func (r *Repository) Get(ctx context.Context) (Catalogo, error) {
	const q = `
		SELECT a.id, a.nombre, c.id, c.nombre
		FROM areas a
		LEFT JOIN cursos c ON c.area_id = a.id
		ORDER BY a.id, c.id`

	rows, err := r.pool.Query(ctx, q)
	if err != nil {
		return Catalogo{}, fmt.Errorf("consultar catálogo: %w", err)
	}
	defer rows.Close()

	cat := Catalogo{Areas: []Area{}}
	for rows.Next() {
		var (
			areaID      int
			areaNombre  string
			cursoID     *int
			cursoNombre *string
		)
		if err := rows.Scan(&areaID, &areaNombre, &cursoID, &cursoNombre); err != nil {
			return Catalogo{}, fmt.Errorf("leer fila de catálogo: %w", err)
		}
		if n := len(cat.Areas); n == 0 || cat.Areas[n-1].ID != areaID {
			cat.Areas = append(cat.Areas, Area{ID: areaID, Nombre: areaNombre, Cursos: []Curso{}})
		}
		if cursoID != nil {
			last := &cat.Areas[len(cat.Areas)-1]
			last.Cursos = append(last.Cursos, Curso{ID: *cursoID, Nombre: *cursoNombre})
		}
	}
	if err := rows.Err(); err != nil {
		return Catalogo{}, fmt.Errorf("recorrer catálogo: %w", err)
	}
	return cat, nil
}
