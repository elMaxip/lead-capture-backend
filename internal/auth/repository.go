package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNoEncontrado indica que no existe el usuario o la sesión buscada.
var ErrNoEncontrado = errors.New("no encontrado")

// Repository persiste usuarios y sesiones en PostgreSQL.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository crea un repositorio de usuarios y sesiones.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// PorUsername devuelve el usuario y el hash bcrypt de su contraseña.
func (r *Repository) PorUsername(ctx context.Context, username string) (Usuario, string, error) {
	const q = `SELECT id, username, nombre, rol, password_hash FROM usuarios WHERE username = $1`

	var u Usuario
	var hash string
	err := r.pool.QueryRow(ctx, q, username).Scan(&u.ID, &u.Username, &u.Nombre, &u.Rol, &hash)
	if errors.Is(err, pgx.ErrNoRows) {
		return Usuario{}, "", ErrNoEncontrado
	}
	if err != nil {
		return Usuario{}, "", fmt.Errorf("consultar usuario: %w", err)
	}
	return u, hash, nil
}

// CrearSesion registra una sesión nueva para el usuario.
func (r *Repository) CrearSesion(ctx context.Context, tokenHash string, usuarioID int) error {
	const q = `INSERT INTO sesiones (token_hash, usuario_id) VALUES ($1, $2)`
	if _, err := r.pool.Exec(ctx, q, tokenHash, usuarioID); err != nil {
		return fmt.Errorf("crear sesión: %w", err)
	}
	return nil
}

// UsuarioPorSesion devuelve el dueño de la sesión indicada.
func (r *Repository) UsuarioPorSesion(ctx context.Context, tokenHash string) (Usuario, error) {
	const q = `
		SELECT u.id, u.username, u.nombre, u.rol
		FROM sesiones s
		JOIN usuarios u ON u.id = s.usuario_id
		WHERE s.token_hash = $1`

	var u Usuario
	err := r.pool.QueryRow(ctx, q, tokenHash).Scan(&u.ID, &u.Username, &u.Nombre, &u.Rol)
	if errors.Is(err, pgx.ErrNoRows) {
		return Usuario{}, ErrNoEncontrado
	}
	if err != nil {
		return Usuario{}, fmt.Errorf("consultar sesión: %w", err)
	}
	return u, nil
}

// EliminarSesion borra la sesión; si no existe no hace nada.
func (r *Repository) EliminarSesion(ctx context.Context, tokenHash string) error {
	if _, err := r.pool.Exec(ctx, `DELETE FROM sesiones WHERE token_hash = $1`, tokenHash); err != nil {
		return fmt.Errorf("eliminar sesión: %w", err)
	}
	return nil
}
