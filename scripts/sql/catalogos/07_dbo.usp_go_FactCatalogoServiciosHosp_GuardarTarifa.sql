-- ============================================================================
-- Procedimiento: dbo.usp_go_FactCatalogoServiciosHosp_GuardarTarifa
-- Módulo: Catálogos / Servicios Médicos y Tarifarios (HNSEB)
-- Objetivo: Registrar o actualizar la tarifa hospitalaria de un servicio
--           para un tipo de financiamiento específico (Particular, SIS, etc.).
-- Base de Datos: SIGH / Galenos Pro
-- Motor: SQL Server 2016 / 2019 / 2022
-- ============================================================================

SET ANSI_NULLS ON;
GO
SET QUOTED_IDENTIFIER ON;
GO

CREATE OR ALTER PROCEDURE dbo.usp_go_FactCatalogoServiciosHosp_GuardarTarifa
    @IdProducto INT,
    @IdTipoFinanciamiento INT,
    @PrecioUnitario DECIMAL(12,4),
    @Activo BIT = 1
AS
BEGIN
    SET NOCOUNT ON;

    BEGIN TRY
        IF @IdProducto IS NULL OR @IdProducto <= 0
        BEGIN
            RAISERROR('El identificador de producto es requerido.', 16, 1);
            RETURN;
        END;

        IF @IdTipoFinanciamiento IS NULL OR @IdTipoFinanciamiento <= 0
        BEGIN
            RAISERROR('El tipo de financiamiento es requerido.', 16, 1);
            RETURN;
        END;

        IF EXISTS (
            SELECT 1 
            FROM dbo.FactCatalogoServiciosHosp 
            WHERE IdProducto = @IdProducto AND IdTipoFinanciamiento = @IdTipoFinanciamiento
        )
        BEGIN
            UPDATE dbo.FactCatalogoServiciosHosp
            SET 
                PrecioUnitario = @PrecioUnitario,
                Activo = CASE WHEN @Activo = 1 THEN 1 ELSE 0 END
            WHERE IdProducto = @IdProducto AND IdTipoFinanciamiento = @IdTipoFinanciamiento;
        END
        ELSE
        BEGIN
            INSERT INTO dbo.FactCatalogoServiciosHosp (
                IdProducto, 
                IdTipoFinanciamiento, 
                PrecioUnitario, 
                Activo, 
                SeUsaSinPrecio
            )
            VALUES (
                @IdProducto, 
                @IdTipoFinanciamiento, 
                @PrecioUnitario, 
                CASE WHEN @Activo = 1 THEN 1 ELSE 0 END, 
                0
            );
        END;

    END TRY
    BEGIN CATCH
        DECLARE @ErrMsg NVARCHAR(4000) = ERROR_MESSAGE();
        DECLARE @ErrSeverity INT = ERROR_SEVERITY();
        DECLARE @ErrState INT = ERROR_STATE();
        RAISERROR(@ErrMsg, @ErrSeverity, @ErrState);
    END CATCH
END;
GO
