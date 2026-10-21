SET ANSI_NULLS ON;
GO
SET QUOTED_IDENTIFIER ON;
GO

CREATE OR ALTER PROCEDURE dbo.usp_go_AtencionesCE_ObtenerFicha
    @IdAtencion INT
AS
BEGIN
    SET NOCOUNT ON;
    SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;

    BEGIN TRY
        IF @IdAtencion IS NULL OR @IdAtencion <= 0
        BEGIN
            RAISERROR('El identificador de atención es requerido.', 16, 1);
            RETURN;
        END;

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
        WHERE a.IdAtencion = @IdAtencion;

    END TRY
    BEGIN CATCH
        DECLARE @MensajeError NVARCHAR(4000) = ERROR_MESSAGE();
        DECLARE @SeveridadError INT = ERROR_SEVERITY();
        DECLARE @EstadoError INT = ERROR_STATE();
        RAISERROR(@MensajeError, @SeveridadError, @EstadoError);
    END CATCH
END;
GO
