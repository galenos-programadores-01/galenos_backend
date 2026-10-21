package sqlserver

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type atencionCERepository struct {
	db *sql.DB
}

func NewAtencionCERepository(db *sql.DB) output.AtencionCERepository {
	return &atencionCERepository{
		db: db,
	}
}

func (r *atencionCERepository) ListarColaPacientes(ctx context.Context, filtros domain.ColaFiltrosCE) ([]domain.ColaPacienteCEItem, error) {
	query := `EXEC dbo.usp_go_AtencionesCE_ListarCola 
		@IdMedico = @p1, 
		@IdServicio = @p2, 
		@Fecha = @p3, 
		@Filtro = @p4`

	var fechaParam any = nil
	if filtros.Fecha != "" {
		fechaParam = filtros.Fecha
	}

	filtroLimpio := strings.TrimSpace(filtros.Filtro)
	var filtroParam any = nil
	if filtroLimpio != "" {
		filtroParam = filtroLimpio
	}

	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("p1", filtros.IdMedico),
		sql.Named("p2", filtros.IdServicio),
		sql.Named("p3", fechaParam),
		sql.Named("p4", filtroParam),
	)
	if err != nil {
		return nil, fmt.Errorf("error al listar cola de pacientes CE mediante SP: %w", err)
	}
	defer rows.Close()

	cola := make([]domain.ColaPacienteCEItem, 0)
	for rows.Next() {
		var item domain.ColaPacienteCEItem
		err := rows.Scan(
			&item.IdAtencion,
			&item.IdCuentaAtencion,
			&item.IdPaciente,
			&item.NroHistoriaClinica,
			&item.NroDocumento,
			&item.PacienteNombre,
			&item.FechaIngreso,
			&item.HoraIngreso,
			&item.HoraEgreso,
			&item.Edad,
			&item.IdServicioIngreso,
			&item.ServicioNombre,
			&item.IdMedicoIngreso,
			&item.MedicoNombre,
			&item.EstadoAtencion,
			&item.IdCita,
			&item.TipoFinanciador,
			&item.Sexo,
			&item.Telefono,
			&item.TieneTriaje,
			&item.Turno,
		)
		if err != nil {
			return nil, fmt.Errorf("error al escanear cola de paciente CE: %w", err)
		}
		cola = append(cola, item)
	}

	return cola, rows.Err()
}

