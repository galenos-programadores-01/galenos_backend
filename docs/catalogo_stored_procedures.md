# Catálogo de Procedimientos Almacenados (Stored Procedures) — Galenos Pro

Este documento contiene el inventario exhaustivo de los procedimientos almacenados (SPs) de SQL Server utilizados en el sistema Galenos Pro, clasificando los **SPs modernos (`usp_go_`)**, los **antiguos reemplazados**, los **antiguos aún en uso** y el **catálogo completo de los 60 SPs `usp_go_`** verificados directamente en la base de datos.

---

## 1. Procedimientos Almacenados Modernos (`usp_go_`) Usados Actualmente

Estos procedimientos forman parte del nuevo estándar de la API en Go y se encuentran activos en los repositorios del backend:

| Procedimiento Almacenado | Módulo / Pantalla | Descripción y Propósito | Archivo Fuente |
| :--- | :--- | :--- | :--- |
| `[dbo].[usp_go_ListarPacientesSegunTipoServicio]` | Hospitalización / Bandeja | Lista pacientes hospitalizados o de emergencia por fecha y servicio. | [`evolucion_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/evolucion_repository.go) |
| `[dbo].[usp_go_EvolucionesMedicas_Bandeja]` | Evolución / Bandeja | Consulta el historial de evoluciones clínicas estructuradas con filtros de fecha. | [`evolucion_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/evolucion_repository.go) |
| `[dbo].[usp_go_EvolucionesMedicas_Insertar]` | Evolución / Firma y Cierre | Inserta el registro clínico estructurado de la evolución médica (signos vitales, examen físico, pronóstico, indicaciones). | [`evolucion_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/evolucion_repository.go) |
| `[dbo].[usp_go_InsertarEvolucionSintomas]` | Evolución / Síntomas | Persiste los síntomas marcados por el médico para la atención en formato JSON. | [`sintoma_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/sintoma_repository.go) |
| `[dbo].[usp_go_ObtenerAtencionSintomas]` | Evolución / Síntomas | Recupera los síntomas registrados previamente para una atención médica. | [`sintoma_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/sintoma_repository.go) |
| `[dbo].[usp_go_SelectRecetaFrecuenciaSelecionarTodos]` | Plan / Receta Médica | Lista las frecuencias horarias de administración de medicamentos. | [`catalog_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/catalog_repository.go) |
| `[dbo].[Usp_go_SelectRecetaUndDosisSelecionarTodos]` | Plan / Receta Médica | Lista las unidades de dosificación de fármacos. | [`catalog_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/catalog_repository.go) |
| `[dbo].[usp_go_RecetasListadoViasAdministracion]` | Plan / Receta Médica | Lista las vías de administración de medicamentos. | [`catalog_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/catalog_repository.go) |
| `[dbo].[usp_go_SelectMedicamentosFiltro]` | Plan / Receta Médica | Buscador de medicamentos con stock, precio y antecedentes de prescripción. | [`catalog_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/catalog_repository.go) |
| `[dbo].[usp_go_BuscarCatalogoExamenes]` | Plan / Solicitud Exámenes | Buscador de exámenes de laboratorio, imágenes y procedimientos CPT. | [`catalog_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/catalog_repository.go) |
| `[dbo].[usp_go_SelectDiagnosticos]` | Evaluación / CIE-10 | Buscador predictivo de diagnósticos CIE-10 para la atención. | [`diagnostico_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/diagnostico_repository.go) |
| `[dbo].[usp_go_ListarDiagnosticos]` | Evaluación / CIE-10 | Catálogo simple de diagnósticos CIE-10. | [`diagnostico_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/diagnostico_repository.go) |
| `[dbo].[usp_go_ObtenerDiagnosticosAtencion]` | Evaluación / CIE-10 | Lista los diagnósticos previamente vinculados a la atención actual. | [`diagnostico_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/diagnostico_repository.go) |
| `[dbo].[usp_go_AtencionesDiagnosticosAgregar]` | Evaluación / CIE-10 | Registra un diagnóstico para la atención con correlativo y subclasificación HIS. | [`diagnostico_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/diagnostico_repository.go) |
| `[dbo].[usp_go_SelectHistoriaLaboratorio]` | Plan / Resultados Lab | Historial de pruebas de laboratorio clínico del paciente. | [`resultado_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/resultado_repository.go) |
| `[dbo].[usp_go_HistorialExamenLaboratorioResultado]` | Plan / Resultados Lab | Detalle analítico de resultados de laboratorio con rangos referenciales. | [`resultado_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/resultado_repository.go) |
| `[dbo].[usp_go_SelectHistorialExamenImageneologia]` | Plan / Resultados Imagen | Historial de estudios radiológicos y ecográficos del paciente. | [`resultado_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/resultado_repository.go) |
| `[dbo].[usp_go_InterconsultaCrear]` | Plan / Interconsultas | Registra una nueva solicitud de interconsulta médica. | [`interconsulta_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/interconsulta_repository.go) |
| `[dbo].[usp_go_RecetaCabeceraAgregar]` | Órdenes / Recetas | Inserta la cabecera de la receta médica en el punto de carga. | [`orden_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/orden_repository.go) |
| `[dbo].[usp_go_RecetaDetalleAgregar]` | Órdenes / Recetas | Inserta cada ítem farmacológico con dosis, duración y precio. | [`orden_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/orden_repository.go) |
| `[dbo].[usp_go_PacientesDatosAdicionalesIdPaciente]` | Antecedentes Médicos | Consulta antecedentes patológicos, quirúrgicos y comorbilidades. | [`patient_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/patient_repository.go) |
| `[dbo].[usp_go_Cat_ParametroClinico_Listar]` | Catálogos / Parámetros Clínicos | Lista parámetros tipificados por grupo: Estado General (1), Hidratación (2), Conciencia/Glasgow (3), Destino/Alta (4) y Tipo Atención (5). | [`catalog_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/catalog_repository.go) |


