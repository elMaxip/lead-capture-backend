package prospecto

import "time"

// MaxPorLote es la cantidad máxima de prospectos aceptada en una sincronización.
const MaxPorLote = 100

// ProspectoInput es un prospecto tal como lo envía la app al sincronizar.
type ProspectoInput struct {
	ID             string    `json:"id" validate:"required,uuid"`
	Nombre         string    `json:"nombre" validate:"required,min=2,max=100"`
	Email          string    `json:"email" validate:"required,email,max=150"`
	Telefono       string    `json:"telefono" validate:"required,cl_movil"`
	NivelEducativo string    `json:"nivelEducativo" validate:"required,oneof=BASICA MEDIA TECNICO_PROFESIONAL UNIVERSITARIA POSTGRADO"`
	CursoID        int       `json:"cursoId" validate:"required,gt=0"`
	Modalidad      string    `json:"modalidad" validate:"required,oneof=ONLINE PRESENCIAL"`
	CanalContacto  string    `json:"canalContacto" validate:"required,oneof=EMAIL TELEFONO WHATSAPP"`
	Observaciones  *string   `json:"observaciones" validate:"omitempty,max=500"`
	CapturadoEn    time.Time `json:"capturadoEn" validate:"required"`
}

// SyncRequest es el cuerpo de POST /api/prospectos/sync.
type SyncRequest struct {
	Prospectos []ProspectoInput `json:"prospectos"`
}

// Rechazado describe un prospecto que no pasó la validación.
type Rechazado struct {
	ID      string            `json:"id"`
	Errores map[string]string `json:"errores"`
}

// SyncResponse es la respuesta de POST /api/prospectos/sync.
type SyncResponse struct {
	Sincronizados []string    `json:"sincronizados"`
	Rechazados    []Rechazado `json:"rechazados"`
}

// Prospecto es un registro consolidado, con los nombres de curso y área.
type Prospecto struct {
	ID             string    `json:"id"`
	Nombre         string    `json:"nombre"`
	Email          string    `json:"email"`
	Telefono       string    `json:"telefono"`
	NivelEducativo string    `json:"nivelEducativo"`
	CursoID        int       `json:"cursoId"`
	CursoNombre    string    `json:"cursoNombre"`
	AreaNombre     string    `json:"areaNombre"`
	Modalidad      string    `json:"modalidad"`
	CanalContacto  string    `json:"canalContacto"`
	Observaciones  *string   `json:"observaciones"`
	CapturadoEn    time.Time `json:"capturadoEn"`
	RecibidoEn     time.Time `json:"recibidoEn"`
}