func (r *atencionCERepository) ObtenerFichaAtencion(ctx context.Context, identificadorAtencion int) (*domain.FichaAtencionCEDetalle, error) {
	consultaCabecera := `
		SELECT 
			a.IdAtencion,
			ISNULL(a.IdCuentaAtencion, 0) AS IdCuentaAtencion,
			ISNULL(a.IdPaciente, 0) AS IdPaciente,
			LTRIM(RTRIM(CONCAT(
				ISNULL(p.ApellidoPaterno, ''), ' ', 
				ISNULL(p.ApellidoMaterno, ''), ' ', 
				ISNULL(p.PrimerNombre, ''), ' ', 
				ISNULL(p.SegundoNombre, '')
			))) AS PacienteNombre,
			ISNULL(p.NroHistoriaClinica, 0) AS NroHistoriaClinica,
			LTRIM(RTRIM(ISNULL(p.NroDocumento, ''))) AS NroDocumento,
			ISNULL(a.Edad, 0) AS Edad,
			CASE 
				WHEN p.IdTipoSexo = 1 THEN 'M' 
				WHEN p.IdTipoSexo = 2 THEN 'F' 
				ELSE '' 
			END AS Sexo,
			ISNULL(ace.Motivo, '') AS Motivo,
			ISNULL(ace.ExamenClinico, '') AS ExamenClinico,
			ISNULL(ace.Antecedente, '') AS Antecedente,
			ISNULL(ace.InforCompDiagMed, '') AS InforCompDiagMed,
			ISNULL(ace.Tratamiento, '') AS Tratamiento,
			ISNULL(ace.PlanTratamiento, '') AS PlanTratamiento,
			ISNULL(ace.Observaciones, '') AS Observaciones,
			ISNULL(ace.FechaProxCita, '') AS FechaProxCita,
			ISNULL(a.IdDestinoAtencion, 0) AS IdDestinoAtencion,
			LTRIM(RTRIM(ISNULL(tda.Descripcion, ''))) AS DestinoNombre,
			ISNULL(a.HoraIngreso, '') AS HoraIngreso,
			ISNULL(a.HoraEgreso, '') AS HoraEgreso,
			ISNULL(ace.Estado, 0) AS Estado,
			CONVERT(VARCHAR(10), a.FechaIngreso, 120) AS FechaAtencion,
			CONVERT(VARCHAR(10), a.FechaIngreso, 120) AS FechaIngreso,
			LTRIM(RTRIM(ISNULL(s.Nombre, 'CONSULTA EXTERNA'))) AS EspecialidadNombre,
			LTRIM(RTRIM(CONCAT(
				ISNULL(e.ApellidoPaterno, ''), ' ', 
				ISNULL(e.ApellidoMaterno, ''), ' ', 
				ISNULL(e.Nombres, '')
			))) AS MedicoNombre,
			LTRIM(RTRIM(ISNULL(ff.Descripcion, 'SIS'))) AS Iafa,
			ISNULL(NULLIF(LTRIM(RTRIM(p.celular)), ''), ISNULL(LTRIM(RTRIM(p.Telefono)), '')) AS Telefono
		FROM dbo.Atenciones a WITH (NOLOCK)
		LEFT JOIN dbo.Pacientes p WITH (NOLOCK) ON a.IdPaciente = p.IdPaciente
		LEFT JOIN dbo.AtencionesCE ace WITH (NOLOCK) ON a.IdAtencion = ace.idAtencion
		LEFT JOIN dbo.TiposDestinoAtencion tda WITH (NOLOCK) ON a.IdDestinoAtencion = tda.IdDestinoAtencion
		LEFT JOIN dbo.Servicios s WITH (NOLOCK) ON a.IdServicioIngreso = s.IdServicio
		LEFT JOIN dbo.Medicos m WITH (NOLOCK) ON a.IdMedicoIngreso = m.IdMedico
		LEFT JOIN dbo.Empleados e WITH (NOLOCK) ON m.IdEmpleado = e.IdEmpleado
		LEFT JOIN dbo.FuentesFinanciamiento ff WITH (NOLOCK) ON a.idFuenteFinanciamiento = ff.IdFuenteFinanciamiento
		WHERE a.IdAtencion = @pIdentificadorAtencion`

	var ficha domain.FichaAtencionCEDetalle
	errorCabecera := r.db.QueryRowContext(ctx, consultaCabecera, sql.Named("pIdentificadorAtencion", identificadorAtencion)).Scan(
		&ficha.IdAtencion,
		&ficha.IdCuentaAtencion,
		&ficha.IdPaciente,
		&ficha.PacienteNombre,
		&ficha.NroHistoriaClinica,
		&ficha.NroDocumento,
		&ficha.Edad,
		&ficha.Sexo,
		&ficha.Motivo,
		&ficha.ExamenClinico,
		&ficha.Antecedente,
		&ficha.InforCompDiagMed,
		&ficha.Tratamiento,
		&ficha.PlanTratamiento,
		&ficha.Observaciones,
		&ficha.FechaProxCita,
		&ficha.IdDestinoAtencion,
		&ficha.DestinoNombre,
		&ficha.HoraIngreso,
		&ficha.HoraEgreso,
		&ficha.Estado,
		&ficha.FechaAtencion,
		&ficha.FechaIngreso,
		&ficha.EspecialidadNombre,
		&ficha.MedicoNombre,
		&ficha.Iafa,
		&ficha.Telefono,
	)
	if errorCabecera != nil {
		if errorCabecera == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error al obtener cabecera de atencion CE: %w", errorCabecera)
	}

	if ficha.Estado != 2 {
		actualizarEstadoAtendiendo := `
			UPDATE dbo.Atenciones
			SET idEstadoAtencion = 1,
			    HoraIngreso = CASE WHEN HoraIngreso IS NULL OR RTRIM(HoraIngreso) = '' THEN CONVERT(VARCHAR(5), GETDATE(), 108) ELSE HoraIngreso END
			WHERE IdAtencion = @p1 AND idEstadoAtencion <> 2;

			UPDATE dbo.Citas
			SET IdEstadoCita = 8
			WHERE IdAtencion = @p1 AND IdEstadoCita <> 4;
		`
		_, _ = r.db.ExecContext(ctx, actualizarEstadoAtendiendo, sql.Named("p1", identificadorAtencion))
	}

	consultaDiagnosticos := `EXEC [dbo].[usp_go_ObtenerDiagnosticosAtencion] @IdAtencion = @pIdentificadorAtencion, @IdPrimeraAtencion = 0, @IdEvolucion = 0`
	filasDiagnosticos, errorDiagnosticos := r.db.QueryContext(ctx, consultaDiagnosticos, sql.Named("pIdentificadorAtencion", identificadorAtencion))
	if errorDiagnosticos != nil {
		return nil, fmt.Errorf("error al obtener diagnosticos de atencion CE: %w", errorDiagnosticos)
	}
	defer filasDiagnosticos.Close()

	ficha.Diagnosticos = make([]domain.AtencionDiagnosticoCEItem, 0)
	for filasDiagnosticos.Next() {
		var diagnostico domain.AtencionDiagnosticoCEItem
		var idAd, idAtSec, idDx, idSubDx sql.NullInt64
		var cie, desc, descL, tCod, tDx, his, his1, his2, his3 sql.NullString

		if errorEscaneo := filasDiagnosticos.Scan(
			&idAd,
			&idAtSec,
			&idDx,
			&cie,
			&desc,
			&descL,
			&tCod,
			&tDx,
			&idSubDx,
			&his,
			&his1,
			&his2,
			&his3,
		); errorEscaneo != nil {
			return nil, fmt.Errorf("error al escanear diagnostico CE: %w", errorEscaneo)
		}

		diagnostico.IdAtencionDiagnostico = int(idAd.Int64)
		diagnostico.IdDiagnostico = int(idDx.Int64)
		diagnostico.CodigoCIE10 = strings.TrimSpace(cie.String)
		diagnostico.Descripcion = strings.TrimSpace(desc.String)

		codigoClasificacion := strings.ToUpper(strings.TrimSpace(tCod.String))
		subclasificacionVal := int(idSubDx.Int64)
		if subclasificacionVal == 101 || subclasificacionVal == 401 || subclasificacionVal == 404 || subclasificacionVal == 407 || subclasificacionVal == 410 {
			codigoClasificacion = "P"
		} else if subclasificacionVal == 103 || subclasificacionVal == 403 || subclasificacionVal == 406 || subclasificacionVal == 412 {
			codigoClasificacion = "R"
		} else if subclasificacionVal == 102 || subclasificacionVal == 402 || subclasificacionVal == 405 || subclasificacionVal == 408 || subclasificacionVal == 411 {
			codigoClasificacion = "D"
		}
		if codigoClasificacion == "" {
			codigoClasificacion = "D"
		}

		diagnostico.Clasificacion = codigoClasificacion
		if diagnostico.Clasificacion == "P" {
			diagnostico.IdClasificacionDx = 1
		} else if diagnostico.Clasificacion == "R" {
			diagnostico.IdClasificacionDx = 3
		} else {
			diagnostico.IdClasificacionDx = 2
		}

		if len(ficha.Diagnosticos) == 0 || subclasificacionVal == 1 {
			diagnostico.IdSubclasificacionDx = 1
		} else {
			diagnostico.IdSubclasificacionDx = 2
		}

		diagnostico.LabConfHIS = strings.TrimSpace(his.String)
		diagnostico.LabConfHIS1 = strings.TrimSpace(his1.String)
		diagnostico.LabConfHIS2 = strings.TrimSpace(his2.String)
		diagnostico.LabConfHIS3 = strings.TrimSpace(his3.String)

		if diagnostico.LabConfHIS1 == "" && diagnostico.LabConfHIS != "" && diagnostico.LabConfHIS != "P" && diagnostico.LabConfHIS != "D" && diagnostico.LabConfHIS != "R" {
			diagnostico.LabConfHIS1 = diagnostico.LabConfHIS
		}

		if diagnostico.LabConfHIS1 == "" && diagnostico.LabConfHIS2 == "" && diagnostico.LabConfHIS3 == "" {
			var dHis, dHis1, dHis2, dHis3 sql.NullString
			_ = r.db.QueryRowContext(ctx, `
				SELECT TOP 1 
					ISNULL(labConfHIS, ''), 
					ISNULL(labConfHIS_1, ''), 
					ISNULL(labConfHIS_2, ''), 
					ISNULL(labConfHIS_3, '')
				FROM dbo.AtencionesDiagnosticos WITH (NOLOCK)
				WHERE IdAtencion = @p1 AND IdDiagnostico = @p2`,
				sql.Named("p1", identificadorAtencion),
				sql.Named("p2", diagnostico.IdDiagnostico),
			).Scan(&dHis, &dHis1, &dHis2, &dHis3)
			if dHis1.Valid && strings.TrimSpace(dHis1.String) != "" {
				diagnostico.LabConfHIS1 = strings.TrimSpace(dHis1.String)
			}
			if dHis2.Valid && strings.TrimSpace(dHis2.String) != "" {
				diagnostico.LabConfHIS2 = strings.TrimSpace(dHis2.String)
			}
			if dHis3.Valid && strings.TrimSpace(dHis3.String) != "" {
				diagnostico.LabConfHIS3 = strings.TrimSpace(dHis3.String)
			}
			if diagnostico.LabConfHIS1 == "" && dHis.Valid && strings.TrimSpace(dHis.String) != "" && dHis.String != "P" && dHis.String != "D" && dHis.String != "R" {
				diagnostico.LabConfHIS1 = strings.TrimSpace(dHis.String)
			}
		}

		ficha.Diagnosticos = append(ficha.Diagnosticos, diagnostico)
	}

	ficha.ItemsReceta = make([]domain.ItemRecetaPrescripcion, 0)
	consultaReceta := `
		SELECT 
			ISNULL(rd.idItem, 0) AS IdProducto,
			ISNULL(b.Codigo, '') AS Codigo,
			ISNULL(NULLIF(b.Nombre, ''), ISNULL(rd.DescripcionadicionalReceta, '')) AS Nombre,
			ISNULL(b.Presentacion, '') AS FormaFarmaceutica,
			ISNULL(b.Concentracion, '') AS Concentracion,
			ISNULL(rd.CantidadPedida, 0) AS Cantidad,
			ISNULL(CAST(rd.Duracion AS VARCHAR(50)), '') AS Duracion,
			ISNULL(rd.observaciones, '') AS Indicaciones,
			ISNULL(rd.DescripcionadicionalReceta, '') AS DescripcionAdicional,
			COALESCE(NULLIF(LTRIM(RTRIM(rd.CodigoDiagnostico)), ''), dId.CodigoCIE10, dCie.CodigoCIE10, '') AS CodigoCIE10,
			COALESCE(NULLIF(rd.IdDiagnostico, 0), dId.IdDiagnostico, dCie.IdDiagnostico, 0) AS IdDiagnostico,
			ISNULL(va.Descripcion, '') AS ViaNombre,
			ISNULL(frec.Descripcion, '') AS FrecuenciaNombre
		FROM dbo.RecetaCabecera rc WITH (NOLOCK)
		INNER JOIN dbo.RecetaDetalle rd WITH (NOLOCK) ON rc.idReceta = rd.idReceta
		LEFT JOIN dbo.FactCatalogoBienesInsumos b WITH (NOLOCK) ON b.IdProducto = rd.idItem
		LEFT JOIN dbo.RecetaViaAdministracion va WITH (NOLOCK) ON va.IdViaAdministracion = rd.IdViaAdministracion
		LEFT JOIN dbo.TabRecetaFrecuencia frec WITH (NOLOCK) ON frec.idFrecuencia = rd.idFrecuencia
		LEFT JOIN dbo.Diagnosticos dId WITH (NOLOCK) ON dId.IdDiagnostico = rd.IdDiagnostico
		LEFT JOIN dbo.Diagnosticos dCie WITH (NOLOCK) ON (rd.IdDiagnostico IS NULL OR rd.IdDiagnostico = 0) AND dCie.CodigoCIE10 = rd.CodigoDiagnostico
		WHERE (rc.lnIdRegistroAtencion = @p1 OR (rc.idCuentaAtencion = @p2 AND @p2 > 0))
		  AND rc.IdPuntoCarga = 5
		ORDER BY rd.idItem ASC`

	filasReceta, errorReceta := r.db.QueryContext(ctx, consultaReceta,
		sql.Named("p1", identificadorAtencion),
		sql.Named("p2", ficha.IdCuentaAtencion),
	)
	if errorReceta == nil {
		defer filasReceta.Close()
		for filasReceta.Next() {
			var itemReceta domain.ItemRecetaPrescripcion
			var idProd, cant, idDx sql.NullInt64
			var cod, nom, forma, conc, dur, ind, descAd, cie, viaNom, frecNom sql.NullString
			if errorEscaneoReceta := filasReceta.Scan(
				&idProd,
				&cod,
				&nom,
				&forma,
				&conc,
				&cant,
				&dur,
				&ind,
				&descAd,
				&cie,
				&idDx,
				&viaNom,
				&frecNom,
			); errorEscaneoReceta == nil {
				itemReceta.IdProducto = int(idProd.Int64)
				itemReceta.Codigo = strings.TrimSpace(cod.String)
				itemReceta.Nombre = strings.TrimSpace(nom.String)
				itemReceta.FormaFarmaceutica = strings.TrimSpace(forma.String)
				itemReceta.Concentracion = strings.TrimSpace(conc.String)
				itemReceta.Cantidad = int(cant.Int64)
				itemReceta.Indicaciones = strings.TrimSpace(ind.String)
				itemReceta.CodigoCIE10 = strings.TrimSpace(cie.String)
				itemReceta.IdDiagnostico = int(idDx.Int64)

				descripcionExtra := strings.TrimSpace(descAd.String)
				if strings.Contains(descripcionExtra, " | ") {
					partesBarra := strings.Split(descripcionExtra, " | ")
					if len(partesBarra) >= 5 {
						itemReceta.Dosis = strings.TrimSpace(partesBarra[0])
						itemReceta.UnidadMedida = strings.TrimSpace(partesBarra[1])
						itemReceta.Frecuencia = strings.TrimSpace(partesBarra[2])
						itemReceta.Via = strings.TrimSpace(partesBarra[3])
						itemReceta.Duracion = strings.TrimSpace(partesBarra[4])
					} else if len(partesBarra) >= 4 {
						itemReceta.Dosis = strings.TrimSpace(partesBarra[0])
						itemReceta.UnidadMedida = strings.TrimSpace(partesBarra[1])
						itemReceta.Frecuencia = strings.TrimSpace(partesBarra[2])
						itemReceta.Via = strings.TrimSpace(partesBarra[3])
					} else if len(partesBarra) >= 3 {
						itemReceta.Dosis = strings.TrimSpace(partesBarra[0])
						itemReceta.UnidadMedida = strings.TrimSpace(partesBarra[1])
						itemReceta.Frecuencia = strings.TrimSpace(partesBarra[2])
					}
				} else if strings.Contains(descripcionExtra, " - ") {
					partesGuion := strings.Split(descripcionExtra, " - ")
					if len(partesGuion) >= 3 {
						itemReceta.Dosis = strings.TrimSpace(partesGuion[0])
						itemReceta.Frecuencia = strings.TrimSpace(partesGuion[1])
						itemReceta.Via = strings.TrimSpace(partesGuion[2])
					} else if len(partesGuion) == 2 {
						itemReceta.Dosis = strings.TrimSpace(partesGuion[0])
						itemReceta.Frecuencia = strings.TrimSpace(partesGuion[1])
						itemReceta.Via = "ORAL"
					} else if len(partesGuion) == 1 && partesGuion[0] != "" {
						itemReceta.Dosis = strings.TrimSpace(partesGuion[0])
						itemReceta.Via = "ORAL"
					}
				}

				if itemReceta.Via == "" && strings.TrimSpace(viaNom.String) != "" {
					itemReceta.Via = strings.ToUpper(strings.TrimSpace(viaNom.String))
				}
				if itemReceta.Frecuencia == "" && strings.TrimSpace(frecNom.String) != "" {
					itemReceta.Frecuencia = strings.ToUpper(strings.TrimSpace(frecNom.String))
				}
				if itemReceta.UnidadMedida == "" {
					if itemReceta.FormaFarmaceutica != "" {
						itemReceta.UnidadMedida = itemReceta.FormaFarmaceutica
					} else {
						itemReceta.UnidadMedida = "UNIDAD"
					}
				}
				if itemReceta.Dosis == "" {
					itemReceta.Dosis = "1"
				}
				if itemReceta.Via == "" {
					itemReceta.Via = "ORAL"
				}
				if itemReceta.Frecuencia == "" {
					itemReceta.Frecuencia = "CADA 24 HORAS"
				}
				if itemReceta.Duracion == "" {
					duracionCadena := strings.TrimSpace(dur.String)
					if duracionCadena != "" && duracionCadena != "0" {
						if !strings.Contains(strings.ToUpper(duracionCadena), "D") {
							itemReceta.Duracion = duracionCadena + " DÍAS"
						} else {
							itemReceta.Duracion = duracionCadena
						}
					} else {
						itemReceta.Duracion = "30 DÍAS"
					}
				}

				if itemReceta.Nombre != "" {
					ficha.ItemsReceta = append(ficha.ItemsReceta, itemReceta)
				}
			}
		}
	}

	ficha.Interconsultas = make([]domain.InterconsultaCEItem, 0)
	consultaInterconsultas := `
		SELECT
			i.idInterconsulta,
			ISNULL(i.IdRegistroAtencion, @pIdentificadorAtencion) AS IdAtencionOrigen,
			ISNULL(i.IdEspecialidadSolicitada, 0) AS IdEspecialidad,
			ISNULL(i.idEmpleadoRespuesta, 0)      AS IdMedicoDestino,
			ISNULL(i.MotivoClinico, '')           AS Motivo,
			i.fechaSolicitud                      AS FechaSolicitud,
			CASE WHEN i.idEstado = 1 THEN 'SOLICITUD' ELSE 'RESPUESTA' END AS Estado,
			ISNULL(esp.Nombre, '')                AS Especialidad,
			ISNULL(i.ResumenClinico, ISNULL(i.ResumenHC, '')) AS ResumenHC
		FROM dbo.Interconsulta i WITH (NOLOCK)
		LEFT JOIN dbo.Especialidades esp WITH (NOLOCK) ON esp.IdEspecialidad = i.IdEspecialidadSolicitada
		WHERE (i.IdRegistroAtencion = @pIdentificadorAtencion 
		   OR i.IdPrimeraAtencion = @pIdentificadorAtencion)

		ORDER BY i.fechaSolicitud ASC, i.idInterconsulta ASC`

	filasInterconsultas, errorInterconsultas := r.db.QueryContext(ctx, consultaInterconsultas, sql.Named("pIdentificadorAtencion", identificadorAtencion))
	if errorInterconsultas == nil {
		defer filasInterconsultas.Close()
		for filasInterconsultas.Next() {
			var interconsulta domain.InterconsultaCEItem
			var identificadorAtencionOrigen int
			var identificadorMedicoDestino int
			var fechaSolicitud sql.NullTime
			var estadoInterconsulta string
			var resumenClinico string
			if errorEscaneoInterconsulta := filasInterconsultas.Scan(
				&interconsulta.IdInterconsulta,
				&identificadorAtencionOrigen,
				&interconsulta.IdEspecialidad,
				&identificadorMedicoDestino,
				&interconsulta.Motivo,
				&fechaSolicitud,
				&estadoInterconsulta,
				&interconsulta.NombreEspecialidad,
				&resumenClinico,
			); errorEscaneoInterconsulta == nil {
				if fechaSolicitud.Valid {
					interconsulta.FechaPropuesta = fechaSolicitud.Time.Format("2006-01-02")
				}
				interconsulta.ResumenEnfermedad = resumenClinico
				interconsulta.IdMedico = identificadorMedicoDestino
				ficha.Interconsultas = append(ficha.Interconsultas, interconsulta)
			}
		}
	}

	ficha.OrdenesApoyoDiagnostico = make([]domain.OrdenApoyoDiagnosticoItem, 0)
	consultaApoyo := `
		WITH ItemsApoyoCTE AS (
			SELECT 
				ISNULL(rd.idItem, 0) AS IdProducto,
				ISNULL(fcs.Codigo, '') AS Codigo,
				ISNULL(NULLIF(LTRIM(RTRIM(rd.DescripcionadicionalReceta)), ''), ISNULL(fcs.Nombre, 'EXAMEN AUXILIAR')) AS NombreProcedimiento,
				CASE 
					WHEN rc.IdPuntoCarga IN (2, 3, 31, 32, 33, 34, 35, 36, 37, 38) THEN 'LABORATORIO CLÍNICO'
					WHEN rc.IdPuntoCarga IN (4, 12, 20, 21, 22, 23, 24) THEN 'DIAGNÓSTICO POR IMÁGENES'
					WHEN s.Nombre IS NOT NULL AND LTRIM(RTRIM(s.Nombre)) <> '' THEN s.Nombre
					WHEN pc.Descripcion IS NOT NULL AND LTRIM(RTRIM(pc.Descripcion)) <> '' THEN pc.Descripcion
					ELSE 'APOYO AL DIAGNÓSTICO'
				END AS Servicio,
				ISNULL(rd.CantidadPedida, 1) AS Cantidad,
				ISNULL(CAST(rd.Justificacion AS VARCHAR(MAX)), '') AS Justificacion,
				ISNULL(rd.observaciones, '') AS Observacion,
				ISNULL(rd.Precio, 0.0) AS Precio,
				COALESCE(NULLIF(LTRIM(RTRIM(rd.CodigoDiagnostico)), ''), dId.CodigoCIE10, dCie.CodigoCIE10, '') AS CodigoCIE10,
				COALESCE(NULLIF(rd.IdDiagnostico, 0), dId.IdDiagnostico, dCie.IdDiagnostico, 0) AS IdDiagnostico,
				ROW_NUMBER() OVER (PARTITION BY rd.idItem ORDER BY rc.idReceta DESC) AS FilaNumero
			FROM dbo.RecetaCabecera rc WITH (NOLOCK)
			INNER JOIN dbo.RecetaDetalle rd WITH (NOLOCK) ON rc.idReceta = rd.idReceta
			LEFT JOIN dbo.FactCatalogoServicios fcs WITH (NOLOCK) ON fcs.IdProducto = rd.idItem
			LEFT JOIN dbo.FactPuntosCarga pc WITH (NOLOCK) ON pc.IdPuntoCarga = rc.IdPuntoCarga
			LEFT JOIN dbo.Servicios s WITH (NOLOCK) ON rc.idServicioReceta = s.IdServicio
			LEFT JOIN dbo.Diagnosticos dId WITH (NOLOCK) ON dId.IdDiagnostico = rd.IdDiagnostico
			LEFT JOIN dbo.Diagnosticos dCie WITH (NOLOCK) ON (rd.IdDiagnostico IS NULL OR rd.IdDiagnostico = 0) AND dCie.CodigoCIE10 = rd.CodigoDiagnostico
			WHERE (rc.lnIdRegistroAtencion = @p1 OR (rc.idCuentaAtencion = @p2 AND @p2 > 0))
			  AND (
			      (rc.IdPuntoCarga NOT IN (5, 26) AND rc.IdPuntoCarga IS NOT NULL)
			      OR (rd.Justificacion IS NOT NULL AND LTRIM(RTRIM(CAST(rd.Justificacion AS VARCHAR(MAX)))) <> '')
			      OR fcs.IdServicioGrupo = 1
			  )
		)
		SELECT IdProducto, Codigo, NombreProcedimiento, Servicio, Cantidad, Justificacion, Observacion, Precio, CodigoCIE10, IdDiagnostico
		FROM ItemsApoyoCTE
		WHERE FilaNumero = 1
		ORDER BY IdProducto ASC`

	identificadoresApoyoRegistrados := make(map[int]bool)
	filasApoyo, errorApoyo := r.db.QueryContext(ctx, consultaApoyo,
		sql.Named("p1", identificadorAtencion),
		sql.Named("p2", ficha.IdCuentaAtencion),
	)
	if errorApoyo == nil {
		defer filasApoyo.Close()
		for filasApoyo.Next() {
			var ordenApoyo domain.OrdenApoyoDiagnosticoItem
			var idProd, cant, idDx sql.NullInt64
			var cod, nom, serv, just, obs, cie sql.NullString
			var prec sql.NullFloat64
			if errorEscaneoApoyo := filasApoyo.Scan(
				&idProd,
				&cod,
				&nom,
				&serv,
				&cant,
				&just,
				&obs,
				&prec,
				&cie,
				&idDx,
			); errorEscaneoApoyo == nil {
				ordenApoyo.IdProducto = int(idProd.Int64)
				ordenApoyo.Codigo = strings.TrimSpace(cod.String)
				ordenApoyo.NombreProcedimiento = strings.TrimSpace(nom.String)
				ordenApoyo.Servicio = strings.TrimSpace(serv.String)
				ordenApoyo.Cantidad = int(cant.Int64)
				ordenApoyo.Justificacion = strings.TrimSpace(just.String)
				ordenApoyo.Observacion = strings.TrimSpace(obs.String)
				ordenApoyo.Precio = prec.Float64
				ordenApoyo.CodigoCIE10 = strings.TrimSpace(cie.String)
				ordenApoyo.IdDiagnostico = int(idDx.Int64)

				if ordenApoyo.CodigoCIE10 == "" && len(ficha.Diagnosticos) > 0 {
					ordenApoyo.CodigoCIE10 = ficha.Diagnosticos[0].CodigoCIE10
					ordenApoyo.IdDiagnostico = ficha.Diagnosticos[0].IdDiagnostico
				}

				identificadoresApoyoRegistrados[ordenApoyo.IdProducto] = true
				ficha.OrdenesApoyoDiagnostico = append(ficha.OrdenesApoyoDiagnostico, ordenApoyo)
			}
		}
	}

	ficha.Procedimientos = make([]domain.AtencionProcedimientoCEItem, 0)
	ficha.OrdenesProcedimientos = make([]domain.OrdenProcedimientoCEItem, 0)
	consultaProcedimientos := `
		WITH ItemsProcCTE AS (
			SELECT 
				ISNULL(rd.idItem, 0) AS IdProducto,
				ISNULL(fcs.Codigo, '') AS CodigoCPT,
				ISNULL(NULLIF(fcs.Nombre, ''), ISNULL(rd.DescripcionadicionalReceta, 'PROCEDIMIENTO MÉDICO')) AS NombreProcedimiento,
				ISNULL(rd.CantidadPedida, 1) AS Cantidad,
				ISNULL(rd.observaciones, ISNULL(rd.DescripcionadicionalReceta, '')) AS Observaciones,
				ISNULL(rd.Precio, 0.0) AS Precio,
				COALESCE(NULLIF(LTRIM(RTRIM(rd.CodigoDiagnostico)), ''), dId.CodigoCIE10, dCie.CodigoCIE10, '') AS CodigoCIE10Asociado,
				COALESCE(dId.Descripcion, dCie.Descripcion, '') AS DiagnosticoDescripcionAsociado,
				COALESCE(NULLIF(rd.IdDiagnostico, 0), dId.IdDiagnostico, dCie.IdDiagnostico, 0) AS IdDiagnosticoAsociado,
				ROW_NUMBER() OVER (PARTITION BY rd.idItem ORDER BY rc.idReceta DESC) AS FilaNumero
			FROM dbo.RecetaCabecera rc WITH (NOLOCK)
			INNER JOIN dbo.RecetaDetalle rd WITH (NOLOCK) ON rc.idReceta = rd.idReceta
			LEFT JOIN dbo.FactCatalogoServicios fcs WITH (NOLOCK) ON fcs.IdProducto = rd.idItem
			LEFT JOIN dbo.Diagnosticos dId WITH (NOLOCK) ON dId.IdDiagnostico = rd.IdDiagnostico
			LEFT JOIN dbo.Diagnosticos dCie WITH (NOLOCK) ON (rd.IdDiagnostico IS NULL OR rd.IdDiagnostico = 0) AND dCie.CodigoCIE10 = rd.CodigoDiagnostico
			WHERE (rc.lnIdRegistroAtencion = @p1 OR (rc.idCuentaAtencion = @p2 AND @p2 > 0))
			  AND rc.IdPuntoCarga = 26
			  AND (rd.Justificacion IS NULL OR LTRIM(RTRIM(CAST(rd.Justificacion AS VARCHAR(MAX)))) = '')
			  AND (fcs.IdServicioGrupo IS NULL OR fcs.IdServicioGrupo <> 1)
		)
		SELECT IdProducto, CodigoCPT, NombreProcedimiento, Cantidad, Observaciones, Precio, CodigoCIE10Asociado, DiagnosticoDescripcionAsociado, IdDiagnosticoAsociado
		FROM ItemsProcCTE
		WHERE FilaNumero = 1
		ORDER BY IdProducto ASC`

	filasProcedimientos, errorProcedimientos := r.db.QueryContext(ctx, consultaProcedimientos,
		sql.Named("p1", identificadorAtencion),
		sql.Named("p2", ficha.IdCuentaAtencion),
	)
	if errorProcedimientos == nil {
		defer filasProcedimientos.Close()
		for filasProcedimientos.Next() {
			var procedimiento domain.AtencionProcedimientoCEItem
			var idProd, cant, idDxAsoc sql.NullInt64
			var cpt, nombre, obs, cieAsoc, descAsoc sql.NullString
			var prec sql.NullFloat64
			if errorEscaneoProcedimiento := filasProcedimientos.Scan(
				&idProd,
				&cpt,
				&nombre,
				&cant,
				&obs,
				&prec,
				&cieAsoc,
				&descAsoc,
				&idDxAsoc,
			); errorEscaneoProcedimiento == nil {
				idProductoEntero := int(idProd.Int64)
				if identificadoresApoyoRegistrados[idProductoEntero] {
					continue
				}

				procedimiento.IdProducto = idProductoEntero
				procedimiento.IdProcedimiento = idProductoEntero
				procedimiento.CodigoCPT = strings.TrimSpace(cpt.String)
				procedimiento.NombreProcedimiento = strings.TrimSpace(nombre.String)
				procedimiento.Cantidad = int(cant.Int64)
				procedimiento.Observaciones = strings.TrimSpace(obs.String)
				procedimiento.Precio = prec.Float64
				procedimiento.CodigoCIE10Asociado = strings.TrimSpace(cieAsoc.String)
				procedimiento.DiagnosticoDescripcionAsociado = strings.TrimSpace(descAsoc.String)
				procedimiento.IdDiagnosticoAsociado = int(idDxAsoc.Int64)

				if procedimiento.CodigoCIE10Asociado == "" && len(ficha.Diagnosticos) > 0 {
					procedimiento.CodigoCIE10Asociado = ficha.Diagnosticos[0].CodigoCIE10
					procedimiento.DiagnosticoDescripcionAsociado = ficha.Diagnosticos[0].Descripcion
					procedimiento.IdDiagnosticoAsociado = ficha.Diagnosticos[0].IdDiagnostico
				}

				ficha.Procedimientos = append(ficha.Procedimientos, procedimiento)

				var ordenProcedimiento domain.OrdenProcedimientoCEItem
				ordenProcedimiento.IdProcedimiento = procedimiento.IdProducto
				ordenProcedimiento.IdProducto = procedimiento.IdProducto
				ordenProcedimiento.CodigoCPMS = procedimiento.CodigoCPT
				ordenProcedimiento.NombreProcedimiento = procedimiento.NombreProcedimiento
				ordenProcedimiento.Cantidad = procedimiento.Cantidad
				ordenProcedimiento.Observacion = procedimiento.Observaciones
				ordenProcedimiento.CodigoCIE10 = procedimiento.CodigoCIE10Asociado
				ordenProcedimiento.DiagnosticoDescripcion = procedimiento.DiagnosticoDescripcionAsociado
				ordenProcedimiento.IdDiagnostico = procedimiento.IdDiagnosticoAsociado
				ficha.OrdenesProcedimientos = append(ficha.OrdenesProcedimientos, ordenProcedimiento)
			}
		}
	}
	return &ficha, nil
}

