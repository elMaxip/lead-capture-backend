package auth

// Roles disponibles en el sistema.
const (
	RolCaptador       = "CAPTADOR"
	RolAdministrativo = "ADMINISTRATIVO"
)

// Usuario es un usuario del sistema, sin su contraseña.
type Usuario struct {
	ID       int    `json:"-"`
	Username string `json:"username"`
	Nombre   string `json:"nombre"`
	Rol      string `json:"rol"`
}

// LoginRequest es el cuerpo de POST /api/auth/login.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse es la respuesta de un inicio de sesión exitoso.
type LoginResponse struct {
	Token   string  `json:"token"`
	Usuario Usuario `json:"usuario"`
}