---

## 2. Procedimientos Almacenados Antiguos que ya fueron Reemplazados

Los siguientes procedimientos antiguos estaban cableados en el código y han sido reemplazados por sus equivalentes `usp_go_`:

| SP Antiguo Reemplazado | SP Moderno `usp_go_` que lo sustituye | Archivo Modificado |
| :--- | :--- | :--- |
| `sp_go_InsertarEvolucionSintomas` | `[dbo].[usp_go_InsertarEvolucionSintomas]` | [`sintoma_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/sintoma_repository.go) |
| `Usp_SelectRecetaFrecuenciaSelecionarTodos` | `[dbo].[usp_go_SelectRecetaFrecuenciaSelecionarTodos]` | [`catalog_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/catalog_repository.go) |
| `Usp_SelectRecetaUndDosisSelecionarTodos` | `[dbo].[Usp_go_SelectRecetaUndDosisSelecionarTodos]` | [`catalog_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/catalog_repository.go) |
| `RecetasListadoViasAdministracion` | `[dbo].[usp_go_RecetasListadoViasAdministracion]` | [`catalog_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/catalog_repository.go) |
| `InterconsultaFiltrrMedicoXIdEspecialidad` | `[dbo].[usp_go_MedicosFiltrarPorIdEspecialidad]` | [`interconsulta_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/interconsulta_repository.go) |
| `webEspecialidadesListarInterConsulta` | `[dbo].[usp_go_EspecialidadesXidTipoServicio]` | [`interconsulta_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/interconsulta_repository.go) |
| `sp_go_AgregarSintomaCatalogo` | `[dbo].[usp_go_AgregarSintomaCatalogo]` | [`sintoma_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/sintoma_repository.go) |
| `sp_go_SelectSintomaCatalogo` | `[dbo].[usp_go_SelectSintomaCatalogo]` | [`sintoma_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/sintoma_repository.go) |
| `usp_selectInformeImagenes` | `[dbo].[usp_go_SelectInformeImagenes]` | [`resultado_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/resultado_repository.go) |
| `Web_sp_GuardarEvolucionFirma` *(Base64 en AtencionesFirma)* | `[dbo].[usp_go_EvolucionesMedicas_Insertar]` *(con `@UbicacionArchivo` y `@EstadoFirma`)* | [`evolucion_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/evolucion_repository.go) |
| `webEvolucionesFirmaListarIdRegAtencion` | `[dbo].[usp_go_EvolucionesMedicas_Bandeja]` | [`evolucion_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/evolucion_repository.go) |
| `CatalogoBienesInsumosHospSeleccionarXIdProducto` | `[dbo].[FactCatalogoBienesInsumos]` / `[dbo].[usp_go_SelectMedicamentosFiltro]` | [`orden_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/orden_repository.go) |
| `webEvolucionMotivoListar` | `[dbo].[usp_go_EvolucionesMedicas_Bandeja]` *(campo `Motivo`)* | [`motivo_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/motivo_repository.go) |
| `usp_Tab_Atencion_Registro_Guardar` | `[dbo].[usp_go_EvolucionesMedicas_Insertar]` | [`motivo_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/motivo_repository.go) |
| `PacientesDatosAdicionalesModificar` | `[dbo].[usp_go_PacientesDatosAdicionalesModificar]` | [`patient_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/patient_repository.go) |
| `PacientesDatosAdicionalesAgregar` | `[dbo].[usp_go_PacientesDatosAdicionalesAgregar]` | [`patient_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/patient_repository.go) |
| `webPacientesDatosAdicionalesActualizarComorbilidades` | `[dbo].[usp_go_PacientesDatosAdicionalesActualizarComorbilidades]` | [`patient_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/patient_repository.go) |
| `AtencionesSeleccionarPorIdAtencion` | `[dbo].[Atenciones]` *(query directo parametrizado)* / `[dbo].[usp_go_SelectEvolucionMedicaId]` | [`orden_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/orden_repository.go) |
| `MedicosXidEmpleado` | `[dbo].[usp_go_MedicosXidEmpleado]` | [`orden_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/orden_repository.go) |
| `WebListarInterconsultasPorAtencion` | `[dbo].[usp_go_ListarInterconsultasPorAtencion]` | [`interconsulta_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/interconsulta_repository.go) |
| `webOrdenesListarIdCuentaAtencion` | `[dbo].[usp_go_OrdenesListarIdCuentaAtencion]` | [`orden_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/orden_repository.go) |
| `webParametrosDatosInstitucion` | `[dbo].[usp_go_ParametrosDatosInstitucion]` | [`catalog_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/catalog_repository.go) |
| `EmpleadosSeleccionarPorIdEmpleado` | `[dbo].[usp_go_EmpleadosSeleccionarPorId]` | [`auth_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/auth_repository.go) |
| `webMenuSeleccionarIdEmpleado` | `[dbo].[usp_go_MenuSeleccionarPorIdEmpleado]` | [`auth_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/auth_repository.go) |
| `web_MenuPermisosIdempleado` | `[dbo].[usp_go_MenuPermisosPorIdEmpleado]` | [`auth_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/auth_repository.go) |
| `webMedicosXidEmpleado` | `[dbo].[usp_go_MedicosDetallePorIdEmpleado]` | [`auth_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/auth_repository.go) |
| `WebListarInterconsultasSegunTipoServicio` | `[dbo].[usp_go_ListarInterconsultasSegunTipoServicio]` | [`interconsulta_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/interconsulta_repository.go) |
| `webListarInterconsultaPorId` | `[dbo].[usp_go_ListarInterconsultaPorId]` | [`interconsulta_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/interconsulta_repository.go) |
| `webInterconsultaActualizaEstadoFirmado` | `[dbo].[usp_go_InterconsultaActualizaEstadoFirmado]` | [`interconsulta_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/interconsulta_repository.go) |
| `webInterconsultaConsultarFirma` | `[dbo].[usp_go_InterconsultaConsultarFirma]` / `[dbo].[usp_go_InterconsultaActualizaEstadoFirmado]` | [`interconsulta_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/interconsulta_repository.go) |