func (r *atencionCERepository) GuardarFichaAtencion(ctx context.Context, payload domain.AtencionCEGuardarPayload) error {
	transaccion, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error al iniciar transaccion para atencion CE: %w", err)
	}
	defer transaccion.Rollback()

	motivo := strings.TrimSpace(payload.Motivo)
	if motivo == "" {
		motivo = "Atención médica en consultorio externo"
	}

	spGuardarAtencion := `
		EXEC [dbo].[AtencionesCEGuardar]
			@idAtencion = @p1,
			@Motivo = @p2,
			@ExamenClinico = @p3,
			@Antecedente = @p4,
			@InforCompDiagMed = @p5,
			@Tratamiento = @p6,
			@PlanTratamiento = @p7,
			@Observaciones = @p8,
			@IdUsuario = @p9
	`
	_, err = transaccion.ExecContext(ctx, spGuardarAtencion,
		sql.Named("p1", payload.IdAtencion),
		sql.Named("p2", motivo),
		sql.Named("p3", strings.TrimSpace(payload.ExamenClinico)),
		sql.Named("p4", strings.TrimSpace(payload.Antecedente)),
		sql.Named("p5", strings.TrimSpace(payload.InforCompDiagMed)),
		sql.Named("p6", strings.TrimSpace(payload.Tratamiento)),
		sql.Named("p7", strings.TrimSpace(payload.PlanTratamiento)),
		sql.Named("p8", strings.TrimSpace(payload.Observaciones)),
		sql.Named("p9", payload.IdUsuario),
	)
	if err != nil {
		return fmt.Errorf("error al ejecutar procedimiento AtencionesCEGuardar: %w", err)
	}

	queryEstadoSP := `EXEC dbo.usp_go_Atenciones_ActualizarEstadoCE 
		@IdAtencion = @p1, 
		@IdDestinoAtencion = @p2, 
		@CerrarAtencion = @p3`
	var exito, identificadorPaciente int
	err = transaccion.QueryRowContext(ctx, queryEstadoSP,
		sql.Named("p1", payload.IdAtencion),
		sql.Named("p2", payload.IdDestinoAtencion),
		sql.Named("p3", payload.CerrarAtencion),
	).Scan(&exito, &identificadorPaciente)
	if err != nil {
		return fmt.Errorf("error al actualizar estado de atencion mediante SP: %w", err)
	}

	if payload.CerrarAtencion {
		actualizarCitaCerrada := `
			UPDATE dbo.Atenciones
			SET idEstadoAtencion = 2,
			    HoraEgreso = CASE WHEN HoraEgreso IS NULL OR RTRIM(HoraEgreso) = '' THEN CONVERT(VARCHAR(5), GETDATE(), 108) ELSE HoraEgreso END
			WHERE IdAtencion = @p1;

			UPDATE dbo.AtencionesCE
			SET Estado = 2
			WHERE idAtencion = @p1;

			UPDATE dbo.Citas
			SET IdEstadoCita = 4
			WHERE IdAtencion = @p1;
		`
		if _, errCita := transaccion.ExecContext(ctx, actualizarCitaCerrada, sql.Named("p1", payload.IdAtencion)); errCita != nil {
			return fmt.Errorf("error al actualizar cita y cierre de atencion CE: %w", errCita)
		}
	}

	if err := guardarDiagnosticosAtencionTx(ctx, transaccion, payload.IdAtencion, identificadorPaciente, payload.Diagnosticos); err != nil {
		return err
	}

	if err := guardarRecetaAtencionTx(ctx, transaccion, payload.IdAtencion, identificadorPaciente, payload.IdUsuario, payload.ItemsReceta); err != nil {
		return err
	}

	if err := guardarInterconsultasAtencionTx(ctx, transaccion, payload.IdAtencion, payload.IdUsuario, payload.Interconsultas); err != nil {
		return err
	}

	if err := guardarProcedimientosYApoyoTx(ctx, transaccion, payload.IdAtencion, identificadorPaciente, payload.IdUsuario, payload.Procedimientos, payload.OrdenesApoyoDiagnostico, payload.OrdenesProcedimientos); err != nil {
		return err
	}

	return transaccion.Commit()
}

