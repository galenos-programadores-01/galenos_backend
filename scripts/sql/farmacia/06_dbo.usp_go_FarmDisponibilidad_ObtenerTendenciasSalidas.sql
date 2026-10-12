-- ============================================================================
-- Procedimiento: dbo.usp_go_FarmDisponibilidad_ObtenerTendenciasSalidas
-- Módulo: Farmacia / Disponibilidad de Medicamentos y Salidas (HNSEB)
-- Objetivo: Retornar la serie temporal mensual agregada de salidas/dispensación
--           por modalidad (SIS, Ventas, Programas, Intervenciones Sanitarias).
-- Base de Datos: SIGH / Galenos Pro
-- Motor: SQL Server 2016 / 2019 / 2022
-- ============================================================================

SET ANSI_NULLS ON;
GO
SET QUOTED_IDENTIFIER ON;
GO

CREATE OR ALTER PROCEDURE dbo.usp_go_FarmDisponibilidad_ObtenerTendenciasSalidas
    @MesesAtras INT = 6,
    @IdAlmacen INT = 0
AS
BEGIN
    SET NOCOUNT ON;
    SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;

    BEGIN TRY
        IF @MesesAtras IS NULL OR @MesesAtras <= 0
            SET @MesesAtras = 6;

        DECLARE @FechaLimite DATETIME = DATEADD(month, -@MesesAtras, GETDATE());

        IF OBJECT_ID('dbo.farmMovimiento') IS NOT NULL AND OBJECT_ID('dbo.farmMovimientoDetalle') IS NOT NULL
        BEGIN
            SELECT 
                CONVERT(VARCHAR(7), m.fechaCreacion, 120) AS Periodo,
                CASE 
                    WHEN ISNULL(d.IdTipoSalidaBienInsumo, 0) = 1 THEN 'Dispensación SIS (Asegurados)'
                    WHEN ISNULL(d.IdTipoSalidaBienInsumo, 0) = 2 THEN 'Venta Contado (Particulares)'
                    WHEN ISNULL(d.IdTipoSalidaBienInsumo, 0) = 3 THEN 'Programas Estratégicos (MINSA)'
                    ELSE 'Distribución Interna a Servicios'
                END AS ModalidadDispensacion,
                ISNULL(SUM(CAST(d.Cantidad AS BIGINT)), 0) AS TotalUnidadesDispensadas
            FROM dbo.farmMovimiento m WITH (NOLOCK)
            INNER JOIN dbo.farmMovimientoDetalle d WITH (NOLOCK) 
                ON m.MovNumero = d.MovNumero AND m.MovTipo = d.MovTipo
            WHERE m.MovTipo = 'S' -- Solo salidas
              AND m.fechaCreacion >= @FechaLimite
              AND (@IdAlmacen = 0 OR m.idAlmacenOrigen = @IdAlmacen)
              AND ISNULL(m.idEstadoMovimiento, 1) != 3 -- Excluir anulados si corresponde
            GROUP BY 
                CONVERT(VARCHAR(7), m.fechaCreacion, 120),
                CASE 
                    WHEN ISNULL(d.IdTipoSalidaBienInsumo, 0) = 1 THEN 'Dispensación SIS (Asegurados)'
                    WHEN ISNULL(d.IdTipoSalidaBienInsumo, 0) = 2 THEN 'Venta Contado (Particulares)'
                    WHEN ISNULL(d.IdTipoSalidaBienInsumo, 0) = 3 THEN 'Programas Estratégicos (MINSA)'
                    ELSE 'Distribución Interna a Servicios'
                END
            ORDER BY Periodo ASC, ModalidadDispensacion ASC;
        END
        ELSE
        BEGIN
            SELECT 
                CONVERT(VARCHAR(7), DATEADD(month, -m.Mes, GETDATE()), 120) AS Periodo,
                v.Modalidad AS ModalidadDispensacion,
                CAST(0.00 AS DECIMAL(14,2)) AS TotalUnidadesDispensadas
            FROM (VALUES (5), (4), (3), (2), (1), (0)) AS m(Mes)
            CROSS JOIN (
                VALUES 
                    ('Dispensación SIS (Asegurados)'),
                    ('Venta Contado (Particulares)'),
                    ('Programas Estratégicos (MINSA)'),
                    ('Distribución Interna a Servicios')
            ) AS v(Modalidad)
            ORDER BY Periodo ASC, ModalidadDispensacion ASC;
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
