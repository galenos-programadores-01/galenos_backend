CREATE OR ALTER PROCEDURE dbo.usp_go_ListarInterconsultasPorAtencion
    @IdAtencion INT,
    @IdEpisodio INT = NULL,
    @IdRegAtencion INT = NULL,
    @IdPrimeraAtencion INT = NULL
AS
BEGIN
    SET NOCOUNT ON;

    SELECT
        i.idInterconsulta,
        ISNULL(a.IdAtencion, i.IdRegistroAtencion) AS IdAtencionOrigen,
        ISNULL(i.IdEspecialidadSolicitada, 0) AS IdEspecialidad,
        ISNULL(i.idEmpleadoRespuesta, 0)      AS IdMedicoDestino,
        ISNULL(i.MotivoClinico, '')           AS Motivo,
        i.fechaSolicitud                      AS FechaSolicitud,
        CASE WHEN i.idEstado = 1 THEN 'SOLICITUD' ELSE 'RESPUESTA' END AS Estado,
        ISNULL(esp.Nombre, '')                AS Especialidad,
        ISNULL(i.ResumenClinico, ISNULL(i.ResumenHC, '')) AS ResumenHC
    FROM dbo.Interconsulta i WITH (NOLOCK)
    LEFT JOIN dbo.Tab_Atencion_Registro re WITH (NOLOCK) ON re.IdRegAtencion = i.IdRegistroAtencion
    LEFT JOIN dbo.EpisodioAtencion ea WITH (NOLOCK)       ON ea.IdEpisodio = re.IdEpisodio
    LEFT JOIN dbo.Atenciones a WITH (NOLOCK)              ON a.IdAtencion = ea.IdAtencion
    LEFT JOIN dbo.Especialidades esp WITH (NOLOCK)        ON esp.IdEspecialidad = i.IdEspecialidadSolicitada
    WHERE a.IdAtencion = @IdAtencion
       OR i.IdRegistroAtencion = @IdAtencion
       OR i.IdPrimeraAtencion = @IdPrimeraAtencion
       OR i.IdPrimeraAtencion = @IdAtencion
       OR ea.IdEpisodio = @IdEpisodio
    ORDER BY i.fechaSolicitud ASC, i.idInterconsulta ASC;
END;
GO
