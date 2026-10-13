SET ANSI_NULLS ON;
GO
SET QUOTED_IDENTIFIER ON;
GO

CREATE OR ALTER PROCEDURE dbo.usp_go_AtencionesCE_ListarCola
    @IdMedico INT = 0,
    @IdServicio INT = 0,
    @Fecha DATE = NULL,
    @Filtro VARCHAR(100) = NULL
AS
BEGIN
    SET NOCOUNT ON;
    SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;

    BEGIN TRY
        DECLARE @FiltroLimpio VARCHAR(100) = NULL;
        IF @Filtro IS NOT NULL AND LTRIM(RTRIM(@Filtro)) <> ''
            SET @FiltroLimpio = '%' + LTRIM(RTRIM(@Filtro)) + '%';

        DECLARE @FechaConsulta DATE = @Fecha;
        IF @FechaConsulta IS NULL AND @FiltroLimpio IS NULL
            SET @FechaConsulta = CAST(GETDATE() AS DATE);

        IF (@FechaConsulta IS NOT NULL)
        BEGIN
            DECLARE @FechaInicio DATETIME = CAST(@FechaConsulta AS DATETIME);
            DECLARE @FechaFin DATETIME = DATEADD(DAY, 1, @FechaInicio);

            SELECT 
                a.IdAtencion,
                a.IdCuentaAtencion,
                a.IdPaciente,
                a.FechaIngreso,
                a.HoraIngreso,
                a.HoraEgreso,
                a.Edad,
                a.IdServicioIngreso,
                a.IdMedicoIngreso,
                a.idEstadoAtencion,
                a.idFuenteFinanciamiento
            INTO #AtencionesDia
            FROM dbo.Atenciones a WITH (NOLOCK, INDEX(IX_Atenciones_Fecha_Tipo))
            WHERE a.IdTipoServicio = 1
              AND a.FechaIngreso >= @FechaInicio AND a.FechaIngreso < @FechaFin
              AND (@IdMedico = 0 OR a.IdMedicoIngreso = @IdMedico)
              AND (@IdServicio = 0 OR a.IdServicioIngreso = @IdServicio);

            CREATE CLUSTERED INDEX CIX_TempAtenciones ON #AtencionesDia(IdAtencion);

            SELECT 
                a.IdAtencion,
                ISNULL(a.IdCuentaAtencion, 0) AS IdCuentaAtencion,
                ISNULL(a.IdPaciente, 0) AS IdPaciente,
                ISNULL(p.NroHistoriaClinica, 0) AS NroHistoriaClinica,
                LTRIM(RTRIM(ISNULL(p.NroDocumento, ''))) AS NroDocumento,
                LTRIM(RTRIM(CONCAT(
                    ISNULL(p.ApellidoPaterno, ''), ' ', 
                    ISNULL(p.ApellidoMaterno, ''), ' ', 
                    ISNULL(p.PrimerNombre, ''), ' ', 
                    ISNULL(p.SegundoNombre, '')
                ))) AS PacienteNombre,
                CONVERT(VARCHAR(10), a.FechaIngreso, 120) AS FechaIngreso,
                ISNULL(a.HoraIngreso, '') AS HoraIngreso,
                ISNULL(a.HoraEgreso, '') AS HoraEgreso,
                ISNULL(a.Edad, 0) AS Edad,
                ISNULL(a.IdServicioIngreso, 0) AS IdServicioIngreso,
                LTRIM(RTRIM(ISNULL(s.Nombre, ''))) AS ServicioNombre,
                ISNULL(a.IdMedicoIngreso, 0) AS IdMedicoIngreso,
                LTRIM(RTRIM(CONCAT(
                    ISNULL(e.ApellidoPaterno, ''), ' ', 
                    ISNULL(e.ApellidoMaterno, ''), ' ', 
                    ISNULL(e.Nombres, '')
                ))) AS MedicoNombre,
                CASE 
                    WHEN ISNULL(a.HoraEgreso, '') <> '' OR a.idEstadoAtencion = 2 THEN 'ATENDIDO'
                    WHEN ace.idAtencion IS NOT NULL THEN 'EN_ATENCION'
                    ELSE 'PENDIENTE'
                END AS EstadoAtencion,
                ISNULL(c.IdCita, 0) AS IdCita,
                LTRIM(RTRIM(ISNULL(ff.Descripcion, 'PARTICULAR'))) AS TipoFinanciador,
                ISNULL(CASE WHEN p.IdTipoSexo = 2 THEN 'F' ELSE 'M' END, 'M') AS Sexo,
                ISNULL(NULLIF(LTRIM(RTRIM(p.celular)), ''), ISNULL(LTRIM(RTRIM(p.Telefono)), '')) AS Telefono,
                CASE WHEN atri.IdAtencion IS NOT NULL THEN 1 ELSE 0 END AS TieneTriaje,
                CASE WHEN a.HoraIngreso < '13:00' THEN 'MAÑANA' ELSE 'TARDE' END AS Turno
            FROM #AtencionesDia a
            LEFT JOIN dbo.Pacientes p WITH (NOLOCK) ON a.IdPaciente = p.IdPaciente
            LEFT JOIN dbo.Servicios s WITH (NOLOCK) ON a.IdServicioIngreso = s.IdServicio
            LEFT JOIN dbo.Citas c WITH (NOLOCK) ON c.IdAtencion = a.IdAtencion
            LEFT JOIN dbo.Medicos m WITH (NOLOCK) ON a.IdMedicoIngreso = m.IdMedico
            LEFT JOIN dbo.Empleados e WITH (NOLOCK) ON m.IdEmpleado = e.IdEmpleado
            LEFT JOIN dbo.FuentesFinanciamiento ff WITH (NOLOCK) ON a.idFuenteFinanciamiento = ff.IdFuenteFinanciamiento
            LEFT JOIN dbo.AtencionesCE ace WITH (NOLOCK) ON ace.idAtencion = a.IdAtencion
            LEFT JOIN dbo.AtencionesTriaje atri WITH (NOLOCK) ON atri.IdAtencion = a.IdAtencion
            WHERE (@FiltroLimpio IS NULL OR (
                p.NroDocumento LIKE @FiltroLimpio
                OR CAST(p.NroHistoriaClinica AS VARCHAR(20)) LIKE @FiltroLimpio
                OR p.ApellidoPaterno LIKE @FiltroLimpio
                OR p.PrimerNombre LIKE @FiltroLimpio
                OR CONCAT(p.ApellidoPaterno, ' ', p.ApellidoMaterno, ' ', p.PrimerNombre) LIKE @FiltroLimpio
            ))
            ORDER BY a.HoraIngreso ASC, a.IdAtencion ASC;

            DROP TABLE #AtencionesDia;
        END
        ELSE
        BEGIN
            SELECT TOP 200
                a.IdAtencion,
                ISNULL(a.IdCuentaAtencion, 0) AS IdCuentaAtencion,
                ISNULL(a.IdPaciente, 0) AS IdPaciente,
                ISNULL(p.NroHistoriaClinica, 0) AS NroHistoriaClinica,
                LTRIM(RTRIM(ISNULL(p.NroDocumento, ''))) AS NroDocumento,
                LTRIM(RTRIM(CONCAT(
                    ISNULL(p.ApellidoPaterno, ''), ' ', 
                    ISNULL(p.ApellidoMaterno, ''), ' ', 
                    ISNULL(p.PrimerNombre, ''), ' ', 
                    ISNULL(p.SegundoNombre, '')
                ))) AS PacienteNombre,
                CONVERT(VARCHAR(10), a.FechaIngreso, 120) AS FechaIngreso,
                ISNULL(a.HoraIngreso, '') AS HoraIngreso,
                ISNULL(a.HoraEgreso, '') AS HoraEgreso,
                ISNULL(a.Edad, 0) AS Edad,
                ISNULL(a.IdServicioIngreso, 0) AS IdServicioIngreso,
                LTRIM(RTRIM(ISNULL(s.Nombre, ''))) AS ServicioNombre,
                ISNULL(a.IdMedicoIngreso, 0) AS IdMedicoIngreso,
                LTRIM(RTRIM(CONCAT(
                    ISNULL(e.ApellidoPaterno, ''), ' ', 
                    ISNULL(e.ApellidoMaterno, ''), ' ', 
                    ISNULL(e.Nombres, '')
                ))) AS MedicoNombre,
                CASE 
                    WHEN ISNULL(a.HoraEgreso, '') <> '' OR a.idEstadoAtencion = 2 THEN 'ATENDIDO'
                    WHEN ace.idAtencion IS NOT NULL THEN 'EN_ATENCION'
                    ELSE 'PENDIENTE'
                END AS EstadoAtencion,
                ISNULL(c.IdCita, 0) AS IdCita,
                LTRIM(RTRIM(ISNULL(ff.Descripcion, 'PARTICULAR'))) AS TipoFinanciador,
                ISNULL(CASE WHEN p.IdTipoSexo = 2 THEN 'F' ELSE 'M' END, 'M') AS Sexo,
                ISNULL(NULLIF(LTRIM(RTRIM(p.celular)), ''), ISNULL(LTRIM(RTRIM(p.Telefono)), '')) AS Telefono,
                CASE WHEN atri.IdAtencion IS NOT NULL THEN 1 ELSE 0 END AS TieneTriaje,
                CASE WHEN a.HoraIngreso < '13:00' THEN 'MAÑANA' ELSE 'TARDE' END AS Turno
            FROM dbo.Atenciones a WITH (NOLOCK)
            INNER JOIN dbo.Pacientes p WITH (NOLOCK) ON a.IdPaciente = p.IdPaciente
            LEFT JOIN dbo.Servicios s WITH (NOLOCK) ON a.IdServicioIngreso = s.IdServicio
            LEFT JOIN dbo.Citas c WITH (NOLOCK) ON c.IdAtencion = a.IdAtencion
            LEFT JOIN dbo.Medicos m WITH (NOLOCK) ON a.IdMedicoIngreso = m.IdMedico
            LEFT JOIN dbo.Empleados e WITH (NOLOCK) ON m.IdEmpleado = e.IdEmpleado
            LEFT JOIN dbo.FuentesFinanciamiento ff WITH (NOLOCK) ON a.idFuenteFinanciamiento = ff.IdFuenteFinanciamiento
            LEFT JOIN dbo.AtencionesCE ace WITH (NOLOCK) ON ace.idAtencion = a.IdAtencion
            LEFT JOIN dbo.AtencionesTriaje atri WITH (NOLOCK) ON atri.IdAtencion = a.IdAtencion
            WHERE a.IdTipoServicio = 1
              AND (
                    p.NroDocumento LIKE @FiltroLimpio
                    OR CAST(p.NroHistoriaClinica AS VARCHAR(20)) LIKE @FiltroLimpio
                    OR p.ApellidoPaterno LIKE @FiltroLimpio
                    OR p.PrimerNombre LIKE @FiltroLimpio
                  )
              AND (@IdMedico = 0 OR a.IdMedicoIngreso = @IdMedico)
              AND (@IdServicio = 0 OR a.IdServicioIngreso = @IdServicio)
            ORDER BY a.FechaIngreso DESC, a.HoraIngreso ASC;
        END

    END TRY
    BEGIN CATCH
        DECLARE @ErrMsg NVARCHAR(4000) = ERROR_MESSAGE();
        DECLARE @ErrSeverity INT = ERROR_SEVERITY();
        DECLARE @ErrState INT = ERROR_STATE();
        RAISERROR(@ErrMsg, @ErrSeverity, @ErrState);
    END CATCH
END;
GO