func determinarSubclasificacionDx(diagnostico domain.AtencionDiagnosticoGuardarItem) int {
	if diagnostico.Clasificacion != "" {
		switch strings.ToUpper(strings.TrimSpace(diagnostico.Clasificacion)) {
		case "D", "DEFINITIVO":
			return 102
		case "R", "REPETIDO":
			return 103
		case "P", "PRESUNTIVO":
			return 101
		}
	}

	switch diagnostico.IdClasificacionDx {
	case 2:
		return 102
	case 3:
		return 103
	}

	switch diagnostico.IdSubclasificacionDx {
	case 102, 402, 405, 408, 411:
		return 102
	case 103, 403, 406, 412:
		return 103
	default:
		return 101
	}
}

func guardarDiagnosticosAtencionTx(ctx context.Context, tx *sql.Tx, idAtencion int, idPaciente int, diagnosticos []domain.AtencionDiagnosticoGuardarItem) error {
	_, err := tx.ExecContext(ctx, "EXEC dbo.usp_go_AtencionesDiagnosticos_EliminarPorAtencion @IdAtencion = @pIdAtencion",
		sql.Named("pIdAtencion", idAtencion),
	)
	if err != nil {
		return fmt.Errorf("error al limpiar diagnosticos previos mediante SP: %w", err)
	}

	queryDxSP := `
		EXEC [dbo].[usp_go_AtencionesDiagnosticosAgregar]
			@IdAtencionDiagnostico = @p1 OUTPUT,
			@IdClasificacionDx = @p2,
			@IdSubclasificacionDiagnostico = @p3,
			@IdDiagnostico = @p4,
			@IdAtencion = @p5,
			@LabConfHIS = @p6,
			@LabConfHIS_1 = @p7,
			@LabConfHIS_2 = @p8,
			@LabConfHIS_3 = @p9,
			@GrupoHIS = @p10,
			@SubGrupoHIS = @p11,
			@IdEpisodio = @p12,
			@IdEvolucion = @p13,
			@IdPrimeraAtencion = @p14,
			@IdPaciente = @p15
	`

	for _, diagnostico := range diagnosticos {
		if diagnostico.IdDiagnostico <= 0 {
			continue
		}

		identificadorSubclasificacion := determinarSubclasificacionDx(diagnostico)
		identificadorClasificacion := 1

		valorLabHis := strings.TrimSpace(diagnostico.LabConfHIS)
		if valorLabHis == "" {
			valorLabHis = strings.TrimSpace(diagnostico.Clasificacion)
		}
		valorLab1 := strings.TrimSpace(diagnostico.LabConfHIS1)
		valorLab2 := strings.TrimSpace(diagnostico.LabConfHIS2)
		valorLab3 := strings.TrimSpace(diagnostico.LabConfHIS3)

		var identificadorAtencionDiagnosticoGenerado sql.NullString
		_, err = tx.ExecContext(ctx, queryDxSP,
			sql.Named("p1", sql.Out{Dest: &identificadorAtencionDiagnosticoGenerado}),
			sql.Named("p2", identificadorClasificacion),
			sql.Named("p3", identificadorSubclasificacion),
			sql.Named("p4", diagnostico.IdDiagnostico),
			sql.Named("p5", idAtencion),
			sql.Named("p6", valorLabHis),
			sql.Named("p7", valorLab1),
			sql.Named("p8", valorLab2),
			sql.Named("p9", valorLab3),
			sql.Named("p10", 0),
			sql.Named("p11", 0),
			sql.Named("p12", 0),
			sql.Named("p13", 0),
			sql.Named("p14", 0),
			sql.Named("p15", idPaciente),
		)
		if err != nil {
			return fmt.Errorf("error al ejecutar usp_go_AtencionesDiagnosticosAgregar: %w", err)
		}
		cadenaResultado := strings.TrimSpace(identificadorAtencionDiagnosticoGenerado.String)
		if strings.HasPrefix(cadenaResultado, "ERROR;") {
			return fmt.Errorf("error en base de datos al guardar diagnostico %d: %s", diagnostico.IdDiagnostico, strings.TrimPrefix(cadenaResultado, "ERROR;"))
		}
	}

	return nil
}

func guardarRecetaAtencionTx(
	ctx context.Context,
	tx *sql.Tx,
	identificadorAtencion int,
	identificadorPaciente int,
	identificadorUsuario int,
	itemsReceta []domain.ItemRecetaPrescripcion,
) error {
	if len(itemsReceta) == 0 {
		return nil
	}

	var identificadorCuentaAtencion int
	var identificadorServicioIngreso int
	var identificadorMedicoIngreso int

	consultaCuenta := `
		SELECT 
			ISNULL(a.IdCuentaAtencion, 0),
			ISNULL(a.IdServicioIngreso, 0),
			ISNULL(a.IdMedicoIngreso, 0)
		FROM dbo.Atenciones a WITH (NOLOCK)
		WHERE a.IdAtencion = @p1`
	err := tx.QueryRowContext(ctx, consultaCuenta,
		sql.Named("p1", identificadorAtencion),
	).Scan(&identificadorCuentaAtencion, &identificadorServicioIngreso, &identificadorMedicoIngreso)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("error al obtener cuenta de atencion: %w", err)
	}

	identificadorMedico := identificadorMedicoIngreso
	if identificadorMedico <= 0 && identificadorUsuario > 0 {
		var identificadorMedicoEmpleado int
		consultaMedico := `
			SELECT TOP 1 m.IdMedico
			FROM dbo.Medicos m WITH (NOLOCK)
			WHERE m.IdEmpleado = @p1`
		errMedico := tx.QueryRowContext(ctx, consultaMedico,
			sql.Named("p1", identificadorUsuario),
		).Scan(&identificadorMedicoEmpleado)
		if errMedico == nil && identificadorMedicoEmpleado > 0 {
			identificadorMedico = identificadorMedicoEmpleado
		}
	}
	if identificadorMedico <= 0 {
		identificadorMedico = 1
	}

	primerIdentificadorProducto := itemsReceta[0].IdProducto
	if primerIdentificadorProducto <= 0 {
		primerIdentificadorProducto = 1
	}

	if identificadorPaciente <= 0 {
		_ = tx.QueryRowContext(ctx, `
			SELECT ISNULL(IdPaciente, 0)
			FROM dbo.Atenciones WITH (NOLOCK)
			WHERE IdAtencion = @p1`,
			sql.Named("p1", identificadorAtencion),
		).Scan(&identificadorPaciente)
	}

	const puntoCargaFarmacia = 5
	var identificadorReceta int
	_ = tx.QueryRowContext(ctx, `
		SELECT TOP 1 idReceta
		FROM dbo.RecetaCabecera WITH (NOLOCK)
		WHERE (lnIdRegistroAtencion = @p1 OR (idCuentaAtencion = @p2 AND @p2 > 0))
		  AND IdPuntoCarga = 5
		ORDER BY idReceta DESC`,
		sql.Named("p1", identificadorAtencion),
		sql.Named("p2", identificadorCuentaAtencion),
	).Scan(&identificadorReceta)

	if identificadorReceta <= 0 {
		var respuestaCabecera string
		queryCabecera := `
			EXEC usp_go_RecetaCabeceraAgregar
				@Respuesta = @p1 OUTPUT,
				@IdPuntoCarga = @p2,
				@idCuentaAtencion = @p3,
				@idServicioReceta = @p4,
				@idMedicoReceta = @p5,
				@IdProducto = @p6,
				@Idpaciente = @p7,
				@IdUsuarioAuditoria = @p8,
				@IdEvolucion = @p9,
				@IdPrimeraAtencion = @p10`

		_, err = tx.ExecContext(ctx, queryCabecera,
			sql.Named("p1", sql.Out{Dest: &respuestaCabecera}),
			sql.Named("p2", puntoCargaFarmacia),
			sql.Named("p3", identificadorCuentaAtencion),
			sql.Named("p4", identificadorServicioIngreso),
			sql.Named("p5", identificadorMedico),
			sql.Named("p6", primerIdentificadorProducto),
			sql.Named("p7", identificadorPaciente),
			sql.Named("p8", identificadorUsuario),
			sql.Named("p9", 0),
			sql.Named("p10", identificadorAtencion),
		)
		if err != nil {
			return fmt.Errorf("error al crear cabecera de receta en atencion CE: %w", err)
		}

		cadenaLimpia := strings.TrimSpace(respuestaCabecera)
		if strings.HasPrefix(cadenaLimpia, "OK;") {
			fmt.Sscanf(strings.TrimPrefix(cadenaLimpia, "OK;"), "%d", &identificadorReceta)
		}
		if identificadorReceta <= 0 {
			_ = tx.QueryRowContext(ctx, `
				SELECT TOP 1 idReceta
				FROM dbo.RecetaCabecera WITH (NOLOCK)
				WHERE idCuentaAtencion = @p1 AND IdPuntoCarga = 5
				ORDER BY idReceta DESC`,
			sql.Named("p1", identificadorCuentaAtencion),
			).Scan(&identificadorReceta)
		}
	}

	if identificadorReceta <= 0 {
		return fmt.Errorf("no se pudo generar el identificador de receta medica de farmacia")
	}

	_, _ = tx.ExecContext(ctx, "DELETE FROM dbo.RecetaDetalle WHERE idReceta = @p1", sql.Named("p1", identificadorReceta))

	queryDetalle := `
		EXEC usp_go_RecetaDetalleAgregar
			@Mensaje = @p1 OUTPUT,
			@idReceta = @p2,
			@IdProducto = @p3,
			@Cantidad = @p4,
			@Precio = @p5,
			@SaldoEnRegistroReceta = @p6,
			@idDosisRecetada = @p7,
			@observaciones = @p8,
			@IdViaAdministracion = @p9,
			@CodigoDiagnostico = @p10,
			@Justificacion = @p11,
			@idUNIDDosisReceta = @p12,
			@idFrecuencia = @p13,
			@DescripcionadicionalReceta = @p14,
			@Duracion = @p15,
			@idCuentaAtencionProxCita = @p16`

	for _, itemActual := range itemsReceta {
		identificadorProducto := itemActual.IdProducto
		if identificadorProducto <= 0 {
			continue
		}

		var precioUnitario float64
		_ = tx.QueryRowContext(ctx, `
			SELECT TOP 1 ISNULL(PrecioUnitario, 0)
			FROM dbo.Precios WITH (NOLOCK)
			WHERE IdProducto = @p1`,
			sql.Named("p1", identificadorProducto),
		).Scan(&precioUnitario)

		indicacionesAdicionales := fmt.Sprintf("%s | %s | %s | %s | %s",
			strings.TrimSpace(itemActual.Dosis),
			strings.TrimSpace(itemActual.UnidadMedida),
			strings.TrimSpace(itemActual.Frecuencia),
			strings.TrimSpace(itemActual.Via),
			strings.TrimSpace(itemActual.Duracion),
		)
		idVia := resolverIdentificadorViaAdministracion(itemActual.Via)
		idFrecuencia := resolverIdentificadorFrecuencia(itemActual.Frecuencia)
		duracionDias := parsearDuracionDias(itemActual.Duracion)
		var mensajeDetalle string
		_, err = tx.ExecContext(ctx, queryDetalle,
			sql.Named("p1", sql.Out{Dest: &mensajeDetalle}),
			sql.Named("p2", identificadorReceta),
			sql.Named("p3", identificadorProducto),
			sql.Named("p4", itemActual.Cantidad),
			sql.Named("p5", precioUnitario),
			sql.Named("p6", nil),
			sql.Named("p7", nil),
			sql.Named("p8", strings.TrimSpace(itemActual.Indicaciones)),
			sql.Named("p9", idVia),
			sql.Named("p10", strings.TrimSpace(itemActual.CodigoCIE10)),
			sql.Named("p11", nil),
			sql.Named("p12", nil),
			sql.Named("p13", idFrecuencia),
			sql.Named("p14", indicacionesAdicionales),
			sql.Named("p15", duracionDias),
			sql.Named("p16", nil),
		)
		if err != nil {
			return fmt.Errorf("error al agregar producto %d a receta de atencion: %w", identificadorProducto, err)
		}
		cadenaMensaje := strings.TrimSpace(mensajeDetalle)
		if !strings.HasPrefix(cadenaMensaje, "OK") && !strings.Contains(cadenaMensaje, "ya Existe") {
			return fmt.Errorf("el sistema rechazó el producto %d: %s", identificadorProducto, cadenaMensaje)
		}
	}

	return nil
}

