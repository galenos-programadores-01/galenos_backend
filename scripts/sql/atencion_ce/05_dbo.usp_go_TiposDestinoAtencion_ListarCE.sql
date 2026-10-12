SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

-- ====================================================================================================
-- Procedimiento : dbo.usp_go_TiposDestinoAtencion_ListarCE
-- Módulo        : Consulta Externa (CE) - Destinos de Atención
-- Hospital      : Hospital Nacional Sergio E. Bernales (HNSEB)
-- Propósito     : Lista los destinos de atención / egreso clínico válidos para el servicio de Consulta Externa.
-- ====================================================================================================

CREATE OR ALTER PROCEDURE dbo.usp_go_TiposDestinoAtencion_ListarCE
AS
BEGIN
    SET NOCOUNT ON;

    BEGIN TRY
        SELECT 
            IdDestinoAtencion,
            RTRIM(LTRIM(ISNULL(Codigo, ''))) + ' - ' + RTRIM(LTRIM(ISNULL(Descripcion, ''))) AS Descripcion
        FROM dbo.TiposDestinoAtencion WITH (NOLOCK)
        WHERE IdTipoServicio = 1
        ORDER BY Codigo ASC, Descripcion ASC;
    END TRY
    BEGIN CATCH
        DECLARE @ErrorMessage NVARCHAR(4000) = ERROR_MESSAGE();
        DECLARE @ErrorSeverity INT = ERROR_SEVERITY();
        DECLARE @ErrorState INT = ERROR_STATE();
        RAISERROR (@ErrorMessage, @ErrorSeverity, @ErrorState);
    END CATCH
END;
GO
