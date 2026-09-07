package sqlserver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
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
		if len(parts) > 1 {
			idEmpleado, convErr := strconv.Atoi(parts[1])
			if convErr == nil && idEmpleado > 0 {
				return idEmpleado, nil
			}
		}
		return 0, fmt.Errorf("el procedimiento usp_go_Login no devolvió un IdEmpleado válido: %s", res)
	}

	return 0, fmt.Errorf("formato de respuesta de login inesperado: %s", res)
}

func (r *authRepository) GetMenus(ctx context.Context, idEmpleado int) ([]domain.Menu, error) {
	rows, err := r.db.QueryContext(ctx, "EXEC usp_go_MenuSeleccionarPorIdEmpleado @IdEmpleado = @p1", sql.Named("p1", idEmpleado))
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
	rows, err := r.db.QueryContext(ctx, "EXEC usp_go_MenuPermisosPorIdEmpleado @IdUsuario = @p1", sql.Named("p1", idEmpleado))
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

	rowsEmp, err := r.db.QueryContext(
		ctx,
		"EXEC usp_go_EmpleadosSeleccionarPorId @IdEmpleado = @p1",
		sql.Named("p1", idEmpleado),
	)
	if err != nil {
		return profile, fmt.Errorf("error ejecutando usp_go_EmpleadosSeleccionarPorId: %w", err)
	}
	defer rowsEmp.Close()

	if rowsEmp.Next() {
		var idEmp, idTipoEmp int
		var codPlanilla, apePat, apeMat, nombres, dni, usuario, foto sql.NullString
		var activo bool

		err := rowsEmp.Scan(
			&idEmp, &codPlanilla, &apePat, &apeMat, &nombres,
			&dni, &idTipoEmp, &activo, &usuario, &foto,
		)
		if err == nil {
			profile.Username = usuario.String
			profile.Nombres = nombres.String
			profile.ApellidoPaterno = apePat.String
			profile.ApellidoMaterno = apeMat.String
			profile.DNI = strings.TrimSpace(dni.String)
			profile.Foto = strings.TrimSpace(foto.String)

			nombreComp := strings.TrimSpace(apePat.String + " " + apeMat.String + " " + nombres.String)
			if nombreComp == "" {
				nombreComp = usuario.String
			}
			profile.NombreCompleto = nombreComp
		}
	}

	rowsMed, err := r.db.QueryContext(
		ctx,
		"EXEC usp_go_MedicosDetallePorIdEmpleado @IdEmpleado = @p1",
		sql.Named("p1", idEmpleado),
	)
	if err == nil {
		defer rowsMed.Close()
		if rowsMed.Next() {
			var idMedico, idEmpMed int
			var colegiatura, rne, idColegioHis, loteHis, rneEsp, rneEst, firmaMed, especialidad sql.NullString
			var egresado sql.NullBool

			err := rowsMed.Scan(
				&idMedico, &idEmpMed, &colegiatura, &rne, &idColegioHis,
				&loteHis, &rneEsp, &rneEst, &egresado, &firmaMed, &especialidad,
			)
			if err == nil {
				profile.Colegiatura = strings.TrimSpace(colegiatura.String)
				profile.RNE = strings.TrimSpace(rne.String)
				esp := strings.TrimSpace(especialidad.String)
				if esp != "" {
					profile.Especialidad = esp
					profile.Rol = fmt.Sprintf("Médico - %s", esp)
				}
			}
		}
	}

	if profile.Rol == "" {
		profile.Rol = "Médico"
	}

	return profile, nil
}
