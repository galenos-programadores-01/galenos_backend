package sqlserver

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/galenos-pro/appointments-api/internal/domain"
	"github.com/galenos-pro/appointments-api/internal/ports/output"
)

type sqlServerEvolucionRepository struct {
	db *sql.DB
}

func NewSqlServerEvolucionRepository(db *sql.DB) output.EvolucionRepository {
	return &sqlServerEvolucionRepository{db: db}
}

func (r *sqlServerEvolucionRepository) ListPatients(ctx context.Context, fini, ffin string, idUsuario int) ([]domain.PatientListItem, error) {
	fecha, err := time.Parse("2006-01-02", fini)
	if err != nil {
		return nil, fmt.Errorf("fini inválido: %w", err)
	}
	query := `EXEC [dbo].[usp_go_ListarPacientesSegunTipoServicio] @IdTipoServicio = @p1, @Fecha = @p2, @Filtro = @p3, @IdUsuario = @p4`
	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("p1", 2),
		sql.Named("p2", fecha),
		sql.Named("p3", ""),
		sql.Named("p4", idUsuario),
	)
	if err != nil {
		return nil, fmt.Errorf("error querying patients for tray: %w", err)
	}
	defer rows.Close()

	maps, err := rowsToMaps(rows)
	if err != nil {
		return nil, fmt.Errorf("error reading patient tray maps: %w", err)
	}

	var patients []domain.PatientListItem
	for _, m := range maps {
		patients = append(patients, mapToPatientListItem(m))
	}
	return patients, nil
}

func mapToPatientListItem(m map[string]any) domain.PatientListItem {
	var p domain.PatientListItem

	p.IdRegAtencion = getIntFallback(m, "idRegAtencion", "IdAtencion")
	p.IdPaciente = getIntFallback(m, "IdPaciente", "IdEpisodio")
	p.Historia = getStringFallback(m, "historia", "NroHistoriaClinica", "N/A")
	p.Nombre = getNombrePaciente(m)
	p.Edad = getStringFallback(m, "edad", "", "N/A")
	p.Sexo = getSexoFallback(m)
	p.Ubicacion = getStringFallback(m, "ubicacion", "Servicio", "Emergencia")
	p.Servicio = getStringFallback(m, "Servicio", "servicio", p.Ubicacion)
	p.Especialidad = getStringFallback(m, "Especialidad", "especialidad", "")
	p.Cama = getStringFallback(m, "cama", "", "NS")
	p.Estado = getStringFallback(m, "estado", "", "Pendiente")

	return p
}

func getIntFallback(m map[string]any, key1, key2 string) int {
	if val, ok := m[key1]; ok && val != nil {
		return int(val.(int64))
	}
	if key2 != "" {
		if val, ok := m[key2]; ok && val != nil {
			return int(val.(int64))
		}
	}
	return 0
}

func getStringFallback(m map[string]any, key1, key2, fallback string) string {
	if ptr := rowString(m, key1); ptr != nil && *ptr != "" {
		return *ptr
	}
	if key2 != "" {
		if ptr := rowString(m, key2); ptr != nil && *ptr != "" {
			return *ptr
		}
	}
	return fallback
}

func getNombrePaciente(m map[string]any) string {
	if ptr := rowString(m, "nombre", "Paciente"); ptr != nil && *ptr != "" {
		return *ptr
	}
	pat := getStringFallback(m, "ApellidoPaterno", "", "")
	mat := getStringFallback(m, "ApellidoMaterno", "", "")
	nom := getStringFallback(m, "PrimerNombre", "", "")
	return fmt.Sprintf("%s %s, %s", pat, mat, nom)
}

func getSexoFallback(m map[string]any) string {
	if ptr := rowString(m, "sexo"); ptr != nil && *ptr != "" {
		return *ptr
	}
	if ptr := rowInt64(m, "IdTipoSexo"); ptr != nil {
		return fmt.Sprintf("%d", *ptr)
	}
	return "0"
}

