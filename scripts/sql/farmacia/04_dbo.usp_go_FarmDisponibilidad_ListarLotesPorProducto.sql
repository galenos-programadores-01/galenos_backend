-- ============================================================================
-- Procedimiento: dbo.usp_go_FarmDisponibilidad_ListarLotesPorProducto
-- Módulo: Farmacia / Disponibilidad de Medicamentos y Salidas (HNSEB)
-- Objetivo: Obtener exclusivamente los lotes vigentes y con saldo de un producto
--           específico para el modal de lotes del catálogo de existencias.
-- Base de Datos: SIGH / Galenos Pro
-- Motor: SQL Server 2016 / 2019 / 2022
-- ============================================================================

SET ANSI_NULLS ON;
GO
SET QUOTED_IDENTIFIER ON;
GO

CREATE OR ALTER PROCEDURE dbo.usp_go_FarmDisponibilidad_ListarLotesPorProducto
    @IdProducto INT,
    @IdAlmacen INT = 0
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

        DECLARE @AlmParam INT = @IdAlmacen;
        IF @AlmParam <= 0
            SET @AlmParam = 2;

        CREATE TABLE #TmpLotes (
            Codigo VARCHAR(50), Nombre VARCHAR(250), tipo VARCHAR(100), Lote VARCHAR(50),
            FechaVencimiento DATETIME, Saldo DECIMAL(12,2), idAlmacen INT, idProducto INT,
            PrecioUnitario DECIMAL(12,4), Precio DECIMAL(12,4), idTipoSalidaSaldo INT,
            polCanje INT, disp VARCHAR(50), idTipoPresentacion INT, idTipoSalida INT, num INT
        );

        INSERT INTO #TmpLotes
        EXEC dbo.farmDevuelveSaldosConLotesSegunAlmacen 
            @IdAlmacen = @AlmParam, 
            @Orden = 1, 
            @Filtro = '';

        SELECT 
            LTRIM(RTRIM(ISNULL(t.Lote, 'S/L'))) AS Lote,
            CASE 
                WHEN t.FechaVencimiento IS NOT NULL THEN CONVERT(VARCHAR(10), t.FechaVencimiento, 120)
                ELSE ''
            END AS FechaVencimiento,
            ISNULL(t.Saldo, 0) AS Cantidad,
            ISNULL(t.Precio, 0) AS Precio,
            LTRIM(RTRIM(ISNULL(t.Nombre, ''))) AS NombreMarca,
            CASE 
                WHEN t.FechaVencimiento IS NOT NULL THEN DATEDIFF(day, GETDATE(), t.FechaVencimiento)
                ELSE 0
            END AS DiasParaVencer,
            CASE 
                WHEN t.FechaVencimiento IS NULL THEN 'SIN_FECHA'
                WHEN DATEDIFF(day, GETDATE(), t.FechaVencimiento) < 0 THEN 'VENCIDO'
                WHEN DATEDIFF(day, GETDATE(), t.FechaVencimiento) <= 30 THEN 'CRITICO_30'
                WHEN DATEDIFF(day, GETDATE(), t.FechaVencimiento) <= 60 THEN 'ALERTA_60'
                ELSE 'VIGENTE'
            END AS EstadoVencimiento
        FROM #TmpLotes t
        WHERE t.idProducto = @IdProducto
        ORDER BY t.FechaVencimiento ASC, t.Saldo DESC;

        DROP TABLE #TmpLotes;

    END TRY
    BEGIN CATCH
        IF OBJECT_ID('tempdb..#TmpLotes') IS NOT NULL DROP TABLE #TmpLotes;

        DECLARE @ErrMsg NVARCHAR(4000) = ERROR_MESSAGE();
        DECLARE @ErrSeverity INT = ERROR_SEVERITY();
        DECLARE @ErrState INT = ERROR_STATE();
        RAISERROR(@ErrMsg, @ErrSeverity, @ErrState);
    END CATCH
END;
GO
