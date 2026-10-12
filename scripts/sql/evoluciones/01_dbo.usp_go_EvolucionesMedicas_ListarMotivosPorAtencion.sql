CREATE OR ALTER PROCEDURE dbo.usp_go_EvolucionesMedicas_ListarMotivosPorAtencion
    @NroAtencion INT
AS
BEGIN
    SET NOCOUNT ON;

    SELECT TOP 5 
        IdEvolucion AS IdMotivo, 
        NroAtencion AS IdRegAtencion, 
        ISNULL(Motivo, '') AS Motivo, 
        FechaAtencion AS FechaRegistro
    FROM dbo.Tab_EvolucionesMedicas_H_E WITH (NOLOCK)
    WHERE NroAtencion = @NroAtencion 
      AND ISNULL(Motivo, '') <> ''
    ORDER BY FechaAtencion DESC, IdEvolucion DESC;
END;
GO