func (r *sqlServerEvolucionRepository) ListEvoluciones(ctx context.Context, idRegAtencion int) ([]domain.EvolucionFirma, error) {
	// 1. Obtener IdPaciente de la atención con fallbacks robustos
	var idPaciente int
	_ = r.db.QueryRowContext(ctx, "SELECT ISNULL(IdPaciente, 0) FROM dbo.Atenciones WITH (NOLOCK) WHERE IdAtencion = @p1", sql.Named("p1", idRegAtencion)).Scan(&idPaciente)
	if idPaciente == 0 {
		_ = r.db.QueryRowContext(ctx, "SELECT TOP 1 ISNULL(IdEpisodio, 0) FROM dbo.Tab_EvolucionesMedicas_H_E WITH (NOLOCK) WHERE NroAtencion = @p1 AND IdEpisodio > 0", sql.Named("p1", idRegAtencion)).Scan(&idPaciente)
	}
	if idPaciente == 0 {
		var count int
		_ = r.db.QueryRowContext(ctx, "SELECT COUNT(1) FROM dbo.Pacientes WITH (NOLOCK) WHERE IdPaciente = @p1", sql.Named("p1", idRegAtencion)).Scan(&count)
		if count > 0 {
			idPaciente = idRegAtencion
		}
	}
	log.Printf("[ListEvoluciones] idRegAtencion=%d, resolved idPaciente=%d", idRegAtencion, idPaciente)

	// 2. Obtener antecedentes del paciente con el SP moderno
	var antecedentesMap map[string]any
	if idPaciente > 0 {
		antRows, errAnt := r.db.QueryContext(ctx, "EXEC [dbo].[usp_go_PacientesDatosAdicionalesIdPaciente] @IdPaciente = @p1", sql.Named("p1", idPaciente))
		if errAnt == nil {
			if antRows.Next() {
				var idPac int
				var ant, antAlerg, antObst, antQuir, antFam, antPat, otrosComorb sql.NullString
				var fNacCalc sql.NullBool
				var hta, ob, dislip, an, hig, tiro, tb, fuma, ca sql.NullInt64

				if err := antRows.Scan(
					&idPac, &ant, &antAlerg, &antObst, &antQuir, &antFam, &antPat,
					&fNacCalc, &hta, &ob, &dislip, &an, &hig, &tiro, &tb, &fuma, &ca, &otrosComorb,
				); err == nil {
					antecedentesMap = map[string]any{
						"idPaciente":           idPac,
						"antecedentes":         ant.String,
						"antecedAlergico":      antAlerg.String,
						"antecedObstetrico":    antObst.String,
						"antecedQuirurgico":    antQuir.String,
						"antecedFamiliar":      antFam.String,
						"antecedPatologico":    antPat.String,
						"hipertensionArterial": hta.Int64,
						"obesidad":             ob.Int64,
						"dislipidemia":         dislip.Int64,
						"anemia":               an.Int64,
						"higadoGraso":          hig.Int64,
						"enfTiroidea":          tiro.Int64,
						"tuberculosis":         tb.Int64,
						"fumaActualmente":      fuma.Int64,
						"cancer":               ca.Int64,
						"otrosComorbilidad":    otrosComorb.String,
					}
				}
			}
			antRows.Close()
		} else {
			log.Printf("[ListEvoluciones] ERROR ejecutando SP antecedentes: %v", errAnt)
		}
	}
	log.Printf("[ListEvoluciones] antecedentesMap=%v", antecedentesMap != nil)

	// 3. Obtener diagnósticos asociados a la atención médica con el SP moderno
	var dxList []map[string]any
	seenDx := make(map[string]bool)
	dxRows, errDx := r.db.QueryContext(ctx, "EXEC [dbo].[usp_go_ObtenerDiagnosticosAtencion] @IdAtencion = @p1, @IdPrimeraAtencion = 0, @IdEvolucion = 0", sql.Named("p1", idRegAtencion))
	if errDx == nil {
		for dxRows.Next() {
			var idAd, idAt, idDx, idSubDx sql.NullInt64
			var cie, desc, descL, tCod, tDx, his, his1, his2, his3 sql.NullString
			if err := dxRows.Scan(&idAd, &idAt, &idDx, &cie, &desc, &descL, &tCod, &tDx, &idSubDx, &his, &his1, &his2, &his3); err == nil {
				tipoDiag := tDx.String
				if tipoDiag == "" {
					tipoDiag = "Presuntivo"
				}
				key := strings.ToUpper(strings.TrimSpace(cie.String)) + "|" + strings.ToUpper(strings.TrimSpace(desc.String))
				if !seenDx[key] {
					seenDx[key] = true
					dxList = append(dxList, map[string]any{
						"cie10":       cie.String,
						"descripcion": desc.String,
						"tipo":        tipoDiag,
						"condicion":   "Principal",
					})
				}
			}
		}
		dxRows.Close()
	} else {
		log.Printf("[ListEvoluciones] ERROR ejecutando SP diagnosticos: %v", errDx)
	}
	log.Printf("[ListEvoluciones] diagnosticos encontrados: %d", len(dxList))

	// 4. Obtener síntomas referidos de la atención con el SP moderno
	var sintomasList []string
	seenSint := make(map[string]bool)
	sintRows, errSint := r.db.QueryContext(ctx, "EXEC [dbo].[usp_go_ObtenerAtencionSintomas] @IdRegAtencion = @p1", sql.Named("p1", idRegAtencion))
	if errSint == nil {
		for sintRows.Next() {
			var idSint, idAtScan sql.NullInt64
			var sis, sin sql.NullString
			var esPred interface{}
			if err := sintRows.Scan(&idSint, &idAtScan, &sis, &sin, &esPred); err == nil {
				sintTexto := strings.TrimSpace(sin.String)
				if sin.Valid && sintTexto != "" {
					key := strings.ToUpper(sintTexto)
					if !seenSint[key] {
						seenSint[key] = true
						sintomasList = append(sintomasList, sintTexto)
					}
				}
			}
		}
		sintRows.Close()
	} else {
		log.Printf("[ListEvoluciones] ERROR ejecutando SP sintomas: %v", errSint)
	}
	log.Printf("[ListEvoluciones] sintomas encontrados: %d -> %v", len(sintomasList), sintomasList)

	// 5. Obtener IdCuentaAtencion para consultar recetas por evolución
	var idCuenta int
	_ = r.db.QueryRowContext(ctx, "SELECT ISNULL(IdCuentaAtencion, 0) FROM dbo.Atenciones WITH (NOLOCK) WHERE IdAtencion = @p1", sql.Named("p1", idRegAtencion)).Scan(&idCuenta)

	// 6. Obtener interconsultas solicitadas para la atención con el SP moderno
	var interconsultasList []map[string]any
	seenIc := make(map[string]bool)
	icRows, errIc := r.db.QueryContext(ctx, "EXEC [dbo].[usp_go_ListarInterconsultasPorAtencion] @IdAtencion = @p1", sql.Named("p1", idRegAtencion))
	if errIc == nil {
		for icRows.Next() {
			var idIc, idAt, idEsp, idMedDest int
			var mot string
			var fSol sql.NullTime
			var est sql.NullString
			if err := icRows.Scan(&idIc, &idAt, &idEsp, &idMedDest, &mot, &fSol, &est); err == nil {
				var espNombre sql.NullString
				_ = r.db.QueryRowContext(ctx, "SELECT ISNULL(Nombre, '') FROM dbo.Especialidades WITH (NOLOCK) WHERE IdEspecialidad = @p1", sql.Named("p1", idEsp)).Scan(&espNombre)
				nomEsp := espNombre.String
				if nomEsp == "" {
					nomEsp = fmt.Sprintf("Especialidad #%d", idEsp)
				}
				key := strings.ToUpper(strings.TrimSpace(nomEsp)) + "|" + strings.ToUpper(strings.TrimSpace(mot))
				if !seenIc[key] {
					seenIc[key] = true
					interconsultasList = append(interconsultasList, map[string]any{
						"especialidad": nomEsp,
						"motivo":       mot,
						"estado":       est.String,
					})
				}
			}
		}
		icRows.Close()
	}

	// 7. Consultar las notas de evolución clínica de la tabla moderna Tab_EvolucionesMedicas_H_E
	queryModerna := `
		SELECT 
			e.IdEvolucion,
			e.NroAtencion,
			ISNULL(e.UbicacionArchivo, ''),
			ISNULL(e.UsuarioCreacion, 0),
			CONVERT(varchar(19), e.FechaAtencion, 120),
			ISNULL(e.EstadoFirma, 1),
			ISNULL(RTRIM(LTRIM(CONCAT(emp.ApellidoPaterno, ' ', emp.ApellidoMaterno, ' ', emp.Nombres))), ''),
			ISNULL(e.Motivo, ''),
			ISNULL(e.Subjetivo, ''),
			ISNULL(e.EscalaDolor, 0),
			ISNULL(e.PASistolica, 0),
			ISNULL(e.PADiastolica, 0),
			ISNULL(e.FrecuenciaCardiaca, 0),
			ISNULL(e.FrecuenciaRespiratoria, 0),
			ISNULL(e.Temperatura, 0.0),
			ISNULL(e.SaturacionOxigeno, 0),
			ISNULL(e.Peso, 0.0),
			ISNULL(e.Talla, 0.0),
			ISNULL(e.IMC, 0.0),
			ISNULL(e.Glicemia, 0.0),
			ISNULL(e.ExamenFisicoGeneral, ''),
			ISNULL(e.ExamenFisicoPiel, ''),
			ISNULL(e.ExamenFisicoCabezaCuello, ''),
			ISNULL(e.ExamenFisicoToraxPulmon, ''),
			ISNULL(e.ExamenFisicoCorazon, ''),
			ISNULL(e.ExamenFisicoAbdomen, ''),
			ISNULL(e.ExamenFisicoGenitourinario, ''),
			ISNULL(e.ExamenFisicoExtremidadesOsteomuscular, ''),
			ISNULL(e.ExamenFisicoNeurologicoMental, ''),
			ISNULL(e.IdEstadoClinico, 0),
			ISNULL(e.IdPronostico, 0),
			ISNULL(e.IndicacionDieta, ''),
			ISNULL(e.IndicacionReposo, ''),
			ISNULL(e.IndicacionHidratacion, ''),
			ISNULL(e.IndicacionOxigeno, ''),
			ISNULL(e.IndicacionRestriccion, ''),
			ISNULL(e.Sugerencia, '')
		FROM dbo.Tab_EvolucionesMedicas_H_E e WITH (NOLOCK)
		LEFT JOIN dbo.Empleados emp WITH (NOLOCK) ON emp.IdEmpleado = e.UsuarioCreacion
		WHERE e.NroAtencion = @p1
		ORDER BY e.FechaAtencion DESC, e.IdEvolucion DESC
	`
	rowsMod, err := r.db.QueryContext(ctx, queryModerna, sql.Named("p1", idRegAtencion))
	if err != nil {
		return nil, fmt.Errorf("error querying evoluciones: %w", err)
	}
	defer rowsMod.Close()

	var evolutions []domain.EvolucionFirma
	for rowsMod.Next() {
		var idEvol, nroAtenc, usrCreacion, estFirma, escalaDolor, paSis, paDia, fc, fr, sat int
		var idEstClin, idPron int
		var temp, peso, talla, imc, glicemia float64
		var ubicacion, fechaAtenc, medicoNom, motivo, subjetivo, sugerencia string
		var efGen, efPiel, efCabeza, efTorax, efCorazon, efAbdomen, efGenito, efOsteo, efNeuro string
		var dieta, reposo, hidratacion, oxigeno, restriccion string

		if err := rowsMod.Scan(
			&idEvol, &nroAtenc, &ubicacion, &usrCreacion, &fechaAtenc, &estFirma, &medicoNom,
			&motivo, &subjetivo, &escalaDolor, &paSis, &paDia, &fc, &fr, &temp, &sat,
			&peso, &talla, &imc, &glicemia,
			&efGen, &efPiel, &efCabeza, &efTorax, &efCorazon, &efAbdomen, &efGenito, &efOsteo, &efNeuro,
			&idEstClin, &idPron,
			&dieta, &reposo, &hidratacion, &oxigeno, &restriccion, &sugerencia,
		); err == nil {
			var paStr string
			if paSis > 0 || paDia > 0 {
				paStr = fmt.Sprintf("%d/%d", paSis, paDia)
			}

			var estClinStr string
			if idEstClin == 1 {
				estClinStr = "Mejoría"
			} else if idEstClin == 2 {
				estClinStr = "Estacionario"
			} else if idEstClin == 3 {
				estClinStr = "Desfavorable"
			}

			var pronStr string
			if idPron == 1 {
				pronStr = "Bueno"
			} else if idPron == 2 {
				pronStr = "Reservado"
			} else if idPron == 3 {
				pronStr = "Malo"
			}

			fechaStr := fechaAtenc
			horaStr := ""
			if len(fechaAtenc) >= 19 {
				fechaStr = fechaAtenc[:10]
				horaStr = fechaAtenc[11:16]
			}

			sistemas := []string{
				"General", "Piel y Faneras", "Cabeza y Cuello", "Tórax y Pulmones",
				"Cardiovascular", "Abdomen", "Genitourinario", "Osteomuscular", "Neurológico",
			}
			hallazgos := []string{
				efGen, efPiel, efCabeza, efTorax, efCorazon, efAbdomen, efGenito, efOsteo, efNeuro,
			}
			var examenFisicoList []map[string]any
			for k := 0; k < len(sistemas); k++ {
				h := hallazgos[k]
				isNormal := (h == "" || strings.EqualFold(h, "Sin hallazgos") || strings.EqualFold(h, "Normal") || strings.EqualFold(h, "Conservado"))
				examenFisicoList = append(examenFisicoList, map[string]any{
					"sistema":  sistemas[k],
					"normal":   isNormal,
					"hallazgo": h,
				})
			}

			// Consultar recetas farmacológicas específicas de ESTA evolución
			var farmacosEvolucion []map[string]any
			if idEvol > 0 && idCuenta > 0 {
				farmRows, errFarm := r.db.QueryContext(ctx, `
					SELECT DISTINCT
						ISNULL(p.Nombre, ''),
						ISNULL(rd.DescripcionAdicional, ''),
						ISNULL(rd.Cantidad, 0),
						ISNULL(rd.Duracion, '')
					FROM dbo.RecetaDetalle rd WITH (NOLOCK)
					INNER JOIN dbo.RecetaCabecera rc WITH (NOLOCK) ON rc.IdReceta = rd.IdReceta
					INNER JOIN dbo.Producto p WITH (NOLOCK) ON p.IdProducto = rd.IdProducto
					WHERE rc.IdEvolucion = @p1
					  AND rc.IdCuentaAtencion = @p2
				`, sql.Named("p1", idEvol), sql.Named("p2", idCuenta))
				if errFarm == nil {
					for farmRows.Next() {
						var nombreMed, descAdicional, duracion string
						var cantidad int
						if farmRows.Scan(&nombreMed, &descAdicional, &cantidad, &duracion) == nil && nombreMed != "" {
							farmacosEvolucion = append(farmacosEvolucion, map[string]any{
								"medicamento": nombreMed,
								"cantidad":    cantidad,
								"dosis":       descAdicional,
								"duracion":    duracion,
							})
						}
					}
					farmRows.Close()
				}
			}

			jsonPayload := map[string]any{
				"cabecera": map[string]any{
					"numeroEvolucion": idEvol,
					"fecha":           fechaStr,
					"hora":            horaStr,
					"medicoTratante":  medicoNom,
					"tipoAtencion":    "Emergencia / Hospitalización",
					"estado":          "Firmado",
					"firmaDni":        ubicacion,
				},
				"antecedentes": antecedentesMap,
				"motivo": map[string]any{
					"tipo":        "Consulta Médica",
					"descripcion": motivo,
				},
				"subjetivo": map[string]any{
					"evolucionSintomas": subjetivo,
					"escalaDolor":       escalaDolor,
				},
				"sintomas": sintomasList,
				"signosVitales": map[string]any{
					"presionArterial":        paStr,
					"frecuenciaCardiaca":     fc,
					"frecuenciaRespiratoria": fr,
					"temperatura":            temp,
					"saturacionOxigeno":      sat,
					"peso":                   peso,
					"talla":                  talla,
					"imc":                    imc,
					"glucemia":               glicemia,
				},
				"examenFisico": examenFisicoList,
				"diagnosticos": dxList,
				"evaluacion": map[string]any{
					"estadoClinico":  estClinStr,
					"pronostico":     pronStr,
					"evolucionLibre": sugerencia,
				},
				"evolucionLibre": sugerencia,
				"plan": map[string]any{
					"farmacologico":  farmacosEvolucion,
					"interconsultas": interconsultasList,
					"indicacionesGenerales": map[string]any{
						"dieta":         dieta,
						"reposo":        reposo,
						"hidratacion":   hidratacion,
						"oxigeno":       oxigeno,
						"restricciones": restriccion,
					},
				},
			}

			jsonBytes, _ := json.Marshal(jsonPayload)
			b64Str := base64.StdEncoding.EncodeToString(jsonBytes)

			evolutions = append(evolutions, domain.EvolucionFirma{
				IdRegAtencion:      nroAtenc,
				IdFirma:            idEvol,
				NombreDocumento:    "EvolucionMedica",
				NombreArchivo:      fmt.Sprintf("EV-%d", idEvol),
				RutaBase:           ubicacion,
				DataB64:            b64Str,
				IdEmpleadoRegistra: usrCreacion,
				MedicoNombre:       medicoNom,
				FechaRegistro:      fechaAtenc,
				Estado:             estFirma,
			})
		}
	}

	return evolutions, nil
}

