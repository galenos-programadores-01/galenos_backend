CREATE OR ALTER PROCEDURE dbo.usp_go_ListarInterconsultasPorAtencion
    @IdAtencion INT
AS
BEGIN
    SET NOCOUNT ON;

    SELECT
        i.idInterconsulta,
        ISNULL(a.IdAtencion, i.IdRegAtencion) AS IdAtencionOrigen,
        ISNULL(i.IdEspecialidadSolicitada, 0) AS IdEspecialidad,
        ISNULL(i.idEmpleadoRespuesta, 0)      AS IdMedicoDestino,
        ISNULL(i.MotivoClinico, '')           AS Motivo,
        i.fechaSolicitud                      AS FechaSolicitud,
        CASE WHEN i.idEstado = 1 THEN 'SOLICITUD' ELSE 'RESPUESTA' END AS Estado,
        ISNULL(esp.Nombre, '')                AS Especialidad,
        ISNULL(i.ResumenClinico, ISNULL(i.ResumenHC, '')) AS ResumenHC
    FROM dbo.Interconsulta i
    LEFT JOIN dbo.Tab_Atencion_Registro re ON re.IdRegAtencion = i.IdRegAtencion
    LEFT JOIN dbo.EpisodioAtencion ea       ON ea.IdEpisodio = re.IdEpisodio
    LEFT JOIN dbo.Atenciones a              ON a.IdAtencion = ea.IdAtencion
    LEFT JOIN dbo.Especialidades esp        ON esp.IdEspecialidad = i.IdEspecialidadSolicitada
    WHERE a.IdAtencion = @IdAtencion
       OR i.IdRegAtencion = @IdAtencion
       OR ea.IdEpisodio = @IdAtencion
    ORDER BY i.fechaSolicitud;
END;
GO
