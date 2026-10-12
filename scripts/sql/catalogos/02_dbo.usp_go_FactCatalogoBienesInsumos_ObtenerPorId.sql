-- ============================================================================
-- Procedimiento: dbo.usp_go_FactCatalogoBienesInsumos_ObtenerPorId
-- Módulo: Catálogos / Bienes e Insumos Médicos (HNSEB)
-- Objetivo: Obtener la ficha completa con especificaciones técnicas, precios,
--           presupuesto, grupos y banderas asistenciales de un bien o insumo médico.
-- Base de Datos: SIGH / Galenos Pro
-- Motor: SQL Server 2016 / 2019 / 2022
-- ============================================================================

SET ANSI_NULLS ON;
GO
SET QUOTED_IDENTIFIER ON;
GO

CREATE OR ALTER PROCEDURE dbo.usp_go_FactCatalogoBienesInsumos_ObtenerPorId
    @IdProducto INT
AS
BEGIN
    SET NOCOUNT ON;
    SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;

    BEGIN TRY
        IF @IdProducto IS NULL OR @IdProducto <= 0
        BEGIN
            RAISERROR('El identificador de producto es requerido.', 16, 1);
            RETURN;
        END;

        SELECT 
            b.IdProducto,
            LTRIM(RTRIM(ISNULL(b.Codigo, ''))) AS Codigo,
            LTRIM(RTRIM(ISNULL(b.CodigoSig, ''))) AS CodigoSiga,
            LTRIM(RTRIM(ISNULL(b.CodigoSunat, ''))) AS CodigoSunat,
            LTRIM(RTRIM(ISNULL(b.Nombre, ''))) AS Nombre,
            LTRIM(RTRIM(ISNULL(b.NombreComercial, ''))) AS NombreComercial,
            ISNULL(b.IdGrupoFarmacologico, 0) AS IdGrupoFarmacologico,
            LTRIM(RTRIM(ISNULL(g.Descripcion, 'DESCONOCIDO'))) AS GrupoFarmacologico,
            ISNULL(b.IdSubGrupoFarmacologico, 0) AS IdSubGrupoFarmacologico,
            LTRIM(RTRIM(ISNULL(sg.Descripcion, 'DESCONOCIDO'))) AS SubGrupoFarmacologico,
            ISNULL(b.IdPartida, 0) AS IdPartida,
            LTRIM(RTRIM(ISNULL(p.Codigo + ' ' + p.Descripcion, ''))) AS Partida,
            ISNULL(b.IdCentroCosto, 0) AS IdCentroCosto,
            LTRIM(RTRIM(ISNULL(cc.Descripcion, ''))) AS CentroCosto,
            ISNULL(b.PrecioCompra, 0) AS PrecioCompra,
            ISNULL(b.PrecioDistribucion, 0) AS PrecioDistribucion,
            ISNULL(b.PrecioDonacion, 0) AS PrecioDonacion,
            ISNULL(b.PrecioUltCompra, 0) AS PrecioUltimaCompra,
            ISNULL(b.StockMinimo, 0) AS StockMinimo,
            ISNULL(b.TipoProducto, 0) AS TipoProducto,
            LTRIM(RTRIM(ISNULL(b.Denominacion, ''))) AS Denominacion,
            LTRIM(RTRIM(ISNULL(b.Concentracion, ''))) AS Concentracion,
            LTRIM(RTRIM(ISNULL(b.Presentacion, ''))) AS Presentacion,
            LTRIM(RTRIM(ISNULL(b.FormaFarmaceutica, ''))) AS FormaFarmaceutica,
            LTRIM(RTRIM(ISNULL(b.MaterialEnvase, ''))) AS MaterialEnvase,
            LTRIM(RTRIM(ISNULL(b.PresentacionEnvase, ''))) AS PresentacionEnvase,
            LTRIM(RTRIM(ISNULL(b.Fabricante, ''))) AS Fabricante,
            ISNULL(b.IdPaisOrigen, 0) AS IdPaisOrigen,
            LTRIM(RTRIM(ISNULL(pais.Nombre, ''))) AS PaisOrigen,
            CAST(ISNULL(b.Petitorio, 0) AS BIT) AS Petitorio,
            LTRIM(RTRIM(ISNULL(b.TipoProductoSismed, ''))) AS TipoProductoSismed,
            ISNULL(b.Estado, 0) AS Estado,
            ISNULL(b.estadoDigemid, 0) AS EstadoDigemid,
            ISNULL(b.idTipoPrefijo, 0) AS IdTipoPrefijo,
            LTRIM(RTRIM(ISNULL(pref.Descripcion, ''))) AS TipoPrefijo,
            ISNULL(b.PoliticaCanje, 0) AS PoliticaCanje,
            CAST(CASE WHEN ISNULL(b.PetitorioIst, 0) = 1 THEN 1 ELSE 0 END AS BIT) AS PetitorioIst,
            CAST(CASE WHEN ISNULL(b.Comunes, 0) = 1 THEN 1 ELSE 0 END AS BIT) AS Comunes,
            CAST(CASE WHEN ISNULL(b.ComiteFarmacologico, 0) = 1 THEN 1 ELSE 0 END AS BIT) AS ComiteFarmacologico,
            CAST(CASE WHEN ISNULL(b.Especializados, 0) = 1 THEN 1 ELSE 0 END AS BIT) AS Especializados,
            CAST(CASE WHEN ISNULL(b.Vitales, 0) = 1 THEN 1 ELSE 0 END AS BIT) AS EsVital,
            CAST(CASE WHEN ISNULL(b.Estrategicos, 0) = 1 THEN 1 ELSE 0 END AS BIT) AS EsEstrategico,
            CAST(CASE WHEN ISNULL(b.Oncologicos, 0) = 1 THEN 1 ELSE 0 END AS BIT) AS EsOncologico,
            CAST(CASE WHEN ISNULL(b.CargaEnFUA, 0) = 1 THEN 1 ELSE 0 END AS BIT) AS CargaEnFua,
            CAST(CASE WHEN ISNULL(b.ProductoIndicadorGestion, 0) = 1 THEN 1 ELSE 0 END AS BIT) AS ProductoIndicadorGestion
        FROM dbo.FactCatalogoBienesInsumos b WITH (NOLOCK)
        LEFT JOIN dbo.FactInsumosGrupoFarmacologico g WITH (NOLOCK) 
            ON g.IdGrupoFarmacologico = b.IdGrupoFarmacologico
        LEFT JOIN dbo.FactInsumosSubGrupoFarmacologico sg WITH (NOLOCK)
            ON sg.IdSubGrupoFarmacologico = b.IdSubGrupoFarmacologico
        LEFT JOIN dbo.FactPartidasPresupuestales p WITH (NOLOCK)
            ON p.IdPartidaPresupuestal = b.IdPartida
        LEFT JOIN dbo.CentrosCosto cc WITH (NOLOCK)
            ON cc.IdCentroCosto = b.IdCentroCosto
        LEFT JOIN dbo.Paises pais WITH (NOLOCK)
            ON pais.IdPais = b.IdPaisOrigen
        LEFT JOIN dbo.Tb_FarmtiposPrefijo pref WITH (NOLOCK)
            ON pref.idTipoSufijos = b.idTipoPrefijo
        WHERE b.IdProducto = @IdProducto;

    END TRY
    BEGIN CATCH
        DECLARE @ErrMsg NVARCHAR(4000) = ERROR_MESSAGE();
        DECLARE @ErrSeverity INT = ERROR_SEVERITY();
        DECLARE @ErrState INT = ERROR_STATE();
        RAISERROR(@ErrMsg, @ErrSeverity, @ErrState);
    END CATCH
END;
GO
