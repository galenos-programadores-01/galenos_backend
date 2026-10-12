SET ANSI_NULLS ON;
GO
SET QUOTED_IDENTIFIER ON;
GO

CREATE OR ALTER PROCEDURE dbo.usp_go_FarmDisponibilidad_ObtenerEstadisticasDashboard
    @IdAlmacen INT = 0,
    @CanastaTipo VARCHAR(20) = 'TODOS',
    @FechaInicio VARCHAR(30) = NULL,
    @FechaFin VARCHAR(30) = NULL
AS
BEGIN
    SET NOCOUNT ON;
    SET DATEFORMAT ymd;
    SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;

    BEGIN TRY
        DECLARE @AlmParam INT = @IdAlmacen;
        IF @AlmParam <= 0
            SET @AlmParam = 2;

        DECLARE @CanastaParam VARCHAR(20) = UPPER(ISNULL(LTRIM(RTRIM(@CanastaTipo)), 'TODOS'));

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
            @Filtro = '';

        IF @FInicioParsed IS NOT NULL AND @FFinParsed IS NOT NULL
        BEGIN
            SELECT TOP 10
                ISNULL(b.IdProducto, d.idProducto) AS IdProducto,
                LTRIM(RTRIM(ISNULL(b.Codigo, ''))) AS Codigo,
                LTRIM(RTRIM(ISNULL(b.Nombre, ''))) AS Nombre,
                ISNULL(SUM(CAST(d.Cantidad AS BIGINT)), 0) AS CantidadTotal
            FROM dbo.farmMovimiento m WITH (NOLOCK)
            INNER JOIN dbo.farmMovimientoDetalle d WITH (NOLOCK) 
                ON m.MovNumero = d.MovNumero AND m.MovTipo = d.MovTipo
            LEFT JOIN dbo.FactCatalogoBienesInsumos b WITH (NOLOCK) 
                ON b.IdProducto = d.idProducto
            WHERE m.MovTipo = 'S'
              AND m.fechaCreacion >= @FInicioParsed 
              AND m.fechaCreacion <= @FFinParsed
              AND (@IdAlmacen = 0 OR m.idAlmacenOrigen = @AlmParam)
              AND ISNULL(m.idEstadoMovimiento, 1) != 3
            GROUP BY 
                ISNULL(b.IdProducto, d.idProducto),
                b.Codigo,
                b.Nombre
            ORDER BY CantidadTotal DESC;
        END
        ELSE
        BEGIN
            SELECT TOP 10
                t.idProducto AS IdProducto,
                LTRIM(RTRIM(ISNULL(t.Codigo, ''))) AS Codigo,
                LTRIM(RTRIM(ISNULL(t.Nombre, ''))) AS Nombre,
                ISNULL(t.saldo, 0) AS CantidadTotal
            FROM #TmpSaldos t
            WHERE ISNULL(t.saldo, 0) > 0
            ORDER BY t.saldo DESC;
        END;

        CREATE TABLE #TmpAlmacenes (
            IdAlmacen INT,
            Descripcion VARCHAR(150)
        );

        INSERT INTO #TmpAlmacenes (IdAlmacen, Descripcion)
        VALUES 
            (2, 'Almacen Especializado'),
            (4, 'Farmacia Principal'),
            (8, 'Farmacia Emergencia'),
            (47, 'Dosis Unitaria'),
            (6, 'Farmacia SOP');

        IF @FInicioParsed IS NOT NULL AND @FFinParsed IS NOT NULL
        BEGIN
            SELECT 
                a.IdAlmacen,
                a.Descripcion,
                ISNULL(SUM(CASE WHEN ISNULL(d.IdTipoSalidaBienInsumo, 0) = 1 THEN d.Cantidad ELSE 0 END), 0) AS Normales,
                ISNULL(SUM(CASE WHEN ISNULL(d.IdTipoSalidaBienInsumo, 0) = 2 THEN d.Cantidad ELSE 0 END), 0) AS EnRiesgo,
                ISNULL(SUM(CASE WHEN ISNULL(d.IdTipoSalidaBienInsumo, 0) NOT IN (1, 2) THEN d.Cantidad ELSE 0 END), 0) AS EnQuiebre,
                ISNULL(SUM(d.Cantidad), 0) AS TotalItems
            FROM #TmpAlmacenes a
            LEFT JOIN dbo.farmMovimiento m WITH (NOLOCK) 
                ON m.idAlmacenOrigen = a.IdAlmacen 
               AND m.MovTipo = 'S' 
               AND m.fechaCreacion >= @FInicioParsed 
               AND m.fechaCreacion <= @FFinParsed
               AND ISNULL(m.idEstadoMovimiento, 1) != 3
            LEFT JOIN dbo.farmMovimientoDetalle d WITH (NOLOCK) 
                ON d.MovNumero = m.MovNumero AND d.MovTipo = m.MovTipo
            GROUP BY a.IdAlmacen, a.Descripcion;
        END
        ELSE
        BEGIN
            CREATE TABLE #TmpComparativa (
                IdAlmacen INT,
                Descripcion VARCHAR(150),
                Normales INT,
                EnRiesgo INT,
                EnQuiebre INT,
                TotalItems INT
            );

            DECLARE @CurIdAlm INT, @CurDesc VARCHAR(150);
            DECLARE curAlm CURSOR LOCAL FAST_FORWARD FOR 
                SELECT IdAlmacen, Descripcion FROM #TmpAlmacenes;

            OPEN curAlm;
            FETCH NEXT FROM curAlm INTO @CurIdAlm, @CurDesc;

            WHILE @@FETCH_STATUS = 0
            BEGIN
                CREATE TABLE #TmpLoopSaldos (
                    idAlmacen INT, idProducto INT, Codigo VARCHAR(50), Nombre VARCHAR(250),
                    tipo VARCHAR(100), saldo DECIMAL(12,2), Precio DECIMAL(12,4),
                    idTipoSalidaBienInsumo INT, PrecioUnitario DECIMAL(12,4),
                    idTipoFinanciamiento INT, idTipoSalidaSaldo INT, disp VARCHAR(50), idTipoPresentacion INT
                );

                INSERT INTO #TmpLoopSaldos
                EXEC dbo.farmDevuelveSaldosSinLotesSegunAlmacen 
                    @IdAlmacen = @CurIdAlm, 
                    @Orden = 1, 
                    @Filtro = '';

                INSERT INTO #TmpComparativa (IdAlmacen, Descripcion, Normales, EnRiesgo, EnQuiebre, TotalItems)
                SELECT 
                    @CurIdAlm,
                    @CurDesc,
                    SUM(CASE WHEN ISNULL(saldo, 0) > 10 THEN 1 ELSE 0 END),
                    SUM(CASE WHEN ISNULL(saldo, 0) > 0 AND ISNULL(saldo, 0) <= 10 THEN 1 ELSE 0 END),
                    SUM(CASE WHEN ISNULL(saldo, 0) <= 0 THEN 1 ELSE 0 END),
                    COUNT(1)
                FROM #TmpLoopSaldos;

                DROP TABLE #TmpLoopSaldos;
                FETCH NEXT FROM curAlm INTO @CurIdAlm, @CurDesc;
            END;

            CLOSE curAlm;
            DEALLOCATE curAlm;

            SELECT 
                IdAlmacen,
                Descripcion,
                ISNULL(Normales, 0) AS Normales,
                ISNULL(EnRiesgo, 0) AS EnRiesgo,
                ISNULL(EnQuiebre, 0) AS EnQuiebre,
                ISNULL(TotalItems, 0) AS TotalItems
            FROM #TmpComparativa;

            DROP TABLE #TmpComparativa;
        END;

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

        CREATE TABLE #TmpTramosVenc (
            Orden INT,
            Rango VARCHAR(50),
            CantidadLotes INT,
            UnidadesTotales DECIMAL(12,2),
            NivelAlerta VARCHAR(20)
        );

        INSERT INTO #TmpTramosVenc (Orden, Rango, CantidadLotes, UnidadesTotales, NivelAlerta)
        VALUES 
            (1, 'Vencen este mes (Urgente)', 0, 0, 'CRITICO'),
            (2, 'Vencen el proximo mes',     0, 0, 'ALERTA'),
            (3, 'De 2 a 3 meses',            0, 0, 'PRECAUCION'),
            (4, 'De 3 a 6 meses',            0, 0, 'SEGUIMIENTO'),
            (5, 'Mas de 6 meses (Vigentes)', 0, 0, 'OPTIMO');

        UPDATE t
        SET 
            t.CantidadLotes = v.CantLotes,
            t.UnidadesTotales = v.TotalUnidades
        FROM #TmpTramosVenc t
        INNER JOIN (
            SELECT 
                CASE 
                    WHEN DATEDIFF(day, GETDATE(), FechaVencimiento) <= 30 THEN 'Vencen este mes (Urgente)'
                    WHEN DATEDIFF(day, GETDATE(), FechaVencimiento) <= 60 THEN 'Vencen el proximo mes'
                    WHEN DATEDIFF(day, GETDATE(), FechaVencimiento) <= 90 THEN 'De 2 a 3 meses'
                    WHEN DATEDIFF(day, GETDATE(), FechaVencimiento) <= 180 THEN 'De 3 a 6 meses'
                    ELSE 'Mas de 6 meses (Vigentes)'
                END AS RangoClave,
                COUNT(1) AS CantLotes,
                SUM(ISNULL(Saldo, 0)) AS TotalUnidades
            FROM #TmpLotes
            WHERE FechaVencimiento IS NOT NULL 
              AND ISNULL(Saldo, 0) > 0
              AND (@FInicioParsed IS NULL OR FechaVencimiento >= @FInicioParsed)
              AND (@FFinParsed IS NULL OR FechaVencimiento <= @FFinParsed)
            GROUP BY 
                CASE 
                    WHEN DATEDIFF(day, GETDATE(), FechaVencimiento) <= 30 THEN 'Vencen este mes (Urgente)'
                    WHEN DATEDIFF(day, GETDATE(), FechaVencimiento) <= 60 THEN 'Vencen el proximo mes'
                    WHEN DATEDIFF(day, GETDATE(), FechaVencimiento) <= 90 THEN 'De 2 a 3 meses'
                    WHEN DATEDIFF(day, GETDATE(), FechaVencimiento) <= 180 THEN 'De 3 a 6 meses'
                    ELSE 'Mas de 6 meses (Vigentes)'
                END
        ) v ON t.Rango = v.RangoClave;

        SELECT 
            Rango,
            ISNULL(CantidadLotes, 0) AS CantidadLotes,
            ISNULL(UnidadesTotales, 0) AS UnidadesTotales,
            NivelAlerta
        FROM #TmpTramosVenc
        ORDER BY Orden ASC;

        CREATE TABLE #TmpConsumoEsenciales (
            IdProducto INT,
            ConsumoPeriodo DECIMAL(12,2)
        );

        IF @FInicioParsed IS NOT NULL AND @FFinParsed IS NOT NULL
        BEGIN
            INSERT INTO #TmpConsumoEsenciales (IdProducto, ConsumoPeriodo)
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

        SELECT 
            MIN(ISNULL(b.IdProducto, 0)) AS IdProducto,
            LTRIM(RTRIM(e.Codigo)) AS Codigo,
            LTRIM(RTRIM(e.Nombre)) AS Nombre,
            ISNULL(NULLIF(LTRIM(RTRIM(e.FormaFarmaceutica)), ''), 'OTROS') AS FormaFarmaceutica,
            SUM(ISNULL(s.saldo, 0)) AS CantidadTotal,
            CASE 
                WHEN @FInicioParsed IS NOT NULL AND @FFinParsed IS NOT NULL 
                THEN ISNULL(ce.ConsumoPeriodo, 0)
                ELSE ISNULL(e.ConsumoPromedio, 0)
            END AS ConsumoMensual,
            CASE 
                WHEN SUM(ISNULL(s.saldo, 0)) <= 0 THEN 'DESABASTECIDO'
                WHEN SUM(ISNULL(s.saldo, 0)) < CASE 
                    WHEN @FInicioParsed IS NOT NULL AND @FFinParsed IS NOT NULL THEN ISNULL(ce.ConsumoPeriodo, 10)
                    ELSE ISNULL(e.ConsumoPromedio, 10)
                END THEN 'SUBSTOCK'
                WHEN SUM(ISNULL(s.saldo, 0)) <= (CASE 
                    WHEN @FInicioParsed IS NOT NULL AND @FFinParsed IS NOT NULL THEN ISNULL(ce.ConsumoPeriodo, 10)
                    ELSE ISNULL(e.ConsumoPromedio, 10)
                END * 6) THEN 'NORMOSTOCK'
                ELSE 'SOBRESTOCK'
            END AS EstadoSemaforo,
            e.EsDigemid,
            e.EsSis
        FROM dbo.Far_ItemsEvaluadosSisDigemid e
        LEFT JOIN dbo.FactCatalogoBienesInsumos b WITH (NOLOCK) ON b.Codigo = e.Codigo
        LEFT JOIN #TmpSaldos s ON s.idProducto = b.IdProducto
        LEFT JOIN #TmpConsumoEsenciales ce ON ce.IdProducto = b.IdProducto
        WHERE (
            @CanastaParam = 'TODOS'
            OR (@CanastaParam = 'DIGEMID' AND e.EsDigemid = 1)
            OR (@CanastaParam = 'SIS' AND e.EsSis = 1)
        )
        GROUP BY 
            e.Codigo, 
            e.Nombre, 
            e.FormaFarmaceutica, 
            e.ConsumoPromedio, 
            ce.ConsumoPeriodo,
            e.EsDigemid, 
            e.EsSis
        ORDER BY 
            CASE 
                WHEN @FInicioParsed IS NOT NULL AND @FFinParsed IS NOT NULL THEN
                    CASE 
                        WHEN ISNULL(ce.ConsumoPeriodo, 0) > 0 THEN 1 
                        ELSE 2 
                    END
                ELSE
                    CASE 
                        WHEN SUM(ISNULL(s.saldo, 0)) <= 0 THEN 1
                        WHEN SUM(ISNULL(s.saldo, 0)) < ISNULL(e.ConsumoPromedio, 10) THEN 2
                        ELSE 3
                    END
            END,
            CASE 
                WHEN @FInicioParsed IS NOT NULL AND @FFinParsed IS NOT NULL THEN ISNULL(ce.ConsumoPeriodo, 0)
                ELSE ISNULL(e.ConsumoPromedio, 0)
            END DESC;

        ;WITH ItemsEvaluadosAgrupados AS (
            SELECT 
                LTRIM(RTRIM(e.Codigo)) AS Codigo,
                ISNULL(NULLIF(LTRIM(RTRIM(e.FormaFarmaceutica)), ''), 'OTROS') AS FormaFarmaceutica,
                SUM(ISNULL(s.saldo, 0)) AS SaldoTotal
            FROM dbo.Far_ItemsEvaluadosSisDigemid e
            LEFT JOIN dbo.FactCatalogoBienesInsumos b WITH (NOLOCK) ON b.Codigo = e.Codigo
            LEFT JOIN #TmpSaldos s ON s.idProducto = b.IdProducto
            WHERE (
                @CanastaParam = 'TODOS'
                OR (@CanastaParam = 'DIGEMID' AND e.EsDigemid = 1)
                OR (@CanastaParam = 'SIS' AND e.EsSis = 1)
            )
            GROUP BY 
                e.Codigo,
                e.FormaFarmaceutica
        )
        SELECT TOP 8
            FormaFarmaceutica,
            COUNT(1) AS TotalItems,
            SUM(CASE WHEN SaldoTotal > 0 THEN 1 ELSE 0 END) AS ConStock,
            SUM(CASE WHEN SaldoTotal <= 0 THEN 1 ELSE 0 END) AS SinStock
        FROM ItemsEvaluadosAgrupados
        GROUP BY FormaFarmaceutica
        ORDER BY TotalItems DESC;

        DROP TABLE #TmpSaldos;
        DROP TABLE #TmpAlmacenes;
        DROP TABLE #TmpLotes;
        DROP TABLE #TmpTramosVenc;
        DROP TABLE #TmpConsumoEsenciales;

    END TRY
    BEGIN CATCH
        IF OBJECT_ID('tempdb..#TmpSaldos') IS NOT NULL DROP TABLE #TmpSaldos;
        IF OBJECT_ID('tempdb..#TmpAlmacenes') IS NOT NULL DROP TABLE #TmpAlmacenes;
        IF OBJECT_ID('tempdb..#TmpComparativa') IS NOT NULL DROP TABLE #TmpComparativa;
        IF OBJECT_ID('tempdb..#TmpLotes') IS NOT NULL DROP TABLE #TmpLotes;
        IF OBJECT_ID('tempdb..#TmpTramosVenc') IS NOT NULL DROP TABLE #TmpTramosVenc;
        IF OBJECT_ID('tempdb..#TmpConsumoEsenciales') IS NOT NULL DROP TABLE #TmpConsumoEsenciales;

        DECLARE @ErrMsg NVARCHAR(4000) = ERROR_MESSAGE();
        DECLARE @ErrSeverity INT = ERROR_SEVERITY();
        DECLARE @ErrState INT = ERROR_STATE();
        RAISERROR(@ErrMsg, @ErrSeverity, @ErrState);
    END CATCH
END;
GO
