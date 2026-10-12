-- ============================================================================
-- Procedimiento: dbo.usp_go_FactCatalogoServicios_ListarPaginado
-- Módulo: Catálogos / Servicios Médicos y Tarifarios (HNSEB)
-- Objetivo: Listar servicios hospitalarios, exámenes y procedimientos CPT
--           con tarifas cruzadas de FactCatalogoServiciosHosp y paginación nativa.
-- Base de Datos: SIGH / Galenos Pro
-- Motor: SQL Server 2016 / 2019 / 2022
-- ============================================================================

SET ANSI_NULLS ON;
GO
SET QUOTED_IDENTIFIER ON;
GO

CREATE OR ALTER PROCEDURE dbo.usp_go_FactCatalogoServicios_ListarPaginado
    @Filtro VARCHAR(150) = NULL,
    @IdServicioGrupo INT = 0,
    @EsCPT BIT = NULL,
    @SoloCpt BIT = NULL,
    @EsQuirurgico BIT = NULL,
    @SoloQuirurgico BIT = NULL,
    @IdEstado INT = NULL,
    @Estado INT = NULL,
    @PageNumber INT = 1,
    @PageSize INT = 20
AS
BEGIN
    SET NOCOUNT ON;
    SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;

    BEGIN TRY
        -- Consolidar alias de parámetros para máxima compatibilidad con frontend y repositorios
        IF @EsCPT IS NULL AND @SoloCpt IS NOT NULL SET @EsCPT = @SoloCpt;
        IF @EsQuirurgico IS NULL AND @SoloQuirurgico IS NOT NULL SET @EsQuirurgico = @SoloQuirurgico;
        IF @IdEstado IS NULL AND @Estado IS NOT NULL SET @IdEstado = @Estado;
        IF @IdEstado IS NULL AND @Estado IS NULL SET @IdEstado = 1;

        IF @PageNumber < 1 SET @PageNumber = 1;
        IF @PageSize < 1 SET @PageSize = 20;

        DECLARE @FiltroLike VARCHAR(160) = NULL;
        IF @Filtro IS NOT NULL AND LTRIM(RTRIM(@Filtro)) <> ''
            SET @FiltroLike = '%' + LTRIM(RTRIM(@Filtro)) + '%';

        -- Cálculo del total de registros
        DECLARE @TotalRegistros INT = 0;

        SELECT @TotalRegistros = COUNT(1)
        FROM dbo.FactCatalogoServicios s WITH (NOLOCK)
        WHERE (@FiltroLike IS NULL OR (
                s.Nombre LIKE @FiltroLike 
                OR s.Codigo LIKE @FiltroLike 
                OR s.CodMINSA LIKE @FiltroLike 
                OR s.CPMS LIKE @FiltroLike 
                OR s.Nombre_Corto LIKE @FiltroLike
              ))
          AND (@IdServicioGrupo = 0 OR s.IdServicioGrupo = @IdServicioGrupo)
          AND (@EsCPT IS NULL OR s.EsCPT = CASE WHEN @EsCPT = 1 THEN 1 ELSE 0 END)
          AND (@EsQuirurgico IS NULL OR s.EsQuirurgico = CASE WHEN @EsQuirurgico = 1 THEN 1 ELSE 0 END)
          AND (@IdEstado IS NULL OR s.idEstado = @IdEstado);

        -- Retorno de la página solicitada
        SELECT 
            s.IdProducto,
            LTRIM(RTRIM(ISNULL(s.Codigo, ''))) AS Codigo,
            LTRIM(RTRIM(ISNULL(s.Codigo_HIS, ''))) AS CodigoHis,
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
            ISNULL(hp.PrecioUnitario, 0) AS PrecioParticular,
            ISNULL(hs.PrecioUnitario, 0) AS PrecioSis,
            ISNULL(s.idEstado, 0) AS IdEstado,
            CAST(CASE WHEN ISNULL(s.CargaEnFUA, 0) = 1 THEN 1 ELSE 0 END AS BIT) AS CargaEnFua,
            CAST(ISNULL(s.CargaEnHIS, 0) AS BIT) AS CargaEnHis,
            @TotalRegistros AS TotalRegistros
        FROM dbo.FactCatalogoServicios s WITH (NOLOCK)
        LEFT JOIN dbo.FactCatalogoServiciosGrupo g WITH (NOLOCK) 
            ON g.IdServicioGrupo = s.IdServicioGrupo
        LEFT JOIN dbo.FactCatalogoServiciosHosp hp WITH (NOLOCK) 
            ON hp.IdProducto = s.IdProducto AND hp.IdTipoFinanciamiento = 1 AND hp.Activo = 1
        LEFT JOIN dbo.FactCatalogoServiciosHosp hs WITH (NOLOCK) 
            ON hs.IdProducto = s.IdProducto AND hs.IdTipoFinanciamiento = 2 AND hs.Activo = 1
        WHERE (@FiltroLike IS NULL OR (
                s.Nombre LIKE @FiltroLike 
                OR s.Codigo LIKE @FiltroLike 
                OR s.CodMINSA LIKE @FiltroLike 
                OR s.CPMS LIKE @FiltroLike 
                OR s.Nombre_Corto LIKE @FiltroLike
              ))
          AND (@IdServicioGrupo = 0 OR s.IdServicioGrupo = @IdServicioGrupo)
          AND (@EsCPT IS NULL OR s.EsCPT = CASE WHEN @EsCPT = 1 THEN 1 ELSE 0 END)
          AND (@EsQuirurgico IS NULL OR s.EsQuirurgico = CASE WHEN @EsQuirurgico = 1 THEN 1 ELSE 0 END)
          AND (@IdEstado IS NULL OR s.idEstado = @IdEstado)
        ORDER BY s.IdProducto DESC
        OFFSET (@PageNumber - 1) * @PageSize ROWS
        FETCH NEXT @PageSize ROWS ONLY;

    END TRY
    BEGIN CATCH
        DECLARE @ErrMsg NVARCHAR(4000) = ERROR_MESSAGE();
        DECLARE @ErrSeverity INT = ERROR_SEVERITY();
        DECLARE @ErrState INT = ERROR_STATE();
        RAISERROR(@ErrMsg, @ErrSeverity, @ErrState);
    END CATCH
END;
GO
