-- ============================================================================
-- Procedimiento: dbo.usp_go_Recetas_ListarDetallesPorEvolucion
-- Módulo: Evoluciones Médicas / Recetas Farmacológicas (HNSEB)
-- Base de Datos: SIGH / Galenos Pro
-- ============================================================================

SET ANSI_NULLS ON;
GO
SET QUOTED_IDENTIFIER ON;
GO

CREATE OR ALTER PROCEDURE dbo.usp_go_Recetas_ListarDetallesPorEvolucion
    @IdEvolucion INT = 0,
    @IdCuentaAtencion INT = 0
AS
BEGIN
    SET NOCOUNT ON;
    SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;

    SELECT DISTINCT
        ISNULL(b.Nombre, '') AS NombreMedicamento,
        ISNULL(rd.DescripcionadicionalReceta, '') AS DescripcionAdicional,
        ISNULL(rd.CantidadPedida, 0) AS Cantidad,
        ISNULL(CAST(rd.Duracion AS VARCHAR(50)), '') AS Duracion
    FROM dbo.RecetaDetalle rd WITH (NOLOCK)
    INNER JOIN dbo.RecetaCabecera rc WITH (NOLOCK) ON rc.idReceta = rd.idReceta
    LEFT JOIN dbo.FactCatalogoBienesInsumos b WITH (NOLOCK) ON b.IdProducto = rd.idItem
    WHERE (@IdCuentaAtencion > 0 AND rc.idCuentaAtencion = @IdCuentaAtencion)
       OR (@IdEvolucion > 0 AND rc.lnIdRegistroAtencion = @IdEvolucion);
END;
GO
