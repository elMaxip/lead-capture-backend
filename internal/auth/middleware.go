package auth

import (
	"errors"
	"log/slog"
	"net/http"
	"slices"

	"github.com/therionlabs/lead-capture-backend/internal/httpx"
)

// RequireRol exige una sesión válida cuyo usuario tenga alguno de los roles indicados.
// Responde 401 si no hay sesión o no es válida, y 403 si el rol no tiene acceso.
func RequireRol(repo Store, roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := bearerToken(r)
			if token == "" {
				httpx.WriteError(w, http.StatusUnauthorized, "debes iniciar sesión")
				return
			}

			u, err := repo.UsuarioPorSesion(r.Context(), hashToken(token))
			if errors.Is(err, ErrNoEncontrado) {
				httpx.WriteError(w, http.StatusUnauthorized, "la sesión no es válida")
				return
			}
			if err != nil {
				slog.ErrorContext(r.Context(), "validar sesión", "error", err)
				httpx.WriteError(w, http.StatusInternalServerError, "no se pudo validar la sesión")
				return
			}

			if !slices.Contains(roles, u.Rol) {
				httpx.WriteError(w, http.StatusForbidden, "tu perfil no tiene acceso a este recurso")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
