package output

import (
	"context"

	"github.com/galenos-pro/appointments-api/internal/domain"
)

// AuthRepository define las operaciones de base de datos para autenticación y roles.
type AuthRepository interface {
	// Login verifica las credenciales y devuelve el IdEmpleado si son válidas.
	Login(ctx context.Context, username, password string) (int, error)
	// GetMenus obtiene los menús básicos a los que tiene acceso el empleado.
	GetMenus(ctx context.Context, idEmpleado int) ([]domain.Menu, error)
	// GetMenuPermisos obtiene los permisos detallados de los menús para el empleado.
	GetMenuPermisos(ctx context.Context, idEmpleado int) ([]domain.MenuPermiso, error)
	// GetUserProfile obtiene el perfil completo con nombres, apellidos y foto del empleado.
	GetUserProfile(ctx context.Context, idEmpleado int) (domain.UserProfile, error)
}