func parsearDuracionDias(duracion string) *int {
	texto := strings.TrimSpace(duracion)
	if texto == "" {
		return nil
	}
	var dias int
	for _, caracter := range texto {
		if caracter >= '0' && caracter <= '9' {
			dias = dias*10 + int(caracter-'0')
		} else if dias > 0 {
			break
		}
	}
	if dias > 0 {
		return &dias
	}
	return nil
}

func resolverIdentificadorViaAdministracion(via string) *int {
	viaNormalizada := strings.ToUpper(strings.TrimSpace(via))
	var identificador int
	switch {
	case strings.Contains(viaNormalizada, "ORAL"):
		identificador = 1
	case strings.Contains(viaNormalizada, "SUBLINGUAL"):
		identificador = 2
	case strings.Contains(viaNormalizada, "RECTAL"):
		identificador = 3
	case strings.Contains(viaNormalizada, "INTRAVENOSA") || strings.Contains(viaNormalizada, "EV") || strings.Contains(viaNormalizada, "IV"):
		identificador = 4
	case strings.Contains(viaNormalizada, "INTRAMUSCULAR") || strings.Contains(viaNormalizada, "IM"):
		identificador = 5
	case strings.Contains(viaNormalizada, "SUBCUT"):
		identificador = 6
	case strings.Contains(viaNormalizada, "DERM") || strings.Contains(viaNormalizada, "TOPIC"):
		identificador = 7
	case strings.Contains(viaNormalizada, "NASAL"):
		identificador = 8
	case strings.Contains(viaNormalizada, "OFTAL") || strings.Contains(viaNormalizada, "OCUL"):
		identificador = 9
	default:
		identificador = 1
	}
	return &identificador
}

func resolverIdentificadorFrecuencia(frecuencia string) *int {
	frecuenciaNormalizada := strings.ToUpper(strings.TrimSpace(frecuencia))
	var identificador int
	switch {
	case strings.Contains(frecuenciaNormalizada, "2 HORAS") || strings.Contains(frecuenciaNormalizada, "2H"):
		identificador = 2
	case strings.Contains(frecuenciaNormalizada, "4 HORAS") || strings.Contains(frecuenciaNormalizada, "4H"):
		identificador = 4
	case strings.Contains(frecuenciaNormalizada, "6 HORAS") || strings.Contains(frecuenciaNormalizada, "6H"):
		identificador = 6
	case strings.Contains(frecuenciaNormalizada, "8 HORAS") || strings.Contains(frecuenciaNormalizada, "8H"):
		identificador = 8
	case strings.Contains(frecuenciaNormalizada, "12 HORAS") || strings.Contains(frecuenciaNormalizada, "12H"):
		identificador = 12
	case strings.Contains(frecuenciaNormalizada, "24 HORAS") || strings.Contains(frecuenciaNormalizada, "24H") || strings.Contains(frecuenciaNormalizada, "DIARIA"):
		identificador = 24
	default:
		identificador = 24
	}
	return &identificador
}

func guardarInterconsultasAtencionTx(
	ctx context.Context,
	transaccion *sql.Tx,
	identificadorAtencion int,
	identificadorUsuario int,
	interconsultas []domain.InterconsultaCEItem,
) error {
	if len(interconsultas) == 0 {
		return nil
	}

	for _, interconsulta := range interconsultas {
		if interconsulta.IdEspecialidad <= 0 {
			continue
		}

		var identificadorExistente int
		_ = transaccion.QueryRowContext(ctx, `
			SELECT TOP 1 idInterconsulta 
			FROM dbo.Interconsulta WITH (NOLOCK) 
			WHERE IdRegistroAtencion = @p1 
			  AND IdEspecialidadSolicitada = @p2`,
			sql.Named("p1", identificadorAtencion),
			sql.Named("p2", interconsulta.IdEspecialidad),
		).Scan(&identificadorExistente)

		if identificadorExistente == 0 {
			_, errorInterconsulta := transaccion.ExecContext(ctx, `
				EXEC dbo.usp_go_InterconsultaCrear
					@IdAtencionOrigen = @p1,
					@IdEspecialidad = @p2,
					@IdMedicoDestino = @p3,
					@Motivo = @p4,
					@IdPrioridad = 1,
					@IdEmpleadoSolicitante = @p5,
					@ResumenHC = @p6`,
				sql.Named("p1", identificadorAtencion),
				sql.Named("p2", interconsulta.IdEspecialidad),
				sql.Named("p3", interconsulta.IdMedico),
				sql.Named("p4", strings.TrimSpace(interconsulta.Motivo)),
				sql.Named("p5", identificadorUsuario),
				sql.Named("p6", strings.TrimSpace(interconsulta.ResumenEnfermedad)),
			)
			if errorInterconsulta != nil {
				return fmt.Errorf("error al registrar interconsulta para especialidad %d: %w", interconsulta.IdEspecialidad, errorInterconsulta)
			}
		} else {
			_, _ = transaccion.ExecContext(ctx, `
				UPDATE dbo.Interconsulta
				SET MotivoClinico = @p1,
				    ResumenClinico = @p2,
				    ResumenHC = @p2
				WHERE idInterconsulta = @p3`,
				sql.Named("p1", strings.TrimSpace(interconsulta.Motivo)),
				sql.Named("p2", strings.TrimSpace(interconsulta.ResumenEnfermedad)),
				sql.Named("p3", identificadorExistente),
			)
		}
	}
	return nil
}

