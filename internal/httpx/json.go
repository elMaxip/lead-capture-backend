package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

// MaxBodyBytes es el tamaño máximo aceptado para el cuerpo de una petición (1 MB).
const MaxBodyBytes = 1 << 20

// WriteJSON serializa v como JSON con el código de estado indicado.
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("escribir respuesta JSON", "error", err)
	}
}

// WriteError responde con el formato {"error": "mensaje"}.
func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]string{"error": msg})
}

// ReadJSON decodifica el cuerpo en dst, con límite de 1 MB y sin campos desconocidos.
// Los mensajes de error nunca incluyen el contenido del cuerpo.
func ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(dst); err != nil {
		var maxErr *http.MaxBytesError
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		switch {
		case errors.As(err, &maxErr):
			return errors.New("el cuerpo supera el máximo de 1 MB")
		case errors.As(err, &syntaxErr), errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("JSON mal formado")
		case errors.As(err, &typeErr):
			return fmt.Errorf("tipo inválido para el campo %q", typeErr.Field)
		case errors.Is(err, io.EOF):
			return errors.New("el cuerpo no puede estar vacío")
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			field := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("campo desconocido %s", field)
		default:
			return errors.New("JSON inválido")
		}
	}
	if dec.More() {
		return errors.New("el cuerpo debe contener un único objeto JSON")
	}
	return nil
}
