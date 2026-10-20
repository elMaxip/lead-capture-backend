package catalogo

// Curso es un curso ofrecido dentro de un área.
type Curso struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

// Area agrupa cursos de una misma temática.
type Area struct {
	ID     int     `json:"id"`
	Nombre string  `json:"nombre"`
	Cursos []Curso `json:"cursos"`
}

// Catalogo es la respuesta de GET /api/catalogo.
type Catalogo struct {
	Areas []Area `json:"areas"`
}
