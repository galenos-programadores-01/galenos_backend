package sqlserver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) output.AuthRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Login(ctx context.Context, username, password string) (int, error) {
	var resultStr sql.NullString
	// Procedimiento almacenado usp_go_Login devuelve un string como salida
	// "SUCCESS;IdEmpleado;..." o "ERROR;Mensaje"
	_, err := r.db.ExecContext(ctx, "EXEC usp_go_Login @Usuario = @p1, @Password = @p2, @Resultado = @p3 OUTPUT",
		sql.Named("p1", username),
		sql.Named("p2", password),
		sql.Named("p3", sql.Out{Dest: &resultStr}),
	)
	if err != nil {
		return 0, fmt.Errorf("error ejecutando usp_go_Login: %w", err)
	}

	res := resultStr.String
	log.Printf("DB Login Response: %q", res)

	parts := strings.Split(res, ";")
	if len(parts) > 0 && parts[0] == "ERROR" {
		msg := "Credenciales inválidas"
		if len(parts) > 1 {
			msg = parts[1]
		}
		return 0, fmt.Errorf("%w: %s", domain.ErrInvalidCredentials, msg)
	}

	if len(parts) > 0 && parts[0] == "OK" {
		var idEmpleado int
		// Buscamos el IdEmpleado porque el SP no lo devuelve en el parámetro OUTPUT
		err = r.db.QueryRowContext(ctx, "SELECT IdEmpleado FROM dbo.Empleados WHERE Usuario = @p1", sql.Named("p1", username)).Scan(&idEmpleado)
		if err != nil {
			return 0, fmt.Errorf("error obteniendo IdEmpleado tras login exitoso: %w", err)
		}
		return idEmpleado, nil
	}

	return 0, fmt.Errorf("formato de respuesta de login inesperado: %s", res)
}

func (r *authRepository) GetMenus(ctx context.Context, idEmpleado int) ([]domain.Menu, error) {
	rows, err := r.db.QueryContext(ctx, "EXEC webMenuSeleccionarIdEmpleado @IdEmpleado = @p1", sql.Named("p1", idEmpleado))
	if err != nil {
		return nil, fmt.Errorf("error consultando menus: %w", err)
	}
	defer rows.Close()

	var menus []domain.Menu
	for rows.Next() {
		var m domain.Menu
		var texto, keyIconWeb, claveWeb sql.NullString
		if err := rows.Scan(&m.IdListGrupo, &texto, &keyIconWeb, &claveWeb, &m.Indice, &m.Estado, &m.NroSubMenu); err != nil {
			return nil, fmt.Errorf("error escaneando menu: %w", err)
		}
		m.Texto = texto.String
		m.KeyIconWeb = keyIconWeb.String
		m.ClaveWeb = claveWeb.String
		menus = append(menus, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando menus: %w", err)
	}

	return menus, nil
}

func (r *authRepository) GetMenuPermisos(ctx context.Context, idEmpleado int) ([]domain.MenuPermiso, error) {
	rows, err := r.db.QueryContext(ctx, "EXEC web_MenuPermisosIdempleado @IdUsuario = @p1", sql.Named("p1", idEmpleado))
	if err != nil {
		return nil, fmt.Errorf("error consultando menu permisos: %w", err)
	}
	defer rows.Close()

	var permisos []domain.MenuPermiso
	for rows.Next() {
		var p domain.MenuPermiso
		var opciones, texto, menu, keyIconWeb, claveWeb sql.NullString
		if err := rows.Scan(&opciones, &p.Indice, &texto, &menu, &p.IdListGrupo, &keyIconWeb, &p.Estado, &claveWeb, &p.Agregar, &p.Modificar, &p.Eliminar); err != nil {
			if !errors.Is(err, sql.ErrNoRows) {
				continue
			}
		}
		p.Opciones = opciones.String
		p.Texto = texto.String
		p.Menu = menu.String
		p.KeyIconWeb = keyIconWeb.String
		p.ClaveWeb = claveWeb.String
		permisos = append(permisos, p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterando permisos: %w", err)
	}

	return permisos, nil
}

func (r *authRepository) GetUserProfile(ctx context.Context, idEmpleado int) (domain.UserProfile, error) {
	var profile domain.UserProfile
	profile.IdEmpleado = idEmpleado

	var username, nombres, apePat, apeMat, foto, cargo sql.NullString
	err := r.db.QueryRowContext(ctx, `
		SELECT 
			ISNULL(em.Usuario, ''),
			ISNULL(em.Nombres, ''),
			ISNULL(em.ApellidoPaterno, ''),
			ISNULL(em.ApellidoMaterno, ''),
			ISNULL(em.Foto, ''),
			ISNULL((SELECT TOP 1 c.Nombre FROM dbo.Cargos c WHERE c.IdCargo = em.IdCargo), 'Médico Tratante')
		FROM dbo.Empleados em
		WHERE em.IdEmpleado = @p1`,
		sql.Named("p1", idEmpleado),
	).Scan(&username, &nombres, &apePat, &apeMat, &foto, &cargo)

	if err != nil {
		err = r.db.QueryRowContext(ctx, `
			SELECT 
				ISNULL(em.Usuario, ''),
				ISNULL(em.Nombres, ''),
				ISNULL(em.ApellidoPaterno, ''),
				ISNULL(em.ApellidoMaterno, ''),
				ISNULL(em.Foto, '')
			FROM dbo.Empleados em
			WHERE em.IdEmpleado = @p1`,
			sql.Named("p1", idEmpleado),
		).Scan(&username, &nombres, &apePat, &apeMat, &foto)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return profile, nil
			}
			return profile, fmt.Errorf("error obteniendo perfil de empleado: %w", err)
		}
	}

	profile.Username = username.String
	profile.Nombres = nombres.String
	profile.ApellidoPaterno = apePat.String
	profile.ApellidoMaterno = apeMat.String

	nombreComp := strings.TrimSpace(apePat.String + " " + apeMat.String + " " + nombres.String)
	if nombreComp == "" {
		nombreComp = username.String
	}
	profile.NombreCompleto = nombreComp
	profile.Foto = foto.String
	if cargo.Valid && cargo.String != "" {
		profile.Rol = cargo.String
	} else {
		profile.Rol = "Médico Tratante"
	}

	return profile, nil
}
