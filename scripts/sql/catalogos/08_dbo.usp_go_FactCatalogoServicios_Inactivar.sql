-- ============================================================================
-- Procedimiento: dbo.usp_go_FactCatalogoServicios_Inactivar
-- Módulo: Catálogos / Servicios Médicos y Tarifarios (HNSEB)
-- Objetivo: Inactivar lógicamente un servicio del catálogo hospitalario (IdEstado = 0).
-- Base de Datos: SIGH / Galenos Pro
-- Motor: SQL Server 2016 / 2019 / 2022
-- ============================================================================

SET ANSI_NULLS ON;
GO
SET QUOTED_IDENTIFIER ON;
GO

CREATE OR ALTER PROCEDURE dbo.usp_go_FactCatalogoServicios_Inactivar
    @IdProducto INT,
    @IdUsuarioAuditoria INT = 1
AS
BEGIN
    SET NOCOUNT ON;

    BEGIN TRY
        IF @IdProducto IS NULL OR @IdProducto <= 0
        BEGIN
            RAISERROR('El identificador de producto es requerido.', 16, 1);
            RETURN;
        END;

        UPDATE dbo.FactCatalogoServicios
        SET 
            idEstado = 0
        WHERE IdProducto = @IdProducto;

        IF @@ROWCOUNT = 0
        BEGIN
            RAISERROR('No se encontró el servicio hospitalario con el ID especificado.', 16, 1);
            RETURN;
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
