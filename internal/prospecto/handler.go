package prospecto

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/therionlabs/lead-capture-backend/internal/httpx"
)

// Store es lo que el handler necesita del repositorio.
type Store interface {
	CursoIDs(ctx context.Context) (map[int]struct{}, error)
	InsertBatch(ctx context.Context, ps []ProspectoInput) error
	List(ctx context.Context) ([]Prospecto, error)
}

// Handler expone la sincronización y el listado de prospectos.
type Handler struct {
	repo      Store
	validator *Validator
}

// NewHandler crea el handler de prospectos.
func NewHandler(repo Store, validator *Validator) *Handler {
	return &Handler{repo: repo, validator: validator}
}

// Sync atiende POST /api/prospectos/sync.
func (h *Handler) Sync(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req SyncRequest
	if err := httpx.ReadJSON(w, r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	if len(req.Prospectos) > MaxPorLote {
		httpx.WriteError(w, http.StatusBadRequest,
			fmt.Sprintf("el lote supera el máximo de %d prospectos", MaxPorLote))
		return
	}

	cursos, err := h.repo.CursoIDs(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "cargar cursos", "error", err)
		httpx.WriteError(w, http.StatusInternalServerError, "no se pudo sincronizar")
		return
	}

	resp := SyncResponse{Sincronizados: []string{}, Rechazados: []Rechazado{}}
	validos := make([]ProspectoInput, 0, len(req.Prospectos))
	vistos := make(map[string]struct{}, len(req.Prospectos))
	for _, p := range req.Prospectos {
		Normalize(&p)
		if errs := h.validator.Validate(p, cursos); len(errs) > 0 {
			resp.Rechazados = append(resp.Rechazados, Rechazado{ID: p.ID, Errores: errs})
			continue
		}
		validos = append(validos, p)
		if _, dup := vistos[p.ID]; !dup {
			vistos[p.ID] = struct{}{}
			resp.Sincronizados = append(resp.Sincronizados, p.ID)
		}
	}

	if err := h.repo.InsertBatch(ctx, validos); err != nil {
		slog.ErrorContext(ctx, "guardar lote de prospectos", "error", err, "cantidad", len(validos))
		httpx.WriteError(w, http.StatusInternalServerError, "no se pudo sincronizar")
		return
	}

	// Solo se registran cantidades, nunca datos personales.
	slog.InfoContext(ctx, "lote sincronizado",
		"recibidos", len(req.Prospectos),
		"sincronizados", len(resp.Sincronizados),
		"rechazados", len(resp.Rechazados))
	httpx.WriteJSON(w, http.StatusOK, resp)
}

// List atiende GET /api/prospectos.
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	ps, err := h.repo.List(r.Context())
	if err != nil {
		slog.ErrorContext(r.Context(), "listar prospectos", "error", err)
		httpx.WriteError(w, http.StatusInternalServerError, "no se pudieron obtener los prospectos")
		return
	}
	httpx.WriteJSON(w, http.StatusOK, ps)
}
