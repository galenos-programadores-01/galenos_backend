-- ============================================================================
-- Procedimiento: dbo.usp_go_FactCatalogoServicios_ActualizarNormativa
-- Módulo: Catálogos / Servicios Médicos y Procedimientos CPT (HNSEB)
-- Objetivo: Actualizar los campos normativos, códigos de salud y banderas
--           asistenciales (Nombre Corto, CodMINSA, Codigo_HIS, CPMS, CargaEnFUA,
--           CargaEnHIS) que el procedimiento legado institucional no contempla.
-- Base de Datos: SIGH / Galenos Pro
-- Motor: SQL Server 2016 / 2019 / 2022
-- ============================================================================

SET ANSI_NULLS ON;
GO
SET QUOTED_IDENTIFIER ON;
GO

CREATE OR ALTER PROCEDURE dbo.usp_go_FactCatalogoServicios_ActualizarNormativa
    @IdProducto INT,
    @NombreCorto VARCHAR(250) = '',
    @CodMinsa VARCHAR(20) = '',
    @CodigoHis VARCHAR(20) = '',
    @Cpms VARCHAR(10) = '',
    @CargaEnFua BIT = 0,
    @CargaEnHis BIT = 0
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
            Nombre_Corto = LTRIM(RTRIM(@NombreCorto)),
            CodMINSA = LTRIM(RTRIM(@CodMinsa)),
            Codigo_HIS = LTRIM(RTRIM(@CodigoHis)),
            CPMS = LTRIM(RTRIM(@Cpms)),
            CargaEnFUA = CASE WHEN @CargaEnFua = 1 THEN 1 ELSE 0 END,
            CargaEnHIS = CASE WHEN @CargaEnHis = 1 THEN 1 ELSE 0 END
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
