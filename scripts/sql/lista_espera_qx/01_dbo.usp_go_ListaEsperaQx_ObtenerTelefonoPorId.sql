CREATE OR ALTER PROCEDURE dbo.usp_go_ListaEsperaQx_ObtenerTelefonoPorId
    @Id INT
AS
BEGIN
    SET NOCOUNT ON;

    SELECT Telefono
    FROM dbo.ListaEsperaQx WITH (NOLOCK)
    WHERE Id = @Id;
END;
GO
