package prospecto

import (
	"strings"
	"testing"
	"time"
)

var cursosTest = map[int]struct{}{1: {}, 2: {}, 3: {}}

func valido() ProspectoInput {
	obs := "Interesada en horario vespertino"
	return ProspectoInput{
		ID:             "3f8a2c1e-5b7d-4e9a-8c6f-1a2b3c4d5e6f",
		Nombre:         "Camila Rojas",
		Email:          "camila.rojas@example.com",
		Telefono:       "+56912345678",
		NivelEducativo: "UNIVERSITARIA",
		CursoID:        3,
		Modalidad:      "ONLINE",
		CanalContacto:  "WHATSAPP",
		Observaciones:  &obs,
		CapturadoEn:    time.Date(2026, 10, 19, 14, 32, 0, 0, time.UTC),
	}
}

func ptr(s string) *string { return &s }

func TestValidate_Valido(t *testing.T) {
	v := NewValidator()

	casos := map[string]func(p *ProspectoInput){
		"completo":                 func(p *ProspectoInput) {},
		"sin observaciones":        func(p *ProspectoInput) { p.Observaciones = nil },
		"nombre de 2 caracteres":   func(p *ProspectoInput) { p.Nombre = "Al" },
		"nombre con tildes":        func(p *ProspectoInput) { p.Nombre = "José Núñez" },
		"nivel BASICA":             func(p *ProspectoInput) { p.NivelEducativo = "BASICA" },
		"nivel TECNICO":            func(p *ProspectoInput) { p.NivelEducativo = "TECNICO_PROFESIONAL" },
		"modalidad PRESENCIAL":     func(p *ProspectoInput) { p.Modalidad = "PRESENCIAL" },
		"canal EMAIL":              func(p *ProspectoInput) { p.CanalContacto = "EMAIL" },
		"canal TELEFONO":           func(p *ProspectoInput) { p.CanalContacto = "TELEFONO" },
		"observaciones de 500":     func(p *ProspectoInput) { p.Observaciones = ptr(strings.Repeat("a", 500)) },
		"nombre de 100":            func(p *ProspectoInput) { p.Nombre = strings.Repeat("a", 100) },
		"espacios que se recortan": func(p *ProspectoInput) { p.Nombre = "  Ana Soto  "; p.Email = " ana@example.com " },
	}

	for nombre, mod := range casos {
		t.Run(nombre, func(t *testing.T) {
			p := valido()
			mod(&p)
			Normalize(&p)
			if errs := v.Validate(p, cursosTest); len(errs) != 0 {
				t.Fatalf("se esperaba válido, errores: %v", errs)
			}
		})
	}
}

func TestValidate_Invalido(t *testing.T) {
	v := NewValidator()

	casos := []struct {
		nombre string
		campo  string
		mod    func(p *ProspectoInput)
	}{
		{"id vacío", "id", func(p *ProspectoInput) { p.ID = "" }},
		{"id no uuid", "id", func(p *ProspectoInput) { p.ID = "abc-123" }},
		{"nombre vacío", "nombre", func(p *ProspectoInput) { p.Nombre = "" }},
		{"nombre solo espacios", "nombre", func(p *ProspectoInput) { p.Nombre = "   " }},
		{"nombre de 1 caracter", "nombre", func(p *ProspectoInput) { p.Nombre = "A" }},
		{"nombre de 101", "nombre", func(p *ProspectoInput) { p.Nombre = strings.Repeat("a", 101) }},
		{"email vacío", "email", func(p *ProspectoInput) { p.Email = "" }},
		{"email sin arroba", "email", func(p *ProspectoInput) { p.Email = "camila.example.com" }},
		{"email muy largo", "email", func(p *ProspectoInput) { p.Email = strings.Repeat("a", 140) + "@example.com" }},
		{"teléfono vacío", "telefono", func(p *ProspectoInput) { p.Telefono = "" }},
		{"teléfono sin prefijo", "telefono", func(p *ProspectoInput) { p.Telefono = "912345678" }},
		{"teléfono corto", "telefono", func(p *ProspectoInput) { p.Telefono = "+5691234567" }},
		{"teléfono largo", "telefono", func(p *ProspectoInput) { p.Telefono = "+569123456789" }},
		{"teléfono fijo", "telefono", func(p *ProspectoInput) { p.Telefono = "+56221234567" }},
		{"teléfono con letras", "telefono", func(p *ProspectoInput) { p.Telefono = "+5691234567a" }},
		{"nivel vacío", "nivelEducativo", func(p *ProspectoInput) { p.NivelEducativo = "" }},
		{"nivel desconocido", "nivelEducativo", func(p *ProspectoInput) { p.NivelEducativo = "DOCTORADO" }},
		{"curso cero", "cursoId", func(p *ProspectoInput) { p.CursoID = 0 }},
		{"curso negativo", "cursoId", func(p *ProspectoInput) { p.CursoID = -1 }},
		{"curso inexistente", "cursoId", func(p *ProspectoInput) { p.CursoID = 99 }},
		{"modalidad vacía", "modalidad", func(p *ProspectoInput) { p.Modalidad = "" }},
		{"modalidad desconocida", "modalidad", func(p *ProspectoInput) { p.Modalidad = "HIBRIDA" }},
		{"canal vacío", "canalContacto", func(p *ProspectoInput) { p.CanalContacto = "" }},
		{"canal desconocido", "canalContacto", func(p *ProspectoInput) { p.CanalContacto = "SMS" }},
		{"observaciones de 501", "observaciones", func(p *ProspectoInput) { p.Observaciones = ptr(strings.Repeat("a", 501)) }},
		{"capturadoEn vacío", "capturadoEn", func(p *ProspectoInput) { p.CapturadoEn = time.Time{} }},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			p := valido()
			c.mod(&p)
			Normalize(&p)
			errs := v.Validate(p, cursosTest)
			msg, ok := errs[c.campo]
			if !ok {
				t.Fatalf("se esperaba error en %q, errores: %v", c.campo, errs)
			}
			if msg == "" {
				t.Fatalf("mensaje vacío para %q", c.campo)
			}
			if len(errs) != 1 {
				t.Fatalf("se esperaba un solo error, errores: %v", errs)
			}
		})
	}
}

func TestValidate_MensajesEnEspanol(t *testing.T) {
	v := NewValidator()
	p := valido()
	p.Email = "malo"
	p.Telefono = "123"
	p.CursoID = 99

	errs := v.Validate(p, cursosTest)

	esperados := map[string]string{
		"email":    "formato de correo inválido",
		"telefono": "debe tener el formato +569XXXXXXXX",
		"cursoId":  "el curso no existe",
	}
	for campo, msg := range esperados {
		if errs[campo] != msg {
			t.Errorf("%s: se esperaba %q, se obtuvo %q", campo, msg, errs[campo])
		}
	}
}
