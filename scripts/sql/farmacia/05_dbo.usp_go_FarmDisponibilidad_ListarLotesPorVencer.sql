-- ============================================================================
-- Procedimiento: dbo.usp_go_FarmDisponibilidad_ListarLotesPorVencer
-- Módulo: Farmacia / Disponibilidad de Medicamentos y Salidas (HNSEB)
-- Objetivo: Reporte consolidado de lotes en riesgo de caducidad (30, 60 o 90 días)
--           para el modal de control de vencimientos y auditoría de farmacia.
-- Base de Datos: SIGH / Galenos Pro
-- Motor: SQL Server 2016 / 2019 / 2022
-- ============================================================================

SET ANSI_NULLS ON;
GO
SET QUOTED_IDENTIFIER ON;
GO

CREATE OR ALTER PROCEDURE dbo.usp_go_FarmDisponibilidad_ListarLotesPorVencer
    @IdAlmacen INT = 0,
    @DiasMax INT = 60
AS
BEGIN
    SET NOCOUNT ON;
    SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;

    BEGIN TRY
        DECLARE @AlmParam INT = @IdAlmacen;
        IF @AlmParam <= 0
            SET @AlmParam = 2;

        IF @DiasMax <= 0
            SET @DiasMax = 60;

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
            t.idProducto AS IdProducto,
            LTRIM(RTRIM(ISNULL(t.Codigo, ''))) AS Codigo,
            LTRIM(RTRIM(ISNULL(t.Nombre, ''))) AS Nombre,
            LTRIM(RTRIM(ISNULL(t.Lote, 'S/L'))) AS Lote,
            CONVERT(VARCHAR(10), t.FechaVencimiento, 120) AS FechaVencimiento,
            DATEDIFF(day, GETDATE(), t.FechaVencimiento) AS DiasParaVencer,
            CASE 
                WHEN DATEDIFF(day, GETDATE(), t.FechaVencimiento) <= 30 THEN 'CRITICO'
                ELSE 'ALERTA'
            END AS EstadoVencimiento,
            ISNULL(t.Saldo, 0) AS Cantidad,
            ISNULL(t.Precio, 0) AS Precio,
            LTRIM(RTRIM(ISNULL(t.Nombre, ''))) AS NombreMarca
        FROM #TmpLotes t
        WHERE t.FechaVencimiento IS NOT NULL
          AND ISNULL(t.Saldo, 0) > 0
          AND DATEDIFF(day, GETDATE(), t.FechaVencimiento) BETWEEN 0 AND @DiasMax
        ORDER BY DATEDIFF(day, GETDATE(), t.FechaVencimiento) ASC, t.Saldo DESC;

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
