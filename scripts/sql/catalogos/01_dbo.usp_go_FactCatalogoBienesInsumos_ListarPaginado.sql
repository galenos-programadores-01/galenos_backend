-- ============================================================================
-- Procedimiento: dbo.usp_go_FactCatalogoBienesInsumos_ListarPaginado
-- Módulo: Catálogos / Bienes e Insumos Médicos (HNSEB)
-- Objetivo: Listar el catálogo de fármacos y materiales médicos con paginación
--           nativa en el motor, filtros de búsqueda, petitorio y grupos.
-- Base de Datos: SIGH / Galenos Pro
-- Motor: SQL Server 2016 / 2019 / 2022
-- ============================================================================

SET ANSI_NULLS ON;
GO
SET QUOTED_IDENTIFIER ON;
GO

CREATE OR ALTER PROCEDURE dbo.usp_go_FactCatalogoBienesInsumos_ListarPaginado
    @Filtro VARCHAR(150) = NULL,
    @IdGrupoFarmacologico INT = 0,
    @TipoProducto INT = 0,
    @SoloPetitorio BIT = NULL,
    @SoloIndicadorGestion BIT = NULL,
    @Estado INT = 1,
    @PageNumber INT = 1,
    @PageSize INT = 20
AS
BEGIN
    SET NOCOUNT ON;
    SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;

    BEGIN TRY
        IF @PageNumber < 1 SET @PageNumber = 1;
        IF @PageSize < 1 SET @PageSize = 20;

        DECLARE @FiltroLike VARCHAR(160) = NULL;
        IF @Filtro IS NOT NULL AND LTRIM(RTRIM(@Filtro)) <> ''
            SET @FiltroLike = '%' + LTRIM(RTRIM(@Filtro)) + '%';

        -- Cálculo del total de registros
        DECLARE @TotalRegistros INT = 0;

        SELECT @TotalRegistros = COUNT(1)
        FROM dbo.FactCatalogoBienesInsumos b WITH (NOLOCK)
        WHERE (@FiltroLike IS NULL OR (
                b.Nombre LIKE @FiltroLike 
                OR b.Codigo LIKE @FiltroLike 
                OR b.CodigoSig LIKE @FiltroLike 
                OR b.NombreComercial LIKE @FiltroLike
              ))
          AND (@IdGrupoFarmacologico = 0 OR b.IdGrupoFarmacologico = @IdGrupoFarmacologico)
          AND (@TipoProducto = 0 OR b.TipoProducto = @TipoProducto)
          AND (@SoloPetitorio IS NULL OR b.Petitorio = @SoloPetitorio)
          AND (@SoloIndicadorGestion IS NULL OR ISNULL(b.ProductoIndicadorGestion, 0) = @SoloIndicadorGestion)
          AND (@Estado IS NULL OR b.Estado = @Estado);

        -- Retorno de la página solicitada
        SELECT 
            b.IdProducto,
            LTRIM(RTRIM(ISNULL(b.Codigo, ''))) AS Codigo,
            LTRIM(RTRIM(ISNULL(b.CodigoSig, ''))) AS CodigoSiga,
            LTRIM(RTRIM(ISNULL(b.CodigoSunat, ''))) AS CodigoSunat,
            LTRIM(RTRIM(ISNULL(b.Nombre, ''))) AS Nombre,
            LTRIM(RTRIM(ISNULL(b.NombreComercial, ''))) AS NombreComercial,
            ISNULL(b.IdGrupoFarmacologico, 0) AS IdGrupoFarmacologico,
            LTRIM(RTRIM(ISNULL(g.Descripcion, 'DESCONOCIDO'))) AS GrupoFarmacologico,
            ISNULL(b.PrecioCompra, 0) AS PrecioCompra,
            ISNULL(b.PrecioDistribucion, 0) AS PrecioDistribucion,
            ISNULL(b.PrecioUltCompra, 0) AS PrecioUltimaCompra,
            ISNULL(b.StockMinimo, 0) AS StockMinimo,
            ISNULL(b.TipoProducto, 0) AS TipoProducto,
            LTRIM(RTRIM(ISNULL(b.Concentracion, ''))) AS Concentracion,
            LTRIM(RTRIM(ISNULL(b.Presentacion, ''))) AS Presentacion,
            LTRIM(RTRIM(ISNULL(b.FormaFarmaceutica, ''))) AS FormaFarmaceutica,
            CAST(ISNULL(b.Petitorio, 0) AS BIT) AS Petitorio,
            ISNULL(b.Estado, 0) AS Estado,
            CAST(CASE WHEN ISNULL(b.Vitales, 0) = 1 THEN 1 ELSE 0 END AS BIT) AS EsVital,
            CAST(CASE WHEN ISNULL(b.Estrategicos, 0) = 1 THEN 1 ELSE 0 END AS BIT) AS EsEstrategico,
            CAST(CASE WHEN ISNULL(b.Oncologicos, 0) = 1 THEN 1 ELSE 0 END AS BIT) AS EsOncologico,
            CAST(CASE WHEN ISNULL(b.CargaEnFUA, 0) = 1 THEN 1 ELSE 0 END AS BIT) AS CargaEnFua,
            CAST(CASE WHEN ISNULL(b.ProductoIndicadorGestion, 0) = 1 THEN 1 ELSE 0 END AS BIT) AS EsIndicadorGestion,
            @TotalRegistros AS TotalRegistros
        FROM dbo.FactCatalogoBienesInsumos b WITH (NOLOCK)
        LEFT JOIN dbo.FactInsumosGrupoFarmacologico g WITH (NOLOCK) 
            ON g.IdGrupoFarmacologico = b.IdGrupoFarmacologico
        WHERE (@FiltroLike IS NULL OR (
                b.Nombre LIKE @FiltroLike 
                OR b.Codigo LIKE @FiltroLike 
                OR b.CodigoSig LIKE @FiltroLike 
                OR b.NombreComercial LIKE @FiltroLike
              ))
          AND (@IdGrupoFarmacologico = 0 OR b.IdGrupoFarmacologico = @IdGrupoFarmacologico)
          AND (@TipoProducto = 0 OR b.TipoProducto = @TipoProducto)
          AND (@SoloPetitorio IS NULL OR b.Petitorio = @SoloPetitorio)
          AND (@SoloIndicadorGestion IS NULL OR ISNULL(b.ProductoIndicadorGestion, 0) = @SoloIndicadorGestion)
          AND (@Estado IS NULL OR b.Estado = @Estado)
        ORDER BY b.IdProducto DESC
        OFFSET (@PageNumber - 1) * @PageSize ROWS
        FETCH NEXT @PageSize ROWS ONLY;

    END TRY
    BEGIN CATCH
        DECLARE @ErrMsg NVARCHAR(4000) = ERROR_MESSAGE();
        DECLARE @ErrSeverity INT = ERROR_SEVERITY();
        DECLARE @ErrState INT = ERROR_STATE();
        RAISERROR(@ErrMsg, @ErrSeverity, @ErrState);
    END CATCH
END;
GO