---

## 3. Estado Global de Migración: ¡100% COMPLETADO (0 Pendientes)! 🚀

### 🗺️ Mapa Visual de Navegación en Pantalla (`/hospitalizacion`)

Cuando el usuario ingresa a **Hospitalización** en el menú principal y hace clic en un paciente o en el botón **"Nueva Evolución"**, accede al **Formulario Clínico SOAP**. En la parte izquierda de la pantalla se despliega un **Sidebar de 10 Pasos**:

| Paso | Nombre en Pantalla | ¿Qué ve y hace el médico en este paso? | Estado de Migración |
| :---: | :--- | :--- | :---: |
| **01** | **Información general** | Consulta datos del paciente y edita **Antecedentes / Comorbilidades** con el botón `[ ✏️ Editar ]`. | ✅ **100% `usp_go_`** |
| **02** | **Apreciación clínica** | Selecciona el **Tipo de motivo**, redacta la **Descripción detallada** y avanza con `[ Siguiente → ]`. | ✅ **100% `usp_go_`** |
| **03** | **Subjetivo** | Marca los **Síntomas referidos** por el paciente organizados por sistema y califica la escala del dolor (EVA 0–10). | ✅ **100% `usp_go_`** |
| **04** | **Objetivo** | Ingresa **Signos Vitales** (PA, FC, FR, SatO2, Temp) y el Examen Físico estructurado por aparatos. | ✅ **100% `usp_go_`** |
| **05** | **Resultados** | Visualiza el historial de análisis de **Laboratorio Clínico** y reportes de **Rayos X / Ecografías / Imágenes**. | ✅ **100% `usp_go_`** |
| **06** | **Evaluación** | Asocia los diagnósticos con código **CIE-10** (Presuntivo, Definitivo, Repetido) y redacta la evolución libre. | ✅ **100% `usp_go_`** |
| **07** | **Plan** | Prescribe el **Tratamiento Farmacológico (Recetas/Medicamentos)**, solicita **Exámenes Auxiliares** y emite **Interconsultas**. | ✅ **100% `usp_go_`** |
| **08** | **Alta médica** | Define la orden de destino hospitalario (Alta médica, Ingreso a piso, Observación en cama, Traslado). | ✅ **100% `usp_go_`** |
| **09** | **Adjuntos** | Carga y previsualiza documentos complementarios, consentimientos o imágenes clínicas. | ✅ **100% `usp_go_`** |
| **10** | **Firma y auditoría** | Presiona el botón **"Firmar evolución con DNIe"**, abriendo el modal con el **PDF oficial y Firma Perú**. | ✅ **100% `usp_go_`** |