func guardarProcedimientosYApoyoTx(
	ctx context.Context,
	transaccion *sql.Tx,
	identificadorAtencion int,
	identificadorPaciente int,
	identificadorUsuario int,
	procedimientos []domain.AtencionProcedimientoCEItem,
	ordenesApoyo []domain.OrdenApoyoDiagnosticoItem,
	ordenesProcedimientos []domain.OrdenProcedimientoCEItem,
) error {
	totalProcedimientos := len(procedimientos) + len(ordenesProcedimientos)
	totalApoyo := len(ordenesApoyo)
	if totalProcedimientos == 0 && totalApoyo == 0 {
		return nil
	}

	var identificadorCuentaAtencion int
	var identificadorServicioIngreso int
	var identificadorMedicoIngreso int

	_ = transaccion.QueryRowContext(ctx, `
		SELECT 
			ISNULL(a.IdCuentaAtencion, 0),
			ISNULL(a.IdServicioIngreso, 0),
			ISNULL(a.IdMedicoIngreso, 0)
		FROM dbo.Atenciones a WITH (NOLOCK)
		WHERE a.IdAtencion = @p1`,
		sql.Named("p1", identificadorAtencion),
	).Scan(&identificadorCuentaAtencion, &identificadorServicioIngreso, &identificadorMedicoIngreso)

	identificadorMedico := identificadorMedicoIngreso
	if identificadorMedico <= 0 && identificadorUsuario > 0 {
		var identificadorMedicoEmpleado int
		_ = transaccion.QueryRowContext(ctx, `
			SELECT TOP 1 m.IdMedico
			FROM dbo.Medicos m WITH (NOLOCK)
			WHERE m.IdEmpleado = @p1`,
			sql.Named("p1", identificadorUsuario),
		).Scan(&identificadorMedicoEmpleado)
		if identificadorMedicoEmpleado > 0 {
			identificadorMedico = identificadorMedicoEmpleado
		}
	}
	if identificadorMedico <= 0 {
		identificadorMedico = 1
	}

	if identificadorPaciente <= 0 {
		_ = transaccion.QueryRowContext(ctx, `
			SELECT ISNULL(IdPaciente, 0)
			FROM dbo.Atenciones WITH (NOLOCK)
			WHERE IdAtencion = @p1`,
			sql.Named("p1", identificadorAtencion),
		).Scan(&identificadorPaciente)
	}

	consultaDetalle := `
		EXEC usp_go_RecetaDetalleAgregar
			@Mensaje = @p1 OUTPUT,
			@idReceta = @p2,
			@IdProducto = @p3,
			@Cantidad = @p4,
			@Precio = @p5,
			@SaldoEnRegistroReceta = @p6,
			@idDosisRecetada = @p7,
			@observaciones = @p8,
			@IdViaAdministracion = @p9,
			@CodigoDiagnostico = @p10,
			@Justificacion = @p11,
			@idUNIDDosisReceta = @p12,
			@idFrecuencia = @p13,
			@DescripcionadicionalReceta = @p14,
			@Duracion = @p15,
			@idCuentaAtencionProxCita = @p16`

	// 1. Guardar Procedimientos (Punto de Carga = 26)
	if totalProcedimientos > 0 {
		const puntoCargaProcedimientos = 26
		var identificadorRecetaProc int
		_ = transaccion.QueryRowContext(ctx, `
			SELECT TOP 1 idReceta
			FROM dbo.RecetaCabecera WITH (NOLOCK)
			WHERE (lnIdRegistroAtencion = @p1 OR (idCuentaAtencion = @p2 AND @p2 > 0))
			  AND IdPuntoCarga = @p3
			ORDER BY idReceta DESC`,
			sql.Named("p1", identificadorAtencion),
			sql.Named("p2", identificadorCuentaAtencion),
			sql.Named("p3", puntoCargaProcedimientos),
		).Scan(&identificadorRecetaProc)

		if identificadorRecetaProc <= 0 {
			primerProducto := 1
			if len(procedimientos) > 0 && procedimientos[0].IdProducto > 0 {
				primerProducto = procedimientos[0].IdProducto
			} else if len(ordenesProcedimientos) > 0 && ordenesProcedimientos[0].IdProducto > 0 {
				primerProducto = ordenesProcedimientos[0].IdProducto
			}

			var respuestaCabecera string
			consultaCabecera := `
				EXEC usp_go_RecetaCabeceraAgregar
					@Respuesta = @p1 OUTPUT,
					@IdPuntoCarga = @p2,
					@idCuentaAtencion = @p3,
					@idServicioReceta = @p4,
					@idMedicoReceta = @p5,
					@IdProducto = @p6,
					@Idpaciente = @p7,
					@IdUsuarioAuditoria = @p8,
					@IdEvolucion = @p9,
					@IdPrimeraAtencion = @p10`

			_ = transaccion.QueryRowContext(ctx, consultaCabecera,
				sql.Named("p1", sql.Out{Dest: &respuestaCabecera}),
				sql.Named("p2", puntoCargaProcedimientos),
				sql.Named("p3", identificadorCuentaAtencion),
				sql.Named("p4", identificadorServicioIngreso),
				sql.Named("p5", identificadorMedico),
				sql.Named("p6", primerProducto),
				sql.Named("p7", identificadorPaciente),
				sql.Named("p8", identificadorUsuario),
				sql.Named("p9", 0),
				sql.Named("p10", identificadorAtencion),
			).Scan()

			cadenaLimpia := strings.TrimSpace(respuestaCabecera)
		if strings.HasPrefix(cadenaLimpia, "OK;") {
				fmt.Sscanf(strings.TrimPrefix(cadenaLimpia, "OK;"), "%d", &identificadorRecetaProc)
			}
		}

		if identificadorRecetaProc > 0 {
			_, _ = transaccion.ExecContext(ctx, "DELETE FROM dbo.RecetaDetalle WHERE idReceta = @p1", sql.Named("p1", identificadorRecetaProc))
			productosVistosProc := make(map[int]bool)
			for _, procedimiento := range procedimientos {
				if procedimiento.IdProducto <= 0 || productosVistosProc[procedimiento.IdProducto] {
					continue
				}
				productosVistosProc[procedimiento.IdProducto] = true
				if procedimiento.IdProducto <= 0 {
					continue
				}
				var mensajeProcedimiento string
				_, errorProcedimiento := transaccion.ExecContext(ctx, consultaDetalle,
					sql.Named("p1", sql.Out{Dest: &mensajeProcedimiento}),
					sql.Named("p2", identificadorRecetaProc),
					sql.Named("p3", procedimiento.IdProducto),
					sql.Named("p4", procedimiento.Cantidad),
					sql.Named("p5", procedimiento.Precio),
					sql.Named("p6", nil),
					sql.Named("p7", nil),
					sql.Named("p8", strings.TrimSpace(procedimiento.Observaciones)),
					sql.Named("p9", nil),
					sql.Named("p10", strings.TrimSpace(procedimiento.CodigoCIE10Asociado)),
					sql.Named("p11", nil),
					sql.Named("p12", nil),
					sql.Named("p13", nil),
					sql.Named("p14", strings.TrimSpace(procedimiento.NombreProcedimiento)),
					sql.Named("p15", nil),
					sql.Named("p16", nil),
				)
				if errorProcedimiento != nil {
					return fmt.Errorf("error al registrar procedimiento %s en receta: %w", procedimiento.NombreProcedimiento, errorProcedimiento)
				}
			}

			for _, ordenProcedimiento := range ordenesProcedimientos {
				if ordenProcedimiento.IdProducto <= 0 || productosVistosProc[ordenProcedimiento.IdProducto] {
					continue
				}
				productosVistosProc[ordenProcedimiento.IdProducto] = true
				var mensajeOrdenProcedimiento string
				_, errorOrdenProcedimiento := transaccion.ExecContext(ctx, consultaDetalle,
					sql.Named("p1", sql.Out{Dest: &mensajeOrdenProcedimiento}),
					sql.Named("p2", identificadorRecetaProc),
					sql.Named("p3", ordenProcedimiento.IdProducto),
					sql.Named("p4", ordenProcedimiento.Cantidad),
					sql.Named("p5", 0.0),
					sql.Named("p6", nil),
					sql.Named("p7", nil),
					sql.Named("p8", strings.TrimSpace(ordenProcedimiento.Observacion)),
					sql.Named("p9", nil),
					sql.Named("p10", strings.TrimSpace(ordenProcedimiento.CodigoCIE10)),
					sql.Named("p11", nil),
					sql.Named("p12", nil),
					sql.Named("p13", nil),
					sql.Named("p14", strings.TrimSpace(ordenProcedimiento.NombreProcedimiento)),
					sql.Named("p15", nil),
					sql.Named("p16", nil),
				)
				if errorOrdenProcedimiento != nil {
					return fmt.Errorf("error al registrar orden de procedimiento %s en receta: %w", ordenProcedimiento.NombreProcedimiento, errorOrdenProcedimiento)
				}
			}
		}
	}

	// 2. Guardar Apoyo al Diagnóstico (Punto de Carga = 1)
	if totalApoyo > 0 {
		const puntoCargaApoyo = 1
		var identificadorRecetaApoyo int
		_ = transaccion.QueryRowContext(ctx, `
			SELECT TOP 1 idReceta
			FROM dbo.RecetaCabecera WITH (NOLOCK)
			WHERE (lnIdRegistroAtencion = @p1 OR (idCuentaAtencion = @p2 AND @p2 > 0))
			  AND IdPuntoCarga = @p3
			ORDER BY idReceta DESC`,
			sql.Named("p1", identificadorAtencion),
			sql.Named("p2", identificadorCuentaAtencion),
			sql.Named("p3", puntoCargaApoyo),
		).Scan(&identificadorRecetaApoyo)

		if identificadorRecetaApoyo <= 0 {
			primerProducto := ordenesApoyo[0].IdProducto
			if primerProducto <= 0 {
				primerProducto = 1
			}

			var respuestaCabecera string
			consultaCabecera := `
				EXEC usp_go_RecetaCabeceraAgregar
					@Respuesta = @p1 OUTPUT,
					@IdPuntoCarga = @p2,
					@idCuentaAtencion = @p3,
					@idServicioReceta = @p4,
					@idMedicoReceta = @p5,
					@IdProducto = @p6,
					@Idpaciente = @p7,
					@IdUsuarioAuditoria = @p8,
					@IdEvolucion = @p9,
					@IdPrimeraAtencion = @p10`

			_ = transaccion.QueryRowContext(ctx, consultaCabecera,
				sql.Named("p1", sql.Out{Dest: &respuestaCabecera}),
				sql.Named("p2", puntoCargaApoyo),
				sql.Named("p3", identificadorCuentaAtencion),
				sql.Named("p4", identificadorServicioIngreso),
				sql.Named("p5", identificadorMedico),
				sql.Named("p6", primerProducto),
				sql.Named("p7", identificadorPaciente),
				sql.Named("p8", identificadorUsuario),
				sql.Named("p9", 0),
				sql.Named("p10", identificadorAtencion),
			).Scan()

			cadenaLimpia := strings.TrimSpace(respuestaCabecera)
			if strings.HasPrefix(cadenaLimpia, "OK;") {
				fmt.Sscanf(strings.TrimPrefix(cadenaLimpia, "OK;"), "%d", &identificadorRecetaApoyo)
			}
		}

		if identificadorRecetaApoyo > 0 {
			_, _ = transaccion.ExecContext(ctx, "DELETE FROM dbo.RecetaDetalle WHERE idReceta = @p1", sql.Named("p1", identificadorRecetaApoyo))
			for _, apoyo := range ordenesApoyo {
				if apoyo.IdProducto <= 0 {
					continue
				}
				justificacionTexto := strings.TrimSpace(apoyo.Justificacion)
				if justificacionTexto == "" {
					justificacionTexto = "SOLICITUD DE APOYO AL DIAGNÓSTICO"
				}
				var mensajeApoyo string
				_, errorApoyo := transaccion.ExecContext(ctx, consultaDetalle,
					sql.Named("p1", sql.Out{Dest: &mensajeApoyo}),
					sql.Named("p2", identificadorRecetaApoyo),
					sql.Named("p3", apoyo.IdProducto),
					sql.Named("p4", apoyo.Cantidad),
					sql.Named("p5", apoyo.Precio),
					sql.Named("p6", nil),
					sql.Named("p7", nil),
					sql.Named("p8", strings.TrimSpace(apoyo.Observacion)),
					sql.Named("p9", nil),
					sql.Named("p10", strings.TrimSpace(apoyo.CodigoCIE10)),
					sql.Named("p11", justificacionTexto),
					sql.Named("p12", nil),
					sql.Named("p13", nil),
					sql.Named("p14", strings.TrimSpace(apoyo.NombreProcedimiento)),
					sql.Named("p15", nil),
					sql.Named("p16", nil),
				)
				if errorApoyo != nil {
					return fmt.Errorf("error al registrar apoyo diagnostico %s en receta: %w", apoyo.NombreProcedimiento, errorApoyo)
				}
			}
		}
	}

	return nil
}

func (r *atencionCERepository) ListarDestinosAtencion(ctx context.Context) ([]domain.DestinoAtencionItem, error) {
	query := "EXEC dbo.usp_go_TiposDestinoAtencion_ListarCE"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al listar destinos de atencion CE: %w", err)
	}
	defer rows.Close()

	destinos := make([]domain.DestinoAtencionItem, 0)
	for rows.Next() {
		var item domain.DestinoAtencionItem
		if err := rows.Scan(&item.IdDestinoAtencion, &item.Descripcion); err != nil {
			return nil, fmt.Errorf("error al escanear destino: %w", err)
		}
		destinos = append(destinos, item)
	}

	return destinos, rows.Err()
}

