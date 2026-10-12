SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

-- ====================================================================================================
-- Procedimiento : dbo.usp_go_AtencionesDiagnosticos_EliminarPorAtencion
-- Módulo        : Consulta Externa (CE) - Diagnósticos de Atención
-- Hospital      : Hospital Nacional Sergio E. Bernales (HNSEB)
-- Propósito     : Elimina los diagnósticos previos asignados a una atención para permitir su actualización integral.
-- ====================================================================================================

CREATE OR ALTER PROCEDURE dbo.usp_go_AtencionesDiagnosticos_EliminarPorAtencion
    @IdAtencion          INT,
    @IdUsuarioAuditoria  INT = 1
AS
BEGIN
    SET NOCOUNT ON;

    BEGIN TRY
        BEGIN TRANSACTION;

        DELETE FROM dbo.AtencionesDiagnosticos
        WHERE IdAtencion = @IdAtencion;

        DECLARE @FilasEliminadas INT = @@ROWCOUNT;

        COMMIT TRANSACTION;

        SELECT @FilasEliminadas AS FilasEliminadas;
    END TRY
    BEGIN CATCH
        IF @@TRANCOUNT > 0
            ROLLBACK TRANSACTION;

        DECLARE @ErrorMessage NVARCHAR(4000) = ERROR_MESSAGE();
        DECLARE @ErrorSeverity INT = ERROR_SEVERITY();
        DECLARE @ErrorState INT = ERROR_STATE();
        RAISERROR (@ErrorMessage, @ErrorSeverity, @ErrorState);
    END CATCH
END;
GO
