CREATE OR ALTER PROCEDURE dbo.usp_go_FactCatalogoBienesInsumos_ObtenerPrecioDistribucion
    @IdProducto INT
AS
BEGIN
    SET NOCOUNT ON;

    SELECT ISNULL(PrecioDistribucion, 0) AS PrecioDistribucion
    FROM dbo.FactCatalogoBienesInsumos WITH (NOLOCK)
    WHERE IdProducto = @IdProducto;
END;
GO