func (r *atencionCERepository) ListarCatalogoLabsHIS(ctx context.Context) ([]domain.CatalogoLabHISItem, error) {
	query := "EXEC dbo.usp_go_MaestroHisLab_Listar"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error al consultar catalogo de labs HIS: %w", err)
	}
	defer rows.Close()

	items := make([]domain.CatalogoLabHISItem, 0)
	for rows.Next() {
		var item domain.CatalogoLabHISItem
		if err := rows.Scan(&item.Codigo, &item.Valor, &item.Descripcion); err != nil {
			return nil, fmt.Errorf("error al escanear lab HIS: %w", err)
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *atencionCERepository) ListarHistorialAtenciones(ctx context.Context, idPaciente int) ([]domain.HistorialAtencionItem, error) {
	query := `
		SELECT 
			a.IdAtencion,
			ISNULL(a.IdCuentaAtencion, 0) AS IdCuentaAtencion,
			ISNULL(a.IdPaciente, 0) AS IdPaciente,
			CONVERT(VARCHAR(10), a.FechaIngreso, 103) AS FechaIngreso,
			ISNULL(a.HoraIngreso, '') AS HoraIngreso,
			ISNULL(a.HoraEgreso, '') AS HoraEgreso,
			LTRIM(RTRIM(ISNULL(s.Nombre, 'CONSULTA EXTERNA'))) AS ServicioNombre,
			ISNULL(a.IdMedicoIngreso, 0) AS IdMedico,
			LTRIM(RTRIM(CONCAT(
				ISNULL(e.ApellidoPaterno, ''), ' ', 
				ISNULL(e.ApellidoMaterno, ''), ' ', 
				ISNULL(e.Nombres, '')
			))) AS MedicoNombre,
			ISNULL(a.IdTipoServicio, 1) AS IdTipoServicio,
			CASE 
				WHEN a.IdTipoServicio = 1 THEN 'CONSULTORIOS EXTERNOS'
				ELSE 'EMERGENCIA'
			END AS TipoServicioNombre,
			15 AS TiempoProgramado,
			ISNULL(DATEDIFF(MINUTE, 
				TRY_CAST(a.HoraIngreso AS TIME), 
				TRY_CAST(a.HoraEgreso AS TIME)
			), 15) AS TiempoAtencion
		FROM dbo.Atenciones a WITH (NOLOCK)
		LEFT JOIN dbo.Servicios s WITH (NOLOCK) ON a.IdServicioIngreso = s.IdServicio
		LEFT JOIN dbo.Medicos m WITH (NOLOCK) ON a.IdMedicoIngreso = m.IdMedico
		LEFT JOIN dbo.Empleados e WITH (NOLOCK) ON m.IdEmpleado = e.IdEmpleado
		WHERE a.IdPaciente = @p1
		ORDER BY a.FechaIngreso DESC, a.IdAtencion DESC`

	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idPaciente))
	if err != nil {
		return nil, fmt.Errorf("error al listar historial de atenciones del paciente: %w", err)
	}
	defer rows.Close()

	historial := make([]domain.HistorialAtencionItem, 0)
	for rows.Next() {
		var item domain.HistorialAtencionItem
		var ta sql.NullInt64
		if err := rows.Scan(
			&item.IdAtencion,
			&item.IdCuentaAtencion,
			&item.IdPaciente,
			&item.FechaIngreso,
			&item.HoraIngreso,
			&item.HoraEgreso,
			&item.ServicioNombre,
			&item.IdMedico,
			&item.MedicoNombre,
			&item.IdTipoServicio,
			&item.TipoServicioNombre,
			&item.TiempoProgramado,
			&ta,
		); err != nil {
			return nil, fmt.Errorf("error al escanear historial de atencion: %w", err)
		}
		if ta.Valid && ta.Int64 > 0 {
			item.TiempoAtencion = int(ta.Int64)
		} else {
			item.TiempoAtencion = 15
		}
		historial = append(historial, item)
	}

	return historial, rows.Err()
}

func (r *atencionCERepository) ObtenerHistorialAtencionDetalle(ctx context.Context, idAtencion int) (*domain.HistorialAtencionDetalle, error) {
	queryCab := `
		SELECT 
			a.IdAtencion,
			ISNULL(a.IdCuentaAtencion, 0) AS IdCuentaAtencion,
			ISNULL(a.IdPaciente, 0) AS IdPaciente,
			LTRIM(RTRIM(CONCAT(
				ISNULL(p.ApellidoPaterno, ''), ' ', 
				ISNULL(p.ApellidoMaterno, ''), ' ', 
				ISNULL(p.PrimerNombre, ''), ' ', 
				ISNULL(p.SegundoNombre, '')
			))) AS PacienteNombre,
			ISNULL(p.NroHistoriaClinica, 0) AS NroHistoriaClinica,
			LTRIM(RTRIM(ISNULL(p.NroDocumento, ''))) AS NroDocumento,
			ISNULL(a.Edad, 0) AS Edad,
			CASE WHEN p.IdTipoSexo = 2 THEN 'Femenino' ELSE 'Masculino' END AS Sexo,
			LTRIM(RTRIM(ISNULL(s.Nombre, 'CONSULTA EXTERNA'))) AS ServicioNombre,
			LTRIM(RTRIM(CONCAT(
				ISNULL(e.ApellidoPaterno, ''), ' ', 
				ISNULL(e.ApellidoMaterno, ''), ' ', 
				ISNULL(e.Nombres, '')
			))) AS MedicoNombre,
			CONVERT(VARCHAR(10), a.FechaIngreso, 103) AS FechaAtencion,
			ISNULL(a.HoraIngreso, '') AS HoraInicio,
			ISNULL(a.HoraEgreso, '') AS HoraFin,
			ISNULL(ace.InforCompDiagMed, '') AS InforCompDiagMed,
			ISNULL(ace.Motivo, '') AS Motivo,
			ISNULL(ace.Tratamiento, '') AS Tratamiento,
			ISNULL(ace.PlanTratamiento, '') AS PlanTratamiento,
			ISNULL(ace.Observaciones, '') AS Observaciones,
			ISNULL(ace.ExamenClinico, '') AS ExamenClinico
		FROM dbo.Atenciones a WITH (NOLOCK)
		LEFT JOIN dbo.Pacientes p WITH (NOLOCK) ON a.IdPaciente = p.IdPaciente
		LEFT JOIN dbo.Servicios s WITH (NOLOCK) ON a.IdServicioIngreso = s.IdServicio
		LEFT JOIN dbo.Medicos m WITH (NOLOCK) ON a.IdMedicoIngreso = m.IdMedico
		LEFT JOIN dbo.Empleados e WITH (NOLOCK) ON m.IdEmpleado = e.IdEmpleado
		LEFT JOIN dbo.AtencionesCE ace WITH (NOLOCK) ON a.IdAtencion = ace.idAtencion
		WHERE a.IdAtencion = @p1`

	var d domain.HistorialAtencionDetalle
	err := r.db.QueryRowContext(ctx, queryCab, sql.Named("p1", idAtencion)).Scan(
		&d.IdAtencion,
		&d.IdCuentaAtencion,
		&d.IdPaciente,
		&d.PacienteNombre,
		&d.NroHistoriaClinica,
		&d.NroDocumento,
		&d.Edad,
		&d.Sexo,
		&d.ServicioNombre,
		&d.MedicoNombre,
		&d.FechaAtencion,
		&d.HoraInicio,
		&d.HoraFin,
		&d.InforCompDiagMed,
		&d.Motivo,
		&d.Tratamiento,
		&d.PlanTratamiento,
		&d.Observaciones,
		&d.ExamenClinico,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("error al obtener cabecera de detalle de atencion: %w", err)
	}

	queryTriaje := `
		SELECT TOP 1
			ISNULL(PresionArterial, '')        AS PresionArterial,
			ISNULL(TRY_CAST(FrecCardiaca AS INT), 0)         AS FrecuenciaCardiaca,
			ISNULL(TRY_CAST(FrecRespiratoria AS INT), 0)     AS FrecuenciaRespiratoria,
			ISNULL(TRY_CAST(Temperatura AS DECIMAL(5,2)), 0.0) AS Temperatura,
			ISNULL(TRY_CAST(SAT02 AS INT), 0)                AS Saturacion,
			ISNULL(TRY_CAST(Peso AS DECIMAL(5,2)), 0.0)      AS Peso,
			ISNULL(TRY_CAST(Talla AS DECIMAL(5,2)), 0.0)     AS Talla,
			ISNULL(TRY_CAST(IMC AS DECIMAL(5,2)), 0.0)       AS IMC,
			CONVERT(VARCHAR(10), Fecha, 103)   AS FechaRegistro,
			ISNULL(Observacion, '')            AS Observacion
		FROM dbo.AtencionesTriaje WITH (NOLOCK)
		WHERE IdAtencion = @p1
		ORDER BY IdTriaje DESC`

	var t domain.HistorialTriajeItem
	errT := r.db.QueryRowContext(ctx, queryTriaje, sql.Named("p1", idAtencion)).Scan(
		&t.PresionArterial,
		&t.FrecuenciaCardiaca,
		&t.FrecuenciaRespiratoria,
		&t.Temperatura,
		&t.Saturacion,
		&t.Peso,
		&t.Talla,
		&t.IMC,
		&t.FechaRegistro,
		&t.Observacion,
	)
	if errT == nil {
		d.Triaje = &t
	}

	queryDx := `EXEC [dbo].[usp_go_ObtenerDiagnosticosAtencion] @IdAtencion = @p1, @IdPrimeraAtencion = 0, @IdEvolucion = 0`
	rowsDx, errDx := r.db.QueryContext(ctx, queryDx, sql.Named("p1", idAtencion))
	d.Diagnosticos = make([]domain.AtencionDiagnosticoCEItem, 0)
	if errDx == nil {
		defer rowsDx.Close()
		for rowsDx.Next() {
			var dx domain.AtencionDiagnosticoCEItem
			var idAtDummy int
			var descLarga string
			var tipoCodigo string
			var tipoDx string
			var idSubBD int
			if err := rowsDx.Scan(
				&dx.IdAtencionDiagnostico,
				&idAtDummy,
				&dx.IdDiagnostico,
				&dx.CodigoCIE10,
				&dx.Descripcion,
				&descLarga,
				&tipoCodigo,
				&tipoDx,
				&idSubBD,
				&dx.LabConfHIS,
				&dx.LabConfHIS1,
				&dx.LabConfHIS2,
				&dx.LabConfHIS3,
			); err == nil {
				dx.Clasificacion = strings.ToUpper(strings.TrimSpace(tipoCodigo))
				d.Diagnosticos = append(d.Diagnosticos, dx)
			}
		}
	}

	return &d, nil
}

func (r *atencionCERepository) ListarReportesOperatorios(ctx context.Context, identificadorPaciente int) ([]domain.HistorialReporteOperatorioItem, error) {
	consulta := `
		SELECT 
			a.IdAtencion AS IdReporte,
			a.IdAtencion,
			ISNULL(s.Codigo, '56054') AS CodigoCPT,
			CONVERT(VARCHAR(10), a.FechaIngreso, 103) AS Fecha,
			ISNULL(s.Nombre, 'INTERVENCIÓN QUIRÚRGICA') AS Descripcion,
			ISNULL(LTRIM(RTRIM(CONCAT(e.ApellidoPaterno, ' ', e.ApellidoMaterno, ' ', e.Nombres))), 'CIRUJANO PRINCIPAL') AS Cirujano,
			ISNULL(s.Nombre, 'CIRUGÍA GENERAL') AS Servicio,
			'VALIDADO' AS Estado,
			ISNULL(a.HoraIngreso, '08:00') AS HoraInicio,
			ISNULL(a.HoraEgreso, '09:30') AS HoraFin,
			ISNULL(CONVERT(VARCHAR(10), a.IdCuentaAtencion), '') AS NumeroCuenta
		FROM dbo.Atenciones a WITH (NOLOCK)
		INNER JOIN dbo.Servicios s WITH (NOLOCK) ON a.IdServicioIngreso = s.IdServicio
		LEFT JOIN dbo.Medicos m WITH (NOLOCK) ON a.IdMedicoIngreso = m.IdMedico
		LEFT JOIN dbo.Empleados e WITH (NOLOCK) ON m.IdEmpleado = e.IdEmpleado
		WHERE a.IdPaciente = @p1 
		  AND (s.EsQuirurgico = 1 OR s.Nombre LIKE '%CIRUGIA%' OR s.Nombre LIKE '%QUIRURG%' OR s.Nombre LIKE '%OPERAT%' OR a.IdTipoServicio = 3)

		UNION ALL

		SELECT 
			ISNULL(rd.idItem, a.IdAtencion) AS IdReporte,
			a.IdAtencion,
			ISNULL(fcs.Codigo, '56054') AS CodigoCPT,
			CONVERT(VARCHAR(10), ISNULL(rc.FechaReceta, a.FechaIngreso), 103) AS Fecha,
			ISNULL(fcs.Nombre, 'PROCEDIMIENTO QUIRÚRGICO') AS Descripcion,
			ISNULL(LTRIM(RTRIM(CONCAT(e.ApellidoPaterno, ' ', e.ApellidoMaterno, ' ', e.Nombres))), 'CIRUJANO') AS Cirujano,
			ISNULL(s.Nombre, 'CENTRO QUIRÚRGICO') AS Servicio,
			'VALIDADO' AS Estado,
			ISNULL(a.HoraIngreso, '08:00') AS HoraInicio,
			ISNULL(a.HoraEgreso, '09:30') AS HoraFin,
			ISNULL(CONVERT(VARCHAR(10), a.IdCuentaAtencion), '') AS NumeroCuenta
		FROM dbo.Atenciones a WITH (NOLOCK)
		INNER JOIN dbo.RecetaCabecera rc WITH (NOLOCK) ON rc.idCuentaAtencion = a.IdCuentaAtencion
		INNER JOIN dbo.RecetaDetalle rd WITH (NOLOCK) ON rd.idReceta = rc.idReceta
		INNER JOIN dbo.FactCatalogoServicios fcs WITH (NOLOCK) ON fcs.IdServicio = rd.idItem
		LEFT JOIN dbo.Servicios s WITH (NOLOCK) ON a.IdServicioIngreso = s.IdServicio
		LEFT JOIN dbo.Medicos m WITH (NOLOCK) ON a.IdMedicoIngreso = m.IdMedico
		LEFT JOIN dbo.Empleados e WITH (NOLOCK) ON m.IdEmpleado = e.IdEmpleado
		WHERE a.IdPaciente = @p1 AND fcs.EsQuirurgico = 1
		ORDER BY Fecha DESC`

	filas, errorConsulta := r.db.QueryContext(ctx, consulta, sql.Named("p1", identificadorPaciente))
	if errorConsulta != nil {
		return []domain.HistorialReporteOperatorioItem{}, nil
	}
	defer filas.Close()

	elementos := make([]domain.HistorialReporteOperatorioItem, 0)
	for filas.Next() {
		var item domain.HistorialReporteOperatorioItem
		var horaInicio, horaFin, numeroCuenta string
		if errorEscaneo := filas.Scan(
			&item.IdReporte,
			&item.IdAtencion,
			&item.CodigoCPT,
			&item.Fecha,
			&item.Descripcion,
			&item.Cirujano,
			&item.Servicio,
			&item.Estado,
			&horaInicio,
			&horaFin,
			&numeroCuenta,
		); errorEscaneo == nil {
			item.HoraInicio = horaInicio
			item.HoraFin = horaFin
			item.Cama = ""
			item.TiempoQuirurgico = ""
			item.DiagnosticoPre = ""
			item.DiagnosticoPost = ""
			item.OperacionIndicada = strings.TrimSpace(item.CodigoCPT + " - " + item.Descripcion)
			item.OperacionRealizada = strings.TrimSpace(item.CodigoCPT + " - " + item.Descripcion)
			item.PrimerAyudante = ""
			item.SegundoAyudante = ""
			item.Anestesiologo = ""
			item.Instrumentista = ""
			item.PersonalCirculante = ""
			item.TipoAnestesia = ""
			item.TipoOperacion = ""
			item.CuentaGasas = ""
			item.PiezaOperatoria = ""
			item.TipoIntervencion = ""
			item.RiesgoQuirurgico = ""
			item.Hallazgos = ""
			elementos = append(elementos, item)
		}
	}

	return elementos, nil
}

