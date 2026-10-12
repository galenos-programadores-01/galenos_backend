SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

-- ====================================================================================================
-- Procedimiento : dbo.usp_go_MaestroHisLab_Listar
-- Módulo        : Consulta Externa (CE) - Catálogos Diagnósticos HIS
-- Hospital      : Hospital Nacional Sergio E. Bernales (HNSEB)
-- Propósito     : Lista los códigos y valores del maestro HIS Laboratorio/Complementario para diagnósticos.
-- ====================================================================================================

CREATE OR ALTER PROCEDURE dbo.usp_go_MaestroHisLab_Listar
AS
BEGIN
    SET NOCOUNT ON;

    BEGIN TRY
        SELECT 
            Codigo_Lab                             AS Codigo,
            RTRIM(LTRIM(Valor_Lab))                AS Valor,
            ISNULL(RTRIM(LTRIM(Descripcion_Lab)), '') AS Descripcion
        FROM dbo.MAESTRO_HIS_LAB WITH (NOLOCK)
        WHERE Valor_Lab IS NOT NULL 
          AND RTRIM(LTRIM(Valor_Lab)) != ''
        ORDER BY Codigo_Lab ASC, Valor_Lab ASC;
    END TRY
    BEGIN CATCH
        DECLARE @ErrorMessage NVARCHAR(4000) = ERROR_MESSAGE();
        DECLARE @ErrorSeverity INT = ERROR_SEVERITY();
        DECLARE @ErrorState INT = ERROR_STATE();
        RAISERROR (@ErrorMessage, @ErrorSeverity, @ErrorState);
    END CATCH
END;
GO
