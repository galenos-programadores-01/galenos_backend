-- ============================================================================
-- Procedimiento: dbo.usp_go_FactInsumosGrupoFarmacologico_ListarTodos
-- Módulo: Catálogos / Bienes e Insumos Médicos (HNSEB)
-- Objetivo: Listar todos los grupos farmacológicos para los selectores de búsqueda.
-- Base de Datos: SIGH / Galenos Pro
-- Motor: SQL Server 2016 / 2019 / 2022
-- ============================================================================

SET ANSI_NULLS ON;
GO
SET QUOTED_IDENTIFIER ON;
GO

CREATE OR ALTER PROCEDURE dbo.usp_go_FactInsumosGrupoFarmacologico_ListarTodos
AS
BEGIN
    SET NOCOUNT ON;
    SET TRANSACTION ISOLATION LEVEL READ UNCOMMITTED;

    BEGIN TRY
        SELECT 
            IdGrupoFarmacologico, 
            LTRIM(RTRIM(ISNULL(Descripcion, ''))) AS Descripcion
        FROM dbo.FactInsumosGrupoFarmacologico WITH (NOLOCK)
        ORDER BY Descripcion ASC;
    END TRY
    BEGIN CATCH
        DECLARE @ErrMsg NVARCHAR(4000) = ERROR_MESSAGE();
        DECLARE @ErrSeverity INT = ERROR_SEVERITY();
        DECLARE @ErrState INT = ERROR_STATE();
        RAISERROR(@ErrMsg, @ErrSeverity, @ErrState);
    END CATCH
END;
GO