---

### 🎉 GRUPO 1: FLUJO CENTRAL DE EVOLUCIÓN MÉDICA — 100% COMPLETADO (0 Pendientes)
*Todos los procedimientos almacenados de los 10 pasos del formulario SOAP, de la generación del documento médico legal, de la autenticación de usuarios y de los permisos han sido modernizados, validados y conectados exitosamente.*

---

### 🎉 GRUPO 2: FLUJO EXTERNO DEL ESPECIALISTA INTERCONSULTADO — 100% COMPLETADO (0 Pendientes)
*Todos los procedimientos almacenados para la bandeja del servicio de destino, apertura por ID, actualización de estado firmado y auditoría de firma digital han sido modernizados a `usp_go_*` y conectados en el backend.*

---

## 4. Catálogo Completo de los 90 Procedimientos Almacenados `usp_go_` en la BD

Verificados directamente en la base de datos SQL Server (`SIGH`):

1. `usp_go_AgregarSintomaCatalogo`
2. `usp_go_AtencionesDiagnosticosAgregar`
3. `usp_go_AtencionesDiagnosticosEliminar`
4. `usp_go_AtencionesDiagnosticosListar`
5. `usp_go_BuscarCatalogoExamenes`
6. `usp_go_Cat_ParametroClinico_Listar`
7. `usp_go_Cat_Prioridad_PrimeraAtencion`
8. `usp_go_EmpleadosSeleccionarPorId`
9. `usp_go_EspecialidadesXidTipoServicio`
10. `usp_go_EvolucionesMedicas_Bandeja`
11. `usp_go_EvolucionesMedicas_Insertar`
12. `usp_go_EvolucionMotivoListar`
13. `usp_go_GrabarListaEsperaQx`
14. `usp_go_HistorialExamenLaboratorio`
15. `usp_go_HistorialExamenLaboratorioResultado`
16. `usp_go_InsertarEvolucionSintomas`
17. `usp_go_InsertarPrimeraAtencion`
18. `usp_go_InterconsultaActualizaEstadoFirmado`
19. `usp_go_InterconsultaConsultarFirma`
20. `usp_go_InterconsultaCrear`
21. `usp_go_LabMovimientoAgregar`
22. `usp_go_LabMovimientoLaboratorioAgregar`
23. `usp_go_ListaEsperaQxListar`
24. `usp_go_ListaEsperaQxReporte`
25. `usp_go_ListarCentrosPoblados`
26. `usp_go_ListarDatosMedicoxIdEmpleado`
27. `usp_go_ListarDepartamentos`
28. `usp_go_ListarDiagnosticos`
29. `usp_go_ListarDistritos`
30. `usp_go_ListarEspecialidades`
31. `usp_go_ListarEspecialidadesQx`
32. `usp_go_ListarEspecialidadXDepartamento`
33. `usp_go_ListarEstadosCivil`
34. `usp_go_listarEstadosLlegoPaciente`
35. `usp_go_ListarFuentesFinanciamiento`
36. `usp_go_ListarGradoInstruccion`
37. `usp_go_ListarInterconsultaPorId`
38. `usp_go_ListarInterconsultasPorAtencion`
39. `usp_go_ListarInterconsultasSegunTipoServicio`
40. `usp_go_ListarMedicosListaEspera`
41. `usp_go_ListarOcupaciones`
42. `usp_go_ListarPacienteListaEspera`
43. `usp_go_ListarPacientePorNroDocyTipo`
44. `usp_go_listarpacientes`
45. `usp_go_ListarPacientesIngresoEmergencia`
46. `usp_go_ListarPacientesSegunTipoServicio`
47. `usp_go_ListarPaises`
48. `usp_go_ListarProvincias`
49. `usp_go_ListarServicios`
50. `usp_go_ListarTiposDocumentos`
51. `usp_go_ListarTiposSexos`
52. `usp_go_Login`
53. `usp_go_MedicosDetallePorIdEmpleado`
54. `usp_go_MedicosFiltrarPorIdEspecialidad`
55. `usp_go_MedicosXidEmpleado`
56. `usp_go_MenuPermisosPorIdEmpleado`
57. `usp_go_MenuSeleccionarPorIdEmpleado`
58. `usp_go_ModificarListaEsperaQx`
59. `usp_go_ModificarPaciente`
60. `usp_go_ObtenerAtencionSintomas`
61. `usp_go_ObtenerDiagnosticosAtencion`
62. `usp_go_OrdenesListarIdCuentaAtencion`
63. `usp_go_OrdenesListarPorCuentaAtencion`
64. `usp_go_PacientesDatosAdicionalesActualizarComorbilidades`
65. `usp_go_PacientesDatosAdicionalesAgregar`
66. `usp_go_PacientesDatosAdicionalesIdPaciente`
67. `usp_go_PacientesDatosAdicionalesModificar`
68. `usp_go_ParametrosDatosInstitucion`
69. `usp_go_ParametrosFTP_Obtener`
70. `usp_go_ProcedimentosFiltarNombrePuntoCarga`
71. `usp_go_RecetaCabeceraAgregar`
72. `usp_go_RecetaDetalleAgregar`
73. `usp_go_RecetaDosisSelecionarTodos`
74. `usp_go_RecetasListadoViasAdministracion`
75. `usp_go_ResultadoExamenLaboratorio`
76. `usp_go_SelectDiagnosticos`
77. `usp_go_SelectEvolucionMedicaId`
78. `usp_go_SelectHistoriaLaboratorio`
79. `usp_go_SelectHistorialExamenImageneologia`
80. `usp_go_SelectInformeImagenes`
81. `usp_go_SelectMedicamentosFiltro`
82. `usp_go_SelectRecetaFrecuenciaSelecionarTodos`
83. `Usp_go_SelectRecetaUndDosisSelecionarTodos`
84. `usp_go_SelectSintomaCatalogo`
85. `usp_go_Tab_Atencion_Registro_Guardar`
86. `usp_go_WebAtencionesDiagnosticosAgregar`
87. `usp_go_webFUAgregar`
88. `usp_go_webParametroSeleccionarPorId`
89. `usp_go_webRecetaCabeceraAgregar`

