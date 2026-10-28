// Package auth implementa la lógica del servicio de autenticación.
package auth

// Usuario representa un usuario del sistema.
// No se manejan roles: todos los usuarios tienen las mismas capacidades.
type Usuario struct {
	ID       int    `json:"id"`
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

// UsuarioSinPassword retorna el usuario sin el campo password.
type UsuarioSinPassword struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
	Email  string `json:"email"`
}

// LoginRequest es el body de POST /auth/login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterRequest es el body de POST /auth/register.
type RegisterRequest struct {
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// SinPassword convierte un Usuario a UsuarioSinPassword.
func (u *Usuario) SinPassword() UsuarioSinPassword {
	return UsuarioSinPassword{
		ID:     u.ID,
		Nombre: u.Nombre,
		Email:  u.Email,
	}
}
