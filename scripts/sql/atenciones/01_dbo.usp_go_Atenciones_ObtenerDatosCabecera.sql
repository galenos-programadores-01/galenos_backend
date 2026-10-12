CREATE OR ALTER PROCEDURE dbo.usp_go_Atenciones_ObtenerDatosCabecera
    @IdAtencion INT
AS
BEGIN
    SET NOCOUNT ON;

    DECLARE @IdPaciente INT = 0,
            @IdCuentaAtencion INT = 0,
            @IdServicioIngreso INT = 0,
            @IdTipoServicio INT = 0;

    SELECT 
        @IdPaciente = ISNULL(IdPaciente, 0),
        @IdCuentaAtencion = ISNULL(IdCuentaAtencion, 0),
        @IdServicioIngreso = ISNULL(IdServicioIngreso, 0),
        @IdTipoServicio = ISNULL(IdTipoServicio, 0)
    FROM dbo.Atenciones WITH (NOLOCK)
    WHERE IdAtencion = @IdAtencion;

    IF @IdPaciente = 0
    BEGIN
        SELECT TOP 1 @IdPaciente = ISNULL(IdEpisodio, 0)
        FROM dbo.Tab_EvolucionesMedicas_H_E WITH (NOLOCK)
        WHERE NroAtencion = @IdAtencion AND IdEpisodio > 0;
    END

    IF @IdPaciente = 0 AND EXISTS (SELECT 1 FROM dbo.Pacientes WITH (NOLOCK) WHERE IdPaciente = @IdAtencion)
    BEGIN
        SET @IdPaciente = @IdAtencion;
    END

    SELECT 
        @IdAtencion AS IdAtencion,
        @IdPaciente AS IdPaciente,
        @IdCuentaAtencion AS IdCuentaAtencion,
        @IdServicioIngreso AS IdServicioIngreso,
        @IdTipoServicio AS IdTipoServicio;
END;
GO
