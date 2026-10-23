package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/therionlabs/lead-capture-backend/internal/httpx"
)

// Store es lo que el handler y el middleware necesitan del repositorio.
type Store interface {
	PorUsername(ctx context.Context, username string) (Usuario, string, error)
	CrearSesion(ctx context.Context, tokenHash string, usuarioID int) error
	UsuarioPorSesion(ctx context.Context, tokenHash string) (Usuario, error)
	EliminarSesion(ctx context.Context, tokenHash string) error
}

// Handler expone el inicio y cierre de sesión.
type Handler struct {
	repo Store
}

// NewHandler crea el handler de autenticación.
func NewHandler(repo Store) *Handler {
	return &Handler{repo: repo}
}

// Login atiende POST /api/auth/login.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req LoginRequest
	if err := httpx.ReadJSON(w, r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	req.Username = strings.ToLower(strings.TrimSpace(req.Username))
	if req.Username == "" || req.Password == "" {
		httpx.WriteError(w, http.StatusBadRequest, "usuario y contraseña son obligatorios")
		return
	}

	u, hash, err := h.repo.PorUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, ErrNoEncontrado) {
		slog.ErrorContext(ctx, "buscar usuario", "error", err)
		httpx.WriteError(w, http.StatusInternalServerError, "no se pudo iniciar sesión")
		return
	}
	if err != nil || bcrypt.CompareHashAndPassword([]byte(hash), []byte(req.Password)) != nil {
		httpx.WriteError(w, http.StatusUnauthorized, "usuario o contraseña incorrectos")
		return
	}

	token := rand.Text()
	if err := h.repo.CrearSesion(ctx, hashToken(token), u.ID); err != nil {
		slog.ErrorContext(ctx, "crear sesión", "error", err)
		httpx.WriteError(w, http.StatusInternalServerError, "no se pudo iniciar sesión")
		return
	}

	slog.InfoContext(ctx, "inicio de sesión", "rol", u.Rol)
	httpx.WriteJSON(w, http.StatusOK, LoginResponse{Token: token, Usuario: u})
}

// Logout atiende POST /api/auth/logout.
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	if err := h.repo.EliminarSesion(r.Context(), hashToken(bearerToken(r))); err != nil {
		slog.ErrorContext(r.Context(), "cerrar sesión", "error", err)
		httpx.WriteError(w, http.StatusInternalServerError, "no se pudo cerrar la sesión")
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// hashToken calcula el SHA-256 del token: en la base nunca se guarda el token en claro.
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// bearerToken extrae el token de la cabecera "Authorization: Bearer <token>".
func bearerToken(r *http.Request) string {
	token, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok {
		return ""
	}
	return strings.TrimSpace(token)
}
