-- ============================================================================
-- Procedimiento: dbo.usp_go_FactCatalogoBienesInsumos_ActualizarClasificaciones
-- Módulo: Catálogos / Bienes e Insumos Médicos (HNSEB)
-- Objetivo: Actualizar las banderas asistenciales y clasificadores especiales
--           (Petitorio, Comunes, Comité, Especializados, Vitales, Estratégicos,
--           Oncológicos, Carga en FUA, Prefijo y Política Canje) de un bien o insumo.
-- Base de Datos: SIGH / Galenos Pro
-- Motor: SQL Server 2016 / 2019 / 2022
-- ============================================================================

SET ANSI_NULLS ON;
GO
SET QUOTED_IDENTIFIER ON;
GO

CREATE OR ALTER PROCEDURE dbo.usp_go_FactCatalogoBienesInsumos_ActualizarClasificaciones
    @IdProducto INT,
    @Vitales BIT = 0,
    @Estrategicos BIT = 0,
    @Oncologicos BIT = 0,
    @CargaEnFUA BIT = 0,
    @PetitorioIst BIT = 0,
    @Comunes BIT = 0,
    @ComiteFarmacologico BIT = 0,
    @Especializados BIT = 0,
    @IdTipoPrefijo INT = 0,
    @PoliticaCanje INT = 0,
    @ProductoIndicadorGestion BIT = 0
AS
BEGIN
    SET NOCOUNT ON;

    BEGIN TRY
        IF @IdProducto IS NULL OR @IdProducto <= 0
        BEGIN
            RAISERROR('El identificador de producto es requerido.', 16, 1);
            RETURN;
        END;

        UPDATE dbo.FactCatalogoBienesInsumos
        SET 
            Vitales = CASE WHEN @Vitales = 1 THEN 1 ELSE 0 END,
            Estrategicos = CASE WHEN @Estrategicos = 1 THEN 1 ELSE 0 END,
            Oncologicos = CASE WHEN @Oncologicos = 1 THEN 1 ELSE 0 END,
            CargaEnFUA = CASE WHEN @CargaEnFUA = 1 THEN 1 ELSE 0 END,
            PetitorioIst = CASE WHEN @PetitorioIst = 1 THEN 1 ELSE 0 END,
            Comunes = CASE WHEN @Comunes = 1 THEN 1 ELSE 0 END,
            ComiteFarmacologico = CASE WHEN @ComiteFarmacologico = 1 THEN 1 ELSE 0 END,
            Especializados = CASE WHEN @Especializados = 1 THEN 1 ELSE 0 END,
            idTipoPrefijo = @IdTipoPrefijo,
            PoliticaCanje = @PoliticaCanje,
            ProductoIndicadorGestion = CASE WHEN @ProductoIndicadorGestion = 1 THEN 1 ELSE 0 END
        WHERE IdProducto = @IdProducto;

    END TRY
    BEGIN CATCH
        DECLARE @ErrMsg NVARCHAR(4000) = ERROR_MESSAGE();
        DECLARE @ErrSeverity INT = ERROR_SEVERITY();
        DECLARE @ErrState INT = ERROR_STATE();
        RAISERROR(@ErrMsg, @ErrSeverity, @ErrState);
    END CATCH
END;
GO