---

## 5. Procedimientos Almacenados Utilizados para Obtener Servicio y Especialidad

A continuación se detallan los procedimientos almacenados específicos utilizados en la arquitectura para obtener la información de **Servicio** y **Especialidad**, tanto a nivel del **Paciente / Atención** como del **Médico Tratante / Catálogos**:

### A. Para la Atención y Paciente (Bandeja y Evolución)
- **`[dbo].[usp_go_ListarPacientesSegunTipoServicio]`**
  - **Parámetros**: `@IdTipoServicio INT, @Fecha DATETIME, @Filtro VARCHAR(250), @IdUsuario INT`
  - **Columnas Proyectadas**: `IdEpisodio`, `IdAtencion`, `IdPaciente`, `IdCuentaAtencion`, `Paciente`, `NroHistoriaClinica`, **`Servicio`** (ej. `"Hospitalización - GINECOLOGIA"`), **`Especialidad`** (asociada al servicio, ej. `"GINECOLOGIA Y OBSTETRICIA"`), `Sexo`, `Cama`, `Edad`.
  - **Uso en Backend**: [`evolucion_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/evolucion_repository.go) en el método `ListPatients`. Alimenta la lista lateral, la ficha del paciente, la cabecera del formulario SOAP y el documento PDF.

### B. Para el Médico Tratante / Operador Clínico
- **`[dbo].[usp_go_MedicosDetallePorIdEmpleado]`**
  - **Parámetros**: `@IdEmpleado INT`
  - **Columnas Proyectadas**: `IdMedico`, `IdEmpleado`, `Colegiatura` (CMP), `RNE`, `idColegioHIS`, `LoteHIS`, `RneEspecialidad`, `RneEstado`, `egresado`, `Firma`, **`Especialidad`** (obtenida mediante `JOIN dbo.MedicosEspecialidad` y `dbo.Especialidades`).
  - **Uso en Backend**: [`auth_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/auth_repository.go) en el endpoint `/api/v1/auth/perfil`. Alimenta el modal "Mi Perfil" en la barra superior, la firma digital y los metadatos del médico en el PDF.