func (r *sqlServerEvolucionRepository) SaveEvolucion(ctx context.Context, evolution domain.EvolucionFirma) error {
	var payload map[string]any
	if evolution.DataB64 != "" {
		jsonBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(evolution.DataB64))
		if err == nil {
			_ = json.Unmarshal(jsonBytes, &payload)
		}
	}

	nroAtencion := evolution.IdRegAtencion
	usuarioCreacion := evolution.IdEmpleadoRegistra
	motivo := evolution.NombreDocumento
	estadoRegistro := 1
	estadoFirma := 1
	ubicacion := evolution.RutaBase

	item := domain.EvolucionMedicaInsert{
		NroAtencion:      &nroAtencion,
		UsuarioCreacion:  usuarioCreacion,
		Motivo:           &motivo,
		EstadoRegistro:   &estadoRegistro,
		EstadoFirma:      &estadoFirma,
		UbicacionArchivo: &ubicacion,
	}

	if payload != nil {
		if cab, ok := payload["cabecera"].(map[string]any); ok {
			if fdni, ok := cab["firmaDni"].(string); ok && fdni != "" {
				item.UbicacionArchivo = &fdni
			}
		}
		if mot, ok := payload["motivo"].(map[string]any); ok {
			if desc, ok := mot["descripcion"].(string); ok && desc != "" {
				item.Motivo = &desc
			}
		}
		if subj, ok := payload["subjetivo"].(map[string]any); ok {
			if evSin, ok := subj["evolucionSintomas"].(string); ok && evSin != "" {
				item.Subjetivo = &evSin
			}
			if ed, ok := subj["escalaDolor"].(float64); ok {
				edInt := int(ed)
				item.EscalaDolor = &edInt
			}
		}
		if sv, ok := payload["signosVitales"].(map[string]any); ok {
			if pa, ok := sv["presionArterial"].(string); ok && pa != "" {
				partes := strings.Split(pa, "/")
				if len(partes) == 2 {
					if pas, err := strconv.Atoi(strings.TrimSpace(partes[0])); err == nil {
						item.PASistolica = &pas
					}
					if pad, err := strconv.Atoi(strings.TrimSpace(partes[1])); err == nil {
						item.PADiastolica = &pad
					}
				}
			}
			if fc, ok := sv["frecuenciaCardiaca"].(float64); ok {
				fcInt := int(fc)
				item.FrecuenciaCardiaca = &fcInt
			}
			if fr, ok := sv["frecuenciaRespiratoria"].(float64); ok {
				frInt := int(fr)
				item.FrecuenciaRespiratoria = &frInt
			}
			if temp, ok := sv["temperatura"].(float64); ok {
				item.Temperatura = &temp
			}
			if sat, ok := sv["saturacionOxigeno"].(float64); ok {
				satInt := int(sat)
				item.SaturacionOxigeno = &satInt
			}
			if peso, ok := sv["peso"].(float64); ok {
				item.Peso = &peso
			}
			if talla, ok := sv["talla"].(float64); ok {
				item.Talla = &talla
			}
			if imc, ok := sv["imc"].(float64); ok {
				item.IMC = &imc
			} else if imcStr, ok := sv["imc"].(string); ok && imcStr != "" && imcStr != "—" {
				if imcNum, err := strconv.ParseFloat(imcStr, 64); err == nil {
					item.IMC = &imcNum
				}
			}
			if glu, ok := sv["glucemia"].(float64); ok {
				item.Glicemia = &glu
			}
		}
		if ef, ok := payload["examenFisico"].([]any); ok {
			getHallazgo := func(idx int) *string {
				if idx < len(ef) {
					if efMap, ok := ef[idx].(map[string]any); ok {
						if h, ok := efMap["hallazgo"].(string); ok && h != "" {
							return &h
						}
					}
				}
				return nil
			}
			item.ExamenFisicoGeneral = getHallazgo(0)
			item.ExamenFisicoPiel = getHallazgo(1)
			item.ExamenFisicoCabezaCuello = getHallazgo(2)
			item.ExamenFisicoToraxPulmon = getHallazgo(3)
			item.ExamenFisicoCorazon = getHallazgo(4)
			item.ExamenFisicoAbdomen = getHallazgo(5)
			item.ExamenFisicoGenitourinario = getHallazgo(6)
			item.ExamenFisicoExtremidadesOsteomuscular = getHallazgo(7)
			item.ExamenFisicoNeurologicoMental = getHallazgo(8)
		}
		if ev, ok := payload["evaluacion"].(map[string]any); ok {
			if est, ok := ev["estadoClinico"].(string); ok {
				var idEst int
				switch est {
				case "Mejoría":
					idEst = 1
				case "Estacionario":
					idEst = 2
				case "Desfavorable":
					idEst = 3
				}
				if idEst > 0 {
					item.IdEstadoClinico = &idEst
				}
			}
			if pro, ok := ev["pronostico"].(string); ok {
				var idPro int
				switch pro {
				case "Bueno":
					idPro = 1
				case "Reservado":
					idPro = 2
				case "Malo":
					idPro = 3
				}
				if idPro > 0 {
					item.IdPronostico = &idPro
				}
			}
		}
		if plan, ok := payload["plan"].(map[string]any); ok {
			if ig, ok := plan["indicacionesGenerales"].(map[string]any); ok {
				if d, ok := ig["dieta"].(string); ok && d != "" {
					item.IndicacionDieta = &d
				}
				if r, ok := ig["reposo"].(string); ok && r != "" {
					item.IndicacionReposo = &r
				}
				if h, ok := ig["hidratacion"].(string); ok && h != "" {
					item.IndicacionHidratacion = &h
				}
				if o, ok := ig["oxigeno"].(string); ok && o != "" {
					item.IndicacionOxigeno = &o
				}
				if res, ok := ig["restricciones"].(string); ok && res != "" {
					item.IndicacionRestriccion = &res
				}
			}
		}
		if evLib, ok := payload["evolucionLibre"].(string); ok && evLib != "" {
			item.Sugerencia = &evLib
		}
	}

	_, _, err := r.InsertEvolucionMedica(ctx, item)
	return err
}

