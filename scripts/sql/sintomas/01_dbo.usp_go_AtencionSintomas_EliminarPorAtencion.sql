CREATE OR ALTER PROCEDURE dbo.usp_go_AtencionSintomas_EliminarPorAtencion
    @IdRegAtencion INT
AS
BEGIN
    SET NOCOUNT ON;

    DELETE FROM dbo.Tab_Atencion_Sintomas 
    WHERE IdRegAtencion = @IdRegAtencion;
END;
GO