func (r *atencionCERepository) ListarProcedimientosHistorial(ctx context.Context, identificadorPaciente int) ([]domain.HistorialProcedimientoItem, error) {
	consulta := `
		SELECT 
			a.IdAtencion AS IdProcedimiento,
			a.IdAtencion,
			ISNULL(s.Codigo, '3040603') AS CodigoCPT,
			CONVERT(VARCHAR(10), a.FechaIngreso, 120) AS Fecha,
			ISNULL(s.Nombre, 'PROCEDIMIENTO MÉDICO') AS Descripcion,
			ISNULL(LTRIM(RTRIM(CONCAT(e.ApellidoPaterno, ' ', e.ApellidoMaterno, ' ', e.Nombres))), 'MÉDICO ESPECIALISTA') AS Medico,
			ISNULL(s.Nombre, 'CONSULTA ESPECIALIZADA') AS Servicio,
			'VALIDADO' AS Estado,
			ISNULL(CONVERT(VARCHAR(10), a.IdCuentaAtencion), '') AS NumeroCuenta,
			CONVERT(VARCHAR(16), a.FechaIngreso, 120) AS FechaHoraAtencion,
			ISNULL(m.ColegioMedico, '') AS CMP
		FROM dbo.Atenciones a WITH (NOLOCK)
		INNER JOIN dbo.Servicios s WITH (NOLOCK) ON a.IdServicioIngreso = s.IdServicio
		LEFT JOIN dbo.Medicos m WITH (NOLOCK) ON a.IdMedicoIngreso = m.IdMedico
		LEFT JOIN dbo.Empleados e WITH (NOLOCK) ON m.IdEmpleado = e.IdEmpleado
		WHERE a.IdPaciente = @p1 
		  AND (s.Nombre LIKE '%ENDOSCOP%' OR s.Nombre LIKE '%PROCEDIMIENTO%' OR s.Nombre LIKE '%BIOPSIA%' OR s.Nombre LIKE '%TERAPIA%' OR s.Nombre LIKE '%ECOGRAF%' OR s.Nombre LIKE '%RAYOS%' OR s.Nombre LIKE '%TOMOGRAF%' OR s.Nombre LIKE '%LABORATOR%' OR s.Nombre LIKE '%REHABILIT%' OR s.Nombre LIKE '%ODONTOLOG%' OR s.Nombre LIKE '%OFTALMOLOG%' OR s.Nombre LIKE '%DERMATOLOG%' OR s.Nombre LIKE '%CARDIOLOG%')

		UNION ALL

		SELECT 
			ISNULL(rd.idItem, a.IdAtencion) AS IdProcedimiento,
			a.IdAtencion,
			ISNULL(fcs.Codigo, '3040603') AS CodigoCPT,
			CONVERT(VARCHAR(10), ISNULL(rc.FechaReceta, a.FechaIngreso), 120) AS Fecha,
			ISNULL(fcs.Nombre, 'PROCEDIMIENTO MÉDICO CPT') AS Descripcion,
			ISNULL(LTRIM(RTRIM(CONCAT(e.ApellidoPaterno, ' ', e.ApellidoMaterno, ' ', e.Nombres))), 'MÉDICO TRATANTE') AS Medico,
			ISNULL(s.Nombre, 'CONSULTORIOS EXTERNOS') AS Servicio,
			'VALIDADO' AS Estado,
			ISNULL(CONVERT(VARCHAR(10), a.IdCuentaAtencion), '') AS NumeroCuenta,
			CONVERT(VARCHAR(16), ISNULL(rc.FechaReceta, a.FechaIngreso), 120) AS FechaHoraAtencion,
			ISNULL(m.ColegioMedico, '') AS CMP
		FROM dbo.Atenciones a WITH (NOLOCK)
		INNER JOIN dbo.RecetaCabecera rc WITH (NOLOCK) ON rc.idCuentaAtencion = a.IdCuentaAtencion
		INNER JOIN dbo.RecetaDetalle rd WITH (NOLOCK) ON rd.idReceta = rc.idReceta
		INNER JOIN dbo.FactCatalogoServicios fcs WITH (NOLOCK) ON fcs.IdServicio = rd.idItem
		LEFT JOIN dbo.Servicios s WITH (NOLOCK) ON a.IdServicioIngreso = s.IdServicio
		LEFT JOIN dbo.Medicos m WITH (NOLOCK) ON a.IdMedicoIngreso = m.IdMedico
		LEFT JOIN dbo.Empleados e WITH (NOLOCK) ON m.IdEmpleado = e.IdEmpleado
		WHERE a.IdPaciente = @p1 AND (fcs.EsCPT = 1 OR fcs.EsQuirurgico = 0)
		ORDER BY Fecha DESC`

	filas, errorConsulta := r.db.QueryContext(ctx, consulta, sql.Named("p1", identificadorPaciente))
	if errorConsulta != nil {
		return []domain.HistorialProcedimientoItem{}, nil
	}
	defer filas.Close()

	elementos := make([]domain.HistorialProcedimientoItem, 0)
	for filas.Next() {
		var item domain.HistorialProcedimientoItem
		var numeroCuenta, fechaHoraAtencion, colegioMedico string
		if errorEscaneo := filas.Scan(
			&item.IdProcedimiento,
			&item.IdAtencion,
			&item.CodigoCPT,
			&item.Fecha,
			&item.Descripcion,
			&item.Medico,
			&item.Servicio,
			&item.Estado,
			&numeroCuenta,
			&fechaHoraAtencion,
			&colegioMedico,
		); errorEscaneo == nil {
			item.NumeroCuenta = numeroCuenta
			item.FechaHoraAtencion = fechaHoraAtencion
			item.CMP = colegioMedico
			item.Especialidad = item.Servicio
			item.TipoSeguro = "SIS"
			item.Procedencia = "CONSULTORIOS EXTERNOS"
			item.Hallazgos = ""
			item.Diagnosticos = ""
			item.Observacion = ""
			elementos = append(elementos, item)
		}
	}

	return elementos, nil
}

func (r *atencionCERepository) ListarServiciosPorMedico(ctx context.Context, idTipoServicio int, idEmpleado int) ([]domain.ServicioMedicoItem, error) {
	if idTipoServicio <= 0 {
		idTipoServicio = 1
	}

	query := `EXEC dbo.usp_go_ListarServiciosSegunTipoServicioIdempleado @p1, @p2`
	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idTipoServicio), sql.Named("p2", idEmpleado))
	if err != nil {
		return []domain.ServicioMedicoItem{}, err
	}
	defer rows.Close()

	items := make([]domain.ServicioMedicoItem, 0)
	for rows.Next() {
		var item domain.ServicioMedicoItem
		var descripcionHospitalaria sql.NullString
		var nombreServicio sql.NullString
		if err := rows.Scan(
			&item.IdServicio,
			&descripcionHospitalaria,
			&nombreServicio,
		); err == nil {
			if descripcionHospitalaria.Valid {
				item.DescripcionHospitalaria = strings.TrimSpace(descripcionHospitalaria.String)
			}
			if nombreServicio.Valid {
				item.NombreServicio = strings.TrimSpace(nombreServicio.String)
			}
			items = append(items, item)
		}
	}

	return items, nil
}

func (r *atencionCERepository) ListarHistorialTriajePaciente(ctx context.Context, idPaciente int) ([]domain.PacienteHistorialTriajeItem, error) {
	query := `
		SELECT 
			t.IdTriaje,
			ISNULL(t.IdAtencion, 0) AS IdAtencion,
			CONVERT(VARCHAR(10), t.Fecha, 103) AS FechaRegistro,
			ISNULL(CONVERT(VARCHAR(5), t.Hora, 108), ISNULL(t.Hora, '')) AS HoraRegistro,
			ISNULL(t.PresionArterial, '') AS PresionArterial,
			ISNULL(t.Pulso, '') AS Pulso,
			ISNULL(TRY_CAST(t.FrecCardiaca AS INT), 0) AS FrecuenciaCardiaca,
			ISNULL(TRY_CAST(t.FrecRespiratoria AS INT), 0) AS FrecuenciaRespiratoria,
			ISNULL(TRY_CAST(t.Temperatura AS DECIMAL(5,2)), 0.0) AS Temperatura,
			ISNULL(TRY_CAST(t.SAT02 AS INT), 0) AS Saturacion,
			ISNULL(TRY_CAST(t.Peso AS DECIMAL(5,2)), 0.0) AS Peso,
			ISNULL(TRY_CAST(t.Talla AS DECIMAL(5,2)), 0.0) AS Talla,
			ISNULL(TRY_CAST(t.IMC AS DECIMAL(5,2)), 0.0) AS IMC,
			ISNULL(t.Hemoglobina, '') AS Hemoglobina,
			ISNULL(t.Observacion, '') AS Observacion
		FROM dbo.AtencionesTriaje t WITH (NOLOCK)
		WHERE t.IdPaciente = @p1
		ORDER BY t.Fecha DESC, t.IdTriaje DESC`

	rows, err := r.db.QueryContext(ctx, query, sql.Named("p1", idPaciente))
	if err != nil {
		return []domain.PacienteHistorialTriajeItem{}, fmt.Errorf("error al consultar historial de triaje del paciente: %w", err)
	}
	defer rows.Close()

	items := make([]domain.PacienteHistorialTriajeItem, 0)
	for rows.Next() {
		var item domain.PacienteHistorialTriajeItem
		if err := rows.Scan(
			&item.IdTriaje,
			&item.IdAtencion,
			&item.FechaRegistro,
			&item.HoraRegistro,
			&item.PresionArterial,
			&item.Pulso,
			&item.FrecuenciaCardiaca,
			&item.FrecuenciaRespiratoria,
			&item.Temperatura,
			&item.Saturacion,
			&item.Peso,
			&item.Talla,
			&item.IMC,
			&item.Hemoglobina,
			&item.Observacion,
		); err == nil {
			items = append(items, item)
		}
	}

	return items, nil
}
