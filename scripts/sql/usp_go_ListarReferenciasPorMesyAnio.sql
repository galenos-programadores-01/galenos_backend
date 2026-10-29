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
    @FechaInicio DATE,
    @FechaFin DATE,
    @estadoreferencia INT = 0
AS
BEGIN
    SET NOCOUNT ON;
    SET LANGUAGE spanish;

    SELECT
        r.Descripcion,
        COUNT(1) AS Conteo,
        MONTH(ref.FechaReferencia) AS MesNum,
        DATENAME(MONTH, ref.FechaReferencia) AS Mes,
        COUNT(ref.IdReferencia) AS Cantidad
    FROM dbo.Referencia ref WITH (NOLOCK)
    INNER JOIN dbo.EstadosRefconReferencia r ON ref.IdEstadoReferencia = r.IdEstado
    WHERE CONVERT(DATE, ref.FechaReferencia) BETWEEN @FechaInicio AND @FechaFin
      AND (ref.IdEstadoReferencia = @estadoreferencia OR @estadoreferencia = 0)
    GROUP BY
        r.Descripcion,
        DATENAME(MONTH, ref.FechaReferencia),
        MONTH(ref.FechaReferencia)
    ORDER BY MONTH(ref.FechaReferencia);
END;
GO
