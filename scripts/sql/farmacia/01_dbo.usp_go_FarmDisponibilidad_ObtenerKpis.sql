SET ANSI_NULLS ON;
GO
SET QUOTED_IDENTIFIER ON;
GO

CREATE OR ALTER PROCEDURE dbo.usp_go_FarmDisponibilidad_ObtenerKpis
    @IdAlmacen INT = 0,
    @CanastaTipo VARCHAR(20) = 'TODOS'
AS
BEGIN
    SET NOCOUNT ON;
    SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;

    BEGIN TRY
        DECLARE @AlmParam INT = @IdAlmacen;
        IF @AlmParam <= 0
            SET @AlmParam = 2;

        DECLARE @CanastaParam VARCHAR(20) = UPPER(ISNULL(LTRIM(RTRIM(@CanastaTipo)), 'TODOS'));

        CREATE TABLE #TmpSaldos (
            idAlmacen INT,
            idProducto INT,
            Codigo VARCHAR(50),
            Nombre VARCHAR(250),
            tipo VARCHAR(100),
            saldo DECIMAL(12,2),
            Precio DECIMAL(12,4),
            idTipoSalidaBienInsumo INT,
            PrecioUnitario DECIMAL(12,4),
            idTipoFinanciamiento INT,
            idTipoSalidaSaldo INT,
            disp VARCHAR(50),
            idTipoPresentacion INT
        );

        INSERT INTO #TmpSaldos
        EXEC dbo.farmDevuelveSaldosSinLotesSegunAlmacen 
            @IdAlmacen = @AlmParam, 
            @Orden = 1, 
            @Filtro = '';

        DECLARE @TotalProductos INT = 0;
        DECLARE @Disponibles INT = 0;
        DECLARE @EnRiesgo INT = 0;
        DECLARE @EnQuiebre INT = 0;

        SELECT 
            @TotalProductos = COUNT(1),
            @Disponibles = SUM(CASE WHEN ISNULL(t.saldo, 0) > 10 THEN 1 ELSE 0 END),
            @EnRiesgo = SUM(CASE WHEN ISNULL(t.saldo, 0) > 0 AND ISNULL(t.saldo, 0) <= 10 THEN 1 ELSE 0 END),
            @EnQuiebre = SUM(CASE WHEN ISNULL(t.saldo, 0) <= 0 THEN 1 ELSE 0 END)
        FROM #TmpSaldos t;

        DECLARE @TotalEsenciales INT = 0;
        DECLARE @EsencialesDisponibles INT = 0;
        DECLARE @Normostock INT = 0;
        DECLARE @Substock INT = 0;
        DECLARE @Sobrestock INT = 0;

        SELECT 
            @TotalEsenciales = COUNT(DISTINCT e.Codigo),
            @EsencialesDisponibles = SUM(CASE WHEN ISNULL(s.saldo, 0) > 0 THEN 1 ELSE 0 END),
            @Substock = SUM(CASE WHEN ISNULL(s.saldo, 0) > 0 AND (ISNULL(s.saldo, 0) < ISNULL(e.ConsumoPromedio, 10)) THEN 1 ELSE 0 END),
            @Normostock = SUM(CASE WHEN ISNULL(s.saldo, 0) >= ISNULL(e.ConsumoPromedio, 10) AND (ISNULL(s.saldo, 0) <= (ISNULL(e.ConsumoPromedio, 10) * 6)) THEN 1 ELSE 0 END),
            @Sobrestock = SUM(CASE WHEN ISNULL(s.saldo, 0) > (ISNULL(e.ConsumoPromedio, 10) * 6) THEN 1 ELSE 0 END)
        FROM dbo.Far_ItemsEvaluadosSisDigemid e
        LEFT JOIN dbo.FactCatalogoBienesInsumos b WITH (NOLOCK) ON b.Codigo = e.Codigo
        LEFT JOIN #TmpSaldos s ON s.idProducto = b.IdProducto
        WHERE (
            @CanastaParam = 'TODOS'
            OR (@CanastaParam = 'DIGEMID' AND e.EsDigemid = 1)
            OR (@CanastaParam = 'SIS' AND e.EsSis = 1)
        );

        DECLARE @PctDisponibilidadEsenciales DECIMAL(5,2) = 0.0;
        IF @TotalEsenciales > 0
            SET @PctDisponibilidadEsenciales = ROUND((CAST(@EsencialesDisponibles AS DECIMAL(10,2)) / CAST(@TotalEsenciales AS DECIMAL(10,2))) * 100.0, 2);

        DECLARE @PorcentajeOptimo DECIMAL(5,2) = 0.0;
        IF @TotalProductos > 0
            SET @PorcentajeOptimo = ROUND((CAST(@Disponibles AS DECIMAL(10,2)) / CAST(@TotalProductos AS DECIMAL(10,2))) * 100.0, 2);

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

        DECLARE @LotesPorVencer30 INT = 0;
        DECLARE @LotesPorVencer60 INT = 0;

        SELECT 
            @LotesPorVencer30 = SUM(CASE WHEN DATEDIFF(day, GETDATE(), FechaVencimiento) <= 30 THEN 1 ELSE 0 END),
            @LotesPorVencer60 = SUM(CASE WHEN DATEDIFF(day, GETDATE(), FechaVencimiento) > 30 AND DATEDIFF(day, GETDATE(), FechaVencimiento) <= 60 THEN 1 ELSE 0 END)
        FROM #TmpLotes
        WHERE FechaVencimiento IS NOT NULL AND ISNULL(Saldo, 0) > 0;

        SELECT 
            ISNULL(@TotalProductos, 0) AS TotalProductos,
            ISNULL(@Disponibles, 0) AS Disponibles,
            ISNULL(@EnRiesgo, 0) AS EnRiesgo,
            ISNULL(@EnQuiebre, 0) AS EnQuiebre,
            ISNULL(@Normostock, 0) AS Normostock,
            ISNULL(@Substock, 0) AS Substock,
            ISNULL(@Sobrestock, 0) AS Sobrestock,
            ISNULL(@TotalEsenciales, 0) AS TotalEsenciales,
            ISNULL(@EsencialesDisponibles, 0) AS EsencialesDisponibles,
            ISNULL(@PctDisponibilidadEsenciales, 0) AS PctDisponibilidadEsenciales,
            ISNULL(@PorcentajeOptimo, 0) AS PorcentajeOptimo,
            ISNULL(@LotesPorVencer30, 0) AS LotesPorVencer30,
            ISNULL(@LotesPorVencer60, 0) AS LotesPorVencer60;

        DROP TABLE #TmpSaldos;
        DROP TABLE #TmpLotes;

    END TRY
    BEGIN CATCH
        IF OBJECT_ID('tempdb..#TmpSaldos') IS NOT NULL DROP TABLE #TmpSaldos;
        IF OBJECT_ID('tempdb..#TmpLotes') IS NOT NULL DROP TABLE #TmpLotes;

        DECLARE @ErrMsg NVARCHAR(4000) = ERROR_MESSAGE();
        DECLARE @ErrSeverity INT = ERROR_SEVERITY();
        DECLARE @ErrState INT = ERROR_STATE();
        RAISERROR(@ErrMsg, @ErrSeverity, @ErrState);
    END CATCH
END;
GO