### C. Para Catálogos Maestros e Interconsultas
- **`[dbo].[usp_go_ListarServicios]`**
  - **Parámetros**: `@IdTipoServicio INT`
  - **Columnas Proyectadas**: `IdServicio`, `Nombre`.
  - **Uso**: Catálogo maestro de servicios activos según el tipo de servicio (Emergencia, Hospitalización, Consulta Externa).
- **`[dbo].[usp_go_ListarEspecialidades]`**
  - **Parámetros**: Ninguno.
  - **Columnas Proyectadas**: `id` (`IdEspecialidad`), `descripcion` (`Nombre`).
  - **Uso**: Catálogo general de especialidades hospitalarias activas.
- **`[dbo].[usp_go_EspecialidadesXidTipoServicio]`**
  - **Parámetros**: `@IdTipoServicio INT`
  - **Columnas Proyectadas**: `IdEspecialidad`, `Nombre`, `DescripcionLarga`.
  - **Uso en Backend**: [`interconsulta_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/interconsulta_repository.go). Lista las especialidades de destino disponibles para emitir interconsultas médicas.
- **`[dbo].[usp_go_MedicosFiltrarPorIdEspecialidad]`**
  - **Parámetros**: `@IdEspecialidad INT`
  - **Columnas Proyectadas**: `IdMedico`, `NombreMedico`, `Colegiatura`.
  - **Uso en Backend**: [`interconsulta_repository.go`](file:///c:/willian/galenos_pro/sgp_backend/internal/adapters/output/persistence/sqlserver/interconsulta_repository.go). Filtra los médicos especialistas disponibles a quienes asignar una interconsulta.
- **`[dbo].[usp_go_ListarEspecialidadXDepartamento]`**
  - **Parámetros**: `@IdDepartamento INT`
  - **Columnas Proyectadas**: `IdEspecialidad`, `Nombre`.
  - **Uso**: Lista las especialidades correspondientes a un departamento médico hospitalario específico.

