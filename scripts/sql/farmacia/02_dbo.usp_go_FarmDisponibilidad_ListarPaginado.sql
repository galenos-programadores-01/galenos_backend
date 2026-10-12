SET ANSI_NULLS ON;
GO
SET QUOTED_IDENTIFIER ON;
GO

CREATE OR ALTER PROCEDURE dbo.usp_go_FarmDisponibilidad_ListarPaginado
    @IdAlmacen INT = 0,
    @Filtro VARCHAR(150) = NULL,
    @EstadoSemaforo VARCHAR(20) = NULL,
    @SoloConStock BIT = 0,
    @SoloEsenciales BIT = 0,
    @CanastaTipo VARCHAR(20) = 'TODOS',
    @FechaInicio VARCHAR(30) = NULL,
    @FechaFin VARCHAR(30) = NULL,
    @SoloConMovimiento BIT = 0,
    @PageNumber INT = 1,
    @PageSize INT = 20
AS
BEGIN
    SET NOCOUNT ON;
    SET DATEFORMAT ymd;
    SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;

    BEGIN TRY
        DECLARE @AlmParam INT = @IdAlmacen;
        IF @AlmParam <= 0
            SET @AlmParam = 2;

        IF @PageNumber < 1 SET @PageNumber = 1;
        IF @PageSize < 1 SET @PageSize = 20;

        DECLARE @CanastaParam VARCHAR(20) = UPPER(ISNULL(LTRIM(RTRIM(@CanastaTipo)), 'TODOS'));

        DECLARE @AlmacenNombre VARCHAR(150) = 'FARMACIA PRINCIPAL';
        SELECT TOP 1 @AlmacenNombre = LTRIM(RTRIM(Descripcion))
        FROM dbo.farmAlmacen WITH (NOLOCK)
        WHERE idAlmacen = @AlmParam;

        DECLARE @FiltroSP VARCHAR(50) = ISNULL(LTRIM(RTRIM(@Filtro)), '');

        DECLARE @SemaforoNormalizado VARCHAR(20) = NULL;
        IF @EstadoSemaforo IS NOT NULL AND LTRIM(RTRIM(@EstadoSemaforo)) <> ''
        BEGIN
            DECLARE @SemUpper VARCHAR(20) = UPPER(LTRIM(RTRIM(@EstadoSemaforo)));
            IF @SemUpper IN ('NORMAL', 'DISPONIBLE') SET @SemaforoNormalizado = 'NORMAL';
            ELSE IF @SemUpper IN ('RIESGO', 'EN_RIESGO', 'SUBSTOCK') SET @SemaforoNormalizado = 'RIESGO';
            ELSE IF @SemUpper IN ('QUIEBRE', 'EN_QUIEBRE', 'AGOTADO', 'AGOTADOS') SET @SemaforoNormalizado = 'QUIEBRE';
        END;

        DECLARE @FInicioParsed DATETIME = NULL;
        DECLARE @FFinParsed DATETIME = NULL;

        IF @FechaInicio IS NOT NULL AND LTRIM(RTRIM(@FechaInicio)) != ''
        BEGIN
            IF ISDATE(@FechaInicio) = 1
                SET @FInicioParsed = CONVERT(DATETIME, LTRIM(RTRIM(@FechaInicio)), 120);
        END;

        IF @FechaFin IS NOT NULL AND LTRIM(RTRIM(@FechaFin)) != ''
        BEGIN
            IF ISDATE(@FechaFin) = 1
            BEGIN
                SET @FFinParsed = CONVERT(DATETIME, LTRIM(RTRIM(@FechaFin)), 120);
                IF DATEPART(HOUR, @FFinParsed) = 0 AND DATEPART(MINUTE, @FFinParsed) = 0
                    SET @FFinParsed = DATEADD(SECOND, 86399, @FFinParsed);
            END;
        END;

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
            @Filtro = @FiltroSP;

        CREATE TABLE #TmpMovimientosPeriodo (
            IdProducto INT PRIMARY KEY,
            CantidadDispensada DECIMAL(12,2)
        );

        IF @FInicioParsed IS NOT NULL AND @FFinParsed IS NOT NULL
        BEGIN
            INSERT INTO #TmpMovimientosPeriodo (IdProducto, CantidadDispensada)
            SELECT 
                d.idProducto,
                SUM(CAST(d.Cantidad AS DECIMAL(12,2)))
            FROM dbo.farmMovimiento m WITH (NOLOCK)
            INNER JOIN dbo.farmMovimientoDetalle d WITH (NOLOCK)
                ON m.MovNumero = d.MovNumero AND m.MovTipo = d.MovTipo
            WHERE m.MovTipo = 'S'
              AND m.fechaCreacion >= @FInicioParsed
              AND m.fechaCreacion <= @FFinParsed
              AND (@IdAlmacen = 0 OR m.idAlmacenOrigen = @AlmParam)
              AND ISNULL(m.idEstadoMovimiento, 1) != 3
            GROUP BY d.idProducto;
        END;

        CREATE TABLE #TmpFiltrados (
            IdProducto INT,
            Codigo VARCHAR(50),
            Nombre VARCHAR(250),
            Concentracion VARCHAR(150),
            FormaFarmaceutica VARCHAR(150),
            IdAlmacen INT,
            AlmacenNombre VARCHAR(150),
            Cantidad DECIMAL(12,2),
            Precio DECIMAL(12,4),
            StockMinimo INT,
            EstadoSemaforo VARCHAR(20),
            TotalLotes INT,
            Petitorio BIT,
            SalidasPeriodo DECIMAL(12,2)
        );

        INSERT INTO #TmpFiltrados
        SELECT 
            t.idProducto AS IdProducto,
            LTRIM(RTRIM(ISNULL(t.Codigo, ''))) AS Codigo,
            LTRIM(RTRIM(ISNULL(t.Nombre, ''))) AS Nombre,
            LTRIM(RTRIM(ISNULL(b.Concentracion, ''))) AS Concentracion,
            LTRIM(RTRIM(ISNULL(b.FormaFarmaceutica, ISNULL(t.tipo, '')))) AS FormaFarmaceutica,
            t.idAlmacen AS IdAlmacen,
            @AlmacenNombre AS AlmacenNombre,
            ISNULL(t.saldo, 0) AS Cantidad,
            ISNULL(t.Precio, 0) AS Precio,
            CAST(ISNULL(b.StockMinimo, 10) AS INT) AS StockMinimo,
            CASE 
                WHEN ISNULL(t.saldo, 0) > CAST(ISNULL(b.StockMinimo, 10) AS INT) THEN 'NORMAL'
                WHEN ISNULL(t.saldo, 0) > 0 THEN 'RIESGO'
                ELSE 'QUIEBRE'
            END AS EstadoSemaforo,
            1 AS TotalLotes,
            CAST(ISNULL(b.Petitorio, 1) AS BIT) AS Petitorio,
            ISNULL(mp.CantidadDispensada, 0) AS SalidasPeriodo
        FROM #TmpSaldos t
        LEFT JOIN dbo.FactCatalogoBienesInsumos b WITH (NOLOCK) 
            ON b.IdProducto = t.idProducto
        LEFT JOIN dbo.Far_ItemsEvaluadosSisDigemid e WITH (NOLOCK)
            ON e.Codigo = b.Codigo
        LEFT JOIN #TmpMovimientosPeriodo mp
            ON mp.IdProducto = t.idProducto
        WHERE (@SoloConStock = 0 OR ISNULL(t.saldo, 0) > 0)
          AND (
            @SoloEsenciales = 0 
            OR (@CanastaParam = 'TODOS' AND (e.Codigo IS NOT NULL OR ISNULL(b.ProductoIndicadorGestion, 0) = 1))
            OR (@CanastaParam = 'DIGEMID' AND e.EsDigemid = 1)
            OR (@CanastaParam = 'SIS' AND e.EsSis = 1)
          )
          AND (
            @CanastaParam = 'TODOS'
            OR @SoloEsenciales = 1
            OR (@CanastaParam = 'DIGEMID' AND e.EsDigemid = 1)
            OR (@CanastaParam = 'SIS' AND e.EsSis = 1)
          )
          AND (@SemaforoNormalizado IS NULL OR 
               CASE 
                   WHEN ISNULL(t.saldo, 0) > CAST(ISNULL(b.StockMinimo, 10) AS INT) THEN 'NORMAL'
                   WHEN ISNULL(t.saldo, 0) > 0 THEN 'RIESGO'
                   ELSE 'QUIEBRE'
               END = @SemaforoNormalizado)
          AND (
              @SoloConMovimiento = 0 
              OR @FInicioParsed IS NULL 
              OR ISNULL(mp.CantidadDispensada, 0) > 0
          );

        DECLARE @TotalRegistros INT = 0;
        SELECT @TotalRegistros = COUNT(1) FROM #TmpFiltrados;

        SELECT 
            f.IdProducto,
            f.Codigo,
            f.Nombre,
            f.Concentracion,
            f.FormaFarmaceutica,
            f.IdAlmacen,
            f.AlmacenNombre,
            f.Cantidad,
            f.Precio,
            CAST(f.StockMinimo AS INT) AS StockMinimo,
            f.EstadoSemaforo,
            f.TotalLotes,
            f.Petitorio,
            f.SalidasPeriodo,
            @TotalRegistros AS TotalRegistros
        FROM #TmpFiltrados f
        ORDER BY 
            CASE WHEN @FInicioParsed IS NOT NULL AND @FFinParsed IS NOT NULL THEN f.SalidasPeriodo ELSE 0 END DESC,
            f.Nombre ASC, 
            f.IdProducto ASC
        OFFSET (@PageNumber - 1) * @PageSize ROWS
        FETCH NEXT @PageSize ROWS ONLY;

        DROP TABLE #TmpSaldos;
        DROP TABLE #TmpFiltrados;
        DROP TABLE #TmpMovimientosPeriodo;

    END TRY
    BEGIN CATCH
        IF OBJECT_ID('tempdb..#TmpSaldos') IS NOT NULL DROP TABLE #TmpSaldos;
        IF OBJECT_ID('tempdb..#TmpFiltrados') IS NOT NULL DROP TABLE #TmpFiltrados;
        IF OBJECT_ID('tempdb..#TmpMovimientosPeriodo') IS NOT NULL DROP TABLE #TmpMovimientosPeriodo;

        DECLARE @ErrMsg NVARCHAR(4000) = ERROR_MESSAGE();
        DECLARE @ErrSeverity INT = ERROR_SEVERITY();
        DECLARE @ErrState INT = ERROR_STATE();
        RAISERROR(@ErrMsg, @ErrSeverity, @ErrState);
    END CATCH
END;
GO
