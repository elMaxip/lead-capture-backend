package prospecto

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var movilChileno = regexp.MustCompile(`^\+569\d{8}$`)

// Validator valida prospectos y entrega los errores en español por campo JSON.
type Validator struct {
	v *validator.Validate
}

// NewValidator configura validator/v10 con el validador cl_movil y los nombres JSON.
func NewValidator() *Validator {
	v := validator.New(validator.WithRequiredStructEnabled())
	v.RegisterTagNameFunc(func(f reflect.StructField) string {
		name, _, _ := strings.Cut(f.Tag.Get("json"), ",")
		if name == "-" {
			return ""
		}
		return name
	})
	_ = v.RegisterValidation("cl_movil", func(fl validator.FieldLevel) bool {
		return movilChileno.MatchString(fl.Field().String())
	})
	return &Validator{v: v}
}

// Normalize quita espacios sobrantes de los campos de texto.
func Normalize(p *ProspectoInput) {
	p.ID = strings.TrimSpace(p.ID)
	p.Nombre = strings.TrimSpace(p.Nombre)
	p.Email = strings.TrimSpace(p.Email)
	p.Telefono = strings.TrimSpace(p.Telefono)
	if p.Observaciones != nil {
		obs := strings.TrimSpace(*p.Observaciones)
		if obs == "" {
			p.Observaciones = nil
		} else {
			p.Observaciones = &obs
		}
	}
}

// Validate devuelve un mapa campo → mensaje, vacío si el prospecto es válido.
// cursos es el conjunto de ids de curso existentes.
func (val *Validator) Validate(p ProspectoInput, cursos map[int]struct{}) map[string]string {
	errs := map[string]string{}

	var verrs validator.ValidationErrors
	if err := val.v.Struct(p); errors.As(err, &verrs) {
		for _, fe := range verrs {
			errs[fe.Field()] = mensaje(fe)
		}
	} else if err != nil {
		errs["prospecto"] = "no se pudo validar"
	}

	if _, ok := errs["cursoId"]; !ok {
		if _, existe := cursos[p.CursoID]; !existe {
			errs["cursoId"] = "el curso no existe"
		}
	}
	return errs
}

func mensaje(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "es obligatorio"
	case "uuid":
		return "debe ser un UUID válido"
	case "min":
		return fmt.Sprintf("debe tener al menos %s caracteres", fe.Param())
	case "max":
		return fmt.Sprintf("debe tener como máximo %s caracteres", fe.Param())
	case "email":
		return "formato de correo inválido"
	case "cl_movil":
		return "debe tener el formato +569XXXXXXXX"
	case "oneof":
		return "debe ser uno de: " + strings.ReplaceAll(fe.Param(), " ", ", ")
	case "gt":
		return "debe ser mayor que " + fe.Param()
	default:
		return "valor inválido"
	}
}
