CREATE OR ALTER PROCEDURE dbo.usp_go_InterconsultaCrear
    @IdAtencionOrigen INT,
    @IdEspecialidad INT,
    @IdMedicoDestino INT = 0,
    @Motivo VARCHAR(1000) = '',
    @IdPrioridad INT = 1,
    @IdEmpleadoSolicitante INT = 0,
    @ResumenHC VARCHAR(5000) = ''
AS
BEGIN
    SET NOCOUNT ON;

    INSERT INTO dbo.Interconsulta (
        IdRegAtencion,
        IdRegistroAtencion,
        IdEspecialidadSolicitada,
        idEspecialidadRespuesta,
        idEmpleadoRespuesta,
        idEmpleadoSolicitud,
        IdEmpleadoSolicitante,
        MotivoClinico,
        ResumenClinico,
        ResumenHC,
        fechaSolicitud,
        idEstado,
        IdPrioridad,
        FechaRegistroSistema
    ) VALUES (
        @IdAtencionOrigen,
        @IdAtencionOrigen,
        @IdEspecialidad,
        @IdEspecialidad,
        NULLIF(@IdMedicoDestino, 0),
        NULLIF(@IdEmpleadoSolicitante, 0),
        NULLIF(@IdEmpleadoSolicitante, 0),
        @Motivo,
        @ResumenHC,
        @ResumenHC,
        GETDATE(),
        1,
        @IdPrioridad,
        GETDATE()
    );

    SELECT SCOPE_IDENTITY() AS IdInterconsulta;
END;
GO
