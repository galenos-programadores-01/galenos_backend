CREATE OR ALTER PROCEDURE dbo.usp_go_ActualizarListaEsperaQx
    @IdListaEspera INT
AS
BEGIN
    SET NOCOUNT ON;

    UPDATE dbo.ListaEsperaQx
    SET Tipo = 'HOSPITALIZADO'
    WHERE Id = @IdListaEspera;
END;
GO
