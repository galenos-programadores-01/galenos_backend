-- ============================================================================
-- Procedimiento: dbo.usp_go_FactCatalogoServicios_ObtenerPorId
-- Módulo: Catálogos / Servicios Médicos y Tarifarios (HNSEB)
-- Objetivo: Obtener la ficha del servicio con sus códigos de salud (Recordset 1)
--           y su desglose de tarifas por financiamiento (Recordset 2).
-- Base de Datos: SIGH / Galenos Pro
-- Motor: SQL Server 2016 / 2019 / 2022
-- ============================================================================

SET ANSI_NULLS ON;
GO
SET QUOTED_IDENTIFIER ON;
GO

CREATE OR ALTER PROCEDURE dbo.usp_go_FactCatalogoServicios_ObtenerPorId
    @IdProducto INT
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

        -- ---------------------------------------------------------------------
        -- RECORDSET 1: FICHA PRINCIPAL DEL SERVICIO
        -- ---------------------------------------------------------------------
        SELECT 
            s.IdProducto,
            LTRIM(RTRIM(ISNULL(s.Codigo, ''))) AS Codigo,
            LTRIM(RTRIM(ISNULL(s.Codigo_HIS, ''))) AS CodigoHis,
            LTRIM(RTRIM(ISNULL(s.CodigoSunat, ''))) AS CodigoSunat,
            LTRIM(RTRIM(ISNULL(s.Nombre, ''))) AS Nombre,
            LTRIM(RTRIM(ISNULL(s.Nombre_Corto, ''))) AS NombreCorto,
            ISNULL(s.IdServicioGrupo, 0) AS IdServicioGrupo,
            LTRIM(RTRIM(ISNULL(g.Descripcion, 'SIN GRUPO'))) AS GrupoNombre,
            LTRIM(RTRIM(ISNULL(s.CodMINSA, ''))) AS CodMinsa,
            LTRIM(RTRIM(ISNULL(s.NombreMINSA, ''))) AS NombreMinsa,
            CAST(CASE WHEN ISNULL(s.EsCPT, 0) = 1 THEN 1 ELSE 0 END AS BIT) AS EsCpt,
            CAST(CASE WHEN ISNULL(s.EsQuirurgico, 0) = 1 THEN 1 ELSE 0 END AS BIT) AS EsQuirurgico,
            LTRIM(RTRIM(ISNULL(s.codigoSIS, ''))) AS CodigoSis,
            LTRIM(RTRIM(ISNULL(s.CodigoSOAT, ''))) AS CodigoSoat,
            LTRIM(RTRIM(ISNULL(s.CPMS, ''))) AS Cpms,
            ISNULL(s.idEstado, 0) AS IdEstado,
            CAST(CASE WHEN ISNULL(s.CargaEnFUA, 0) = 1 THEN 1 ELSE 0 END AS BIT) AS CargaEnFua,
            CAST(ISNULL(s.CargaEnHIS, 0) AS BIT) AS CargaEnHis
        FROM dbo.FactCatalogoServicios s WITH (NOLOCK)
        LEFT JOIN dbo.FactCatalogoServiciosGrupo g WITH (NOLOCK) 
            ON g.IdServicioGrupo = s.IdServicioGrupo
        WHERE s.IdProducto = @IdProducto;

        -- ---------------------------------------------------------------------
        -- RECORDSET 2: TARIFAS POR FINANCIAMIENTO (FactCatalogoServiciosHosp)
        -- ---------------------------------------------------------------------
        SELECT 
            h.IdTipoFinanciamiento,
            LTRIM(RTRIM(ISNULL(f.Descripcion, 'FINANCIAMIENTO ' + CAST(h.IdTipoFinanciamiento AS VARCHAR(10))))) AS FinanciamientoNombre,
            ISNULL(h.PrecioUnitario, 0) AS PrecioUnitario,
            CAST(ISNULL(h.Activo, 0) AS BIT) AS Activo
        FROM dbo.FactCatalogoServiciosHosp h WITH (NOLOCK)
        LEFT JOIN dbo.TiposFinanciamiento f WITH (NOLOCK)
            ON f.IdTipoFinanciamiento = h.IdTipoFinanciamiento
        WHERE h.IdProducto = @IdProducto
        ORDER BY h.IdTipoFinanciamiento ASC;

    END TRY
    BEGIN CATCH
        DECLARE @ErrMsg NVARCHAR(4000) = ERROR_MESSAGE();
        DECLARE @ErrSeverity INT = ERROR_SEVERITY();
        DECLARE @ErrState INT = ERROR_STATE();
        RAISERROR(@ErrMsg, @ErrSeverity, @ErrState);
    END CATCH
END;
GO
