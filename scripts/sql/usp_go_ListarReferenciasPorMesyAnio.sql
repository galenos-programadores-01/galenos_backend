-- ============================================================================
-- Procedimiento: dbo.usp_go_ListarReferenciasPorMesyAnio
-- Módulo: Dashboard de Referencias y Contrarreferencias (HNSEB)
-- Base de Datos: SIGH / Galenos Pro
-- ============================================================================

SET ANSI_NULLS ON;
GO
SET QUOTED_IDENTIFIER ON;
GO

CREATE OR ALTER PROCEDURE dbo.usp_go_ListarReferenciasPorMesyAnio
    @mes INT = 0,
    @año INT = 0,
    @anio INT = 0,
    @estado INT = 0,
    @ups VARCHAR(10) = '',
    @opcion INT = 1
AS
BEGIN
    SET NOCOUNT ON;
    SET LANGUAGE spanish;

    -- Consolidar alias de anio / año
    DECLARE @AnioFiltro INT = @año;
    IF (@AnioFiltro IS NULL OR @AnioFiltro = 0) AND (@anio IS NOT NULL AND @anio <> 0)
        SET @AnioFiltro = @anio;
    IF @AnioFiltro IS NULL OR @AnioFiltro = 0
        SET @AnioFiltro = YEAR(GETDATE());

    IF @opcion = 1
    BEGIN
        SELECT 
            MONTH(FechaReferencia) AS MesNum,
            DATENAME(MONTH, FechaReferencia) AS Mes,
            COUNT(IdReferencia) AS Cantidad
        FROM dbo.Referencia WITH (NOLOCK)
        WHERE YEAR(FechaReferencia) = @AnioFiltro
          AND (@mes = 0 OR MONTH(FechaReferencia) = @mes)
        GROUP BY DATENAME(MONTH, FechaReferencia), MONTH(FechaReferencia)
        ORDER BY MONTH(FechaReferencia);
    END
    ELSE IF @opcion = 2
    BEGIN
        SELECT 
            ISNULL(es.Descripcion, 'DESCONOCIDO') AS Descripcion,
            COUNT(es.IdEstado) AS Conteo,
            MONTH(ref.FechaReferencia) AS MesNum,
            DATENAME(MONTH, ref.FechaReferencia) AS Mes,
            COUNT(ref.IdReferencia) AS Cantidad
        FROM dbo.Referencia ref WITH (NOLOCK)
        INNER JOIN dbo.EstadosReferencia es WITH (NOLOCK) ON ref.IdEstadoReferencia = es.IdEstado
        WHERE YEAR(ref.FechaReferencia) = @AnioFiltro
          AND (@mes = 0 OR MONTH(ref.FechaReferencia) = @mes)
          AND es.IdEstado IN (3, 5, 7)
        GROUP BY DATENAME(MONTH, ref.FechaReferencia), MONTH(ref.FechaReferencia), es.IdEstado, es.Descripcion
        ORDER BY MONTH(ref.FechaReferencia);
    END
    ELSE IF @opcion = 3
    BEGIN
        SELECT 
            ISNULL(es.Descripcion, 'DESCONOCIDO') AS Descripcion,
            COUNT(es.IdEstado) AS Conteo,
            MONTH(ref.FechaReferencia) AS MesNum,
            DATENAME(MONTH, ref.FechaReferencia) AS Mes,
            COUNT(ref.IdReferencia) AS Cantidad
        FROM dbo.Referencia ref WITH (NOLOCK)
        INNER JOIN dbo.EstadosReferencia es WITH (NOLOCK) ON ref.IdEstadoReferencia = es.IdEstado
        WHERE YEAR(ref.FechaReferencia) = @AnioFiltro
          AND (@mes = 0 OR MONTH(ref.FechaReferencia) = @mes)
          AND (@estado = 0 OR es.IdEstado = @estado)
        GROUP BY DATENAME(MONTH, ref.FechaReferencia), MONTH(ref.FechaReferencia), es.IdEstado, es.Descripcion
        ORDER BY MONTH(ref.FechaReferencia);
    END
    ELSE IF @opcion = 4
    BEGIN
        SELECT 
            ISNULL(es.Descripcion, 'DESCONOCIDO') AS Descripcion,
            COUNT(es.IdEstado) AS Conteo,
            MONTH(ref.FechaReferencia) AS MesNum,
            DATENAME(MONTH, ref.FechaReferencia) AS Mes,
            COUNT(ref.IdReferencia) AS Cantidad
        FROM dbo.Referencia ref WITH (NOLOCK)
        INNER JOIN dbo.EstadosReferencia es WITH (NOLOCK) ON ref.IdEstadoReferencia = es.IdEstado
        WHERE YEAR(ref.FechaReferencia) = @AnioFiltro
          AND (@mes = 0 OR MONTH(ref.FechaReferencia) = @mes)
          AND (LTRIM(RTRIM(@ups)) = '' OR ref.ServicioDestino = @ups)
          AND (@estado = 0 OR es.IdEstado = @estado)
        GROUP BY DATENAME(MONTH, ref.FechaReferencia), MONTH(ref.FechaReferencia), es.IdEstado, es.Descripcion
        ORDER BY MONTH(ref.FechaReferencia);
    END
END;
GO