func (r *sqlServerEvolucionRepository) ListBandeja(ctx context.Context, fechaInicio, fechaFin, filtro string) ([]domain.EvolucionBandejaItem, error) {
	query := `EXEC dbo.usp_go_EvolucionesMedicas_Bandeja @FechaInicio = @p1, @FechaFin = @p2, @Filtro = @p3`

	var fInicio, fFin interface{} = nil, nil
	if fechaInicio != "" {
		fInicio = fechaInicio
	}
	if fechaFin != "" {
		fFin = fechaFin
	}

	rows, err := r.db.QueryContext(ctx, query,
		sql.Named("p1", fInicio),
		sql.Named("p2", fFin),
		sql.Named("p3", filtro),
	)
	if err != nil {
		return nil, fmt.Errorf("error al consultar bandeja evoluciones: %w", err)
	}
	defer rows.Close()

	var list []domain.EvolucionBandejaItem
	for rows.Next() {
		var item domain.EvolucionBandejaItem
		var idEpisodio, idCuentaAtencion, escalaDolor, glasgow, paSistolica, paDiastolica, fc, fr, satO2, idEstadoClinico, idPronostico, estadoFirma, usrCreacion, estRegistro sql.NullInt64
		var temp sql.NullFloat64
		var fAtencion, paciente, documento, motivo, fFirma, fCreacion, eqCreacion sql.NullString

		if err := rows.Scan(
			&item.IdEvolucion, &idEpisodio, &item.NroAtencion, &fAtencion,
			&item.IdPaciente, &paciente, &documento, &idCuentaAtencion, &motivo,
			&escalaDolor, &glasgow, &paSistolica, &paDiastolica, &fc, &fr,
			&temp, &satO2, &idEstadoClinico, &idPronostico, &estadoFirma,
			&fFirma, &usrCreacion, &fCreacion, &eqCreacion, &estRegistro,
		); err != nil {
			return nil, fmt.Errorf("error escaneando item de bandeja evoluciones: %w", err)
		}

		if idEpisodio.Valid {
			v := int(idEpisodio.Int64)
			item.IdEpisodio = &v
		}
		if fAtencion.Valid {
			item.FechaAtencion = fAtencion.String
		}
		if paciente.Valid {
			item.Paciente = paciente.String
		}
		if documento.Valid {
			item.Documento = documento.String
		}
		if idCuentaAtencion.Valid {
			item.IdCuentaAtencion = int(idCuentaAtencion.Int64)
		}
		if motivo.Valid {
			item.Motivo = motivo.String
		}
		if escalaDolor.Valid {
			v := int(escalaDolor.Int64)
			item.EscalaDolor = &v
		}
		if glasgow.Valid {
			v := int(glasgow.Int64)
			item.Glasgow = &v
		}
		if paSistolica.Valid {
			v := int(paSistolica.Int64)
			item.PASistolica = &v
		}
		if paDiastolica.Valid {
			v := int(paDiastolica.Int64)
			item.PADiastolica = &v
		}
		if fc.Valid {
			v := int(fc.Int64)
			item.FrecuenciaCardiaca = &v
		}
		if fr.Valid {
			v := int(fr.Int64)
			item.FrecuenciaRespiratoria = &v
		}
		if temp.Valid {
			v := temp.Float64
			item.Temperatura = &v
		}
		if satO2.Valid {
			v := int(satO2.Int64)
			item.SaturacionOxigeno = &v
		}
		if idEstadoClinico.Valid {
			v := int(idEstadoClinico.Int64)
			item.IdEstadoClinico = &v
		}
		if idPronostico.Valid {
			v := int(idPronostico.Int64)
			item.IdPronostico = &v
		}
		if estadoFirma.Valid {
			v := int(estadoFirma.Int64)
			item.EstadoFirma = &v
		}
		if fFirma.Valid {
			v := fFirma.String
			item.FechaFirma = &v
		}
		if usrCreacion.Valid {
			v := int(usrCreacion.Int64)
			item.UsuarioCreacion = &v
		}
		if fCreacion.Valid {
			v := fCreacion.String
			item.FechaCreacion = &v
		}
		if eqCreacion.Valid {
			v := eqCreacion.String
			item.EquipoCreacion = &v
		}
		if estRegistro.Valid {
			v := int(estRegistro.Int64)
			item.EstadoRegistro = &v
		}

		list = append(list, item)
	}

	return list, rows.Err()
}

