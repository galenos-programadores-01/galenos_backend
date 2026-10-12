CREATE OR ALTER PROCEDURE dbo.usp_go_EvolucionesMedicas_ListarNotasPorAtencion
    @NroAtencion INT
AS
BEGIN
    SET NOCOUNT ON;

    SELECT 
        e.IdEvolucion,
        e.NroAtencion,
        ISNULL(e.UbicacionArchivo, '') AS UbicacionArchivo,
        ISNULL(e.UsuarioCreacion, 0) AS UsuarioCreacion,
        CONVERT(varchar(19), e.FechaAtencion, 120) AS FechaAtencion,
        ISNULL(e.EstadoFirma, 1) AS EstadoFirma,
        ISNULL(RTRIM(LTRIM(CONCAT(emp.ApellidoPaterno, ' ', emp.ApellidoMaterno, ' ', emp.Nombres))), '') AS MedicoNombre,
        ISNULL(e.Motivo, '') AS Motivo,
        ISNULL(e.Subjetivo, '') AS Subjetivo,
        ISNULL(e.EscalaDolor, 0) AS EscalaDolor,
        ISNULL(e.PASistolica, 0) AS PASistolica,
        ISNULL(e.PADiastolica, 0) AS PADiastolica,
        ISNULL(e.FrecuenciaCardiaca, 0) AS FrecuenciaCardiaca,
        ISNULL(e.FrecuenciaRespiratoria, 0) AS FrecuenciaRespiratoria,
        ISNULL(e.Temperatura, 0.0) AS Temperatura,
        ISNULL(e.SaturacionOxigeno, 0) AS SaturacionOxigeno,
        ISNULL(e.Peso, 0.0) AS Peso,
        ISNULL(e.Talla, 0.0) AS Talla,
        ISNULL(e.IMC, 0.0) AS IMC,
        ISNULL(e.Glicemia, 0.0) AS Glicemia,
        ISNULL(e.ExamenFisicoGeneral, '') AS ExamenFisicoGeneral,
        ISNULL(e.ExamenFisicoPiel, '') AS ExamenFisicoPiel,
        ISNULL(e.ExamenFisicoCabezaCuello, '') AS ExamenFisicoCabezaCuello,
        ISNULL(e.ExamenFisicoToraxPulmon, '') AS ExamenFisicoToraxPulmon,
        ISNULL(e.ExamenFisicoCorazon, '') AS ExamenFisicoCorazon,
        ISNULL(e.ExamenFisicoAbdomen, '') AS ExamenFisicoAbdomen,
        ISNULL(e.ExamenFisicoGenitourinario, '') AS ExamenFisicoGenitourinario,
        ISNULL(e.ExamenFisicoExtremidadesOsteomuscular, '') AS ExamenFisicoExtremidadesOsteomuscular,
        ISNULL(e.ExamenFisicoNeurologicoMental, '') AS ExamenFisicoNeurologicoMental,
        ISNULL(e.IdEstadoClinico, 0) AS IdEstadoClinico,
        ISNULL(e.IdPronostico, 0) AS IdPronostico,
        ISNULL(e.IndicacionDieta, '') AS IndicacionDieta,
        ISNULL(e.IndicacionReposo, '') AS IndicacionReposo,
        ISNULL(e.IndicacionHidratacion, '') AS IndicacionHidratacion,
        ISNULL(e.IndicacionOxigeno, '') AS IndicacionOxigeno,
        ISNULL(e.IndicacionRestriccion, '') AS IndicacionRestriccion,
        ISNULL(e.Sugerencia, '') AS Sugerencia
    FROM dbo.Tab_EvolucionesMedicas_H_E e WITH (NOLOCK)
    LEFT JOIN dbo.Empleados emp WITH (NOLOCK) ON emp.IdEmpleado = e.UsuarioCreacion
    WHERE e.NroAtencion = @NroAtencion
    ORDER BY e.FechaAtencion DESC, e.IdEvolucion DESC;
END;
GO
