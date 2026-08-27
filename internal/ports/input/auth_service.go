package input

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// AuthService es el puerto de entrada para autenticación mediante Bearer
// token JWT.
type AuthService interface {
	// Login valida credenciales y retorna un JWT firmado si son correctas.
	Login(ctx context.Context, username, password string) (string, error)

	// ValidateToken verifica la firma y expiración de un JWT y retorna
	// sus claims. Se usa en cada request para autorizar el acceso.
	ValidateToken(tokenString string) (domain.AuthClaims, error)

	// GetMenus obtiene los menús y permisos asignados al usuario.
	GetMenus(ctx context.Context, idEmpleado int) (domain.AuthMenus, error)

	// GetUserProfile obtiene el perfil con nombres, apellidos y foto del usuario.
	GetUserProfile(ctx context.Context, idEmpleado int) (domain.UserProfile, error)
}