func (r *sqlServerEvolucionRepository) InsertEvolucionMedica(ctx context.Context, item domain.EvolucionMedicaInsert) (int, string, error) {
	query := `
		EXEC dbo.usp_go_EvolucionesMedicas_Insertar
			@IdEpisodio = @p1,
			@NroAtencion = @p2,
			@IdServicioActual = @p3,
			@FechaAtencion = @p4,
			@IdTipoAtencion = @p5,
			@Motivo = @p6,
			@Subjetivo = @p7,
			@EscalaDolor = @p8,
			@Glasgow = @p9,
			@PASistolica = @p10,
			@PADiastolica = @p11,
			@FrecuenciaCardiaca = @p12,
			@FrecuenciaRespiratoria = @p13,
			@Temperatura = @p14,
			@SaturacionOxigeno = @p15,
			@Peso = @p16,
			@Talla = @p17,
			@IMC = @p18,
			@Glicemia = @p19,
			@ExamenFisicoGeneral = @p20,
			@ExamenFisicoPiel = @p21,
			@ExamenFisicoCabezaCuello = @p22,
			@ExamenFisicoToraxPulmon = @p23,
			@ExamenFisicoCorazon = @p24,
			@ExamenFisicoAbdomen = @p25,
			@ExamenFisicoGenitourinario = @p26,
			@ExamenFisicoExtremidadesOsteomuscular = @p27,
			@ExamenFisicoNeurologicoMental = @p28,
			@IdEstadoClinico = @p29,
			@IdPronostico = @p30,
			@IndicacionDieta = @p31,
			@IndicacionReposo = @p32,
			@IndicacionHidratacion = @p33,
			@IndicacionOxigeno = @p34,
			@IndicacionRestriccion = @p35,
			@Sugerencia = @p36,
			@UsuarioCreacion = @p37,
			@EquipoCreacion = @p38,
			@EstadoRegistro = @p39,
			@EstadoFirma = @p40,
			@UbicacionArchivo = @p41,
			@IdEvolucion = @p42 OUTPUT,
			@Mensaje = @p43 OUTPUT
	`

	var idEpisodio *int
	if item.IdEpisodio != nil {
		idEpisodio = item.IdEpisodio
	} else if item.IdPaciente > 0 {
		idEpisodio = &item.IdPaciente
	}

	var nroAtencion *int
	if item.NroAtencion != nil {
		nroAtencion = item.NroAtencion
	} else if item.IdAtencion > 0 {
		nroAtencion = &item.IdAtencion
	}

	var motivo *string
	if item.Motivo != nil {
		motivo = item.Motivo
	} else if item.MotivoConsulta != nil {
		motivo = item.MotivoConsulta
	}

	var subjetivo *string
	if item.Subjetivo != nil {
		subjetivo = item.Subjetivo
	} else if item.Anamnesis != nil {
		subjetivo = item.Anamnesis
	} else if item.TiempoEnfermedad != nil {
		subjetivo = item.TiempoEnfermedad
	}

	var idTipoAtencion *int
	if item.IdTipoAtencion != nil {
		idTipoAtencion = item.IdTipoAtencion
	} else if item.IdTipoGravedad != nil {
		idTipoAtencion = item.IdTipoGravedad
	}

	var idEvolucion sql.NullInt64
	if item.IdEvolucion != nil && *item.IdEvolucion > 0 {
		idEvolucion = sql.NullInt64{Int64: int64(*item.IdEvolucion), Valid: true}
	}
	var mensaje string

	_, err := r.db.ExecContext(ctx, query,
		sql.Named("p1", idEpisodio),
		sql.Named("p2", nroAtencion),
		sql.Named("p3", item.IdServicioActual),
		sql.Named("p4", item.FechaAtencion),
		sql.Named("p5", idTipoAtencion),
		sql.Named("p6", motivo),
		sql.Named("p7", subjetivo),
		sql.Named("p8", item.EscalaDolor),
		sql.Named("p9", item.Glasgow),
		sql.Named("p10", item.PASistolica),
		sql.Named("p11", item.PADiastolica),
		sql.Named("p12", item.FrecuenciaCardiaca),
		sql.Named("p13", item.FrecuenciaRespiratoria),
		sql.Named("p14", item.Temperatura),
		sql.Named("p15", item.SaturacionOxigeno),
		sql.Named("p16", item.Peso),
		sql.Named("p17", item.Talla),
		sql.Named("p18", item.IMC),
		sql.Named("p19", item.Glicemia),
		sql.Named("p20", item.ExamenFisicoGeneral),
		sql.Named("p21", item.ExamenFisicoPiel),
		sql.Named("p22", item.ExamenFisicoCabezaCuello),
		sql.Named("p23", item.ExamenFisicoToraxPulmon),
		sql.Named("p24", item.ExamenFisicoCorazon),
		sql.Named("p25", item.ExamenFisicoAbdomen),
		sql.Named("p26", item.ExamenFisicoGenitourinario),
		sql.Named("p27", item.ExamenFisicoExtremidadesOsteomuscular),
		sql.Named("p28", item.ExamenFisicoNeurologicoMental),
		sql.Named("p29", item.IdEstadoClinico),
		sql.Named("p30", item.IdPronostico),
		sql.Named("p31", item.IndicacionDieta),
		sql.Named("p32", item.IndicacionReposo),
		sql.Named("p33", item.IndicacionHidratacion),
		sql.Named("p34", item.IndicacionOxigeno),
		sql.Named("p35", item.IndicacionRestriccion),
		sql.Named("p36", item.Sugerencia),
		sql.Named("p37", item.UsuarioCreacion),
		sql.Named("p38", item.EquipoCreacion),
		sql.Named("p39", item.EstadoRegistro),
		sql.Named("p40", item.EstadoFirma),
		sql.Named("p41", item.UbicacionArchivo),
		sql.Named("p42", sql.Out{Dest: &idEvolucion}),
		sql.Named("p43", sql.Out{Dest: &mensaje}),
	)

	if err != nil {
		return 0, "", fmt.Errorf("error ejecutando insercion de evolucion medica: %w", err)
	}

	id := 0
	if idEvolucion.Valid {
		id = int(idEvolucion.Int64)
	}

	return id, mensaje, nil
}
