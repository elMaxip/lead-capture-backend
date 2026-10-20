package catalogo

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/therionlabs/lead-capture-backend/internal/httpx"
)

// Store es lo que el handler necesita del repositorio.
type Store interface {
	Get(ctx context.Context) (Catalogo, error)
}

// Handler expone el catálogo por HTTP.
type Handler struct {
	repo Store
}

// NewHandler crea el handler del catálogo.
func NewHandler(repo Store) *Handler {
	return &Handler{repo: repo}
}

// Get atiende GET /api/catalogo.
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	cat, err := h.repo.Get(r.Context())
	if err != nil {
		slog.ErrorContext(r.Context(), "obtener catálogo", "error", err)
		httpx.WriteError(w, http.StatusInternalServerError, "no se pudo obtener el catálogo")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, cat)
}
