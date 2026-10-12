SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

-- ====================================================================================================
-- Procedimiento : dbo.usp_go_Atenciones_ActualizarEstadoCE
-- Módulo        : Consulta Externa (CE) - Gestión de Estados Clínicos
-- Hospital      : Hospital Nacional Sergio E. Bernales (HNSEB)
-- Propósito     : Actualiza el estado de la atención médica (inicio de atención o cierre/egreso) y destino.
-- ====================================================================================================

CREATE OR ALTER PROCEDURE dbo.usp_go_Atenciones_ActualizarEstadoCE
    @IdAtencion        INT,
    @IdDestinoAtencion INT = 0,
    @CerrarAtencion    BIT = 0
AS
BEGIN
    SET NOCOUNT ON;

    BEGIN TRY
        BEGIN TRANSACTION;

        -- 1. Actualizar destino de atención si fue provisto
        IF (@IdDestinoAtencion IS NOT NULL AND @IdDestinoAtencion > 0)
        BEGIN
            UPDATE dbo.Atenciones
            SET IdDestinoAtencion = @IdDestinoAtencion
            WHERE IdAtencion = @IdAtencion;
        END

        -- 2. Actualizar estado y tiempos de atención
        IF (@CerrarAtencion = 1)
        BEGIN
            -- Cierre / Egreso de atención médica (Estado 2 = Atendido)
            UPDATE dbo.Atenciones
            SET HoraEgreso       = CONVERT(VARCHAR(5), GETDATE(), 108),
                idEstadoAtencion = 2
            WHERE IdAtencion = @IdAtencion;
        END
        ELSE
        BEGIN
            -- Inicio o mantenimiento de atención (Estado 1 = En Atención)
            UPDATE dbo.Atenciones
            SET HoraIngreso      = CASE 
                                     WHEN HoraIngreso IS NULL OR RTRIM(HoraIngreso) = '' 
                                     THEN CONVERT(VARCHAR(5), GETDATE(), 108) 
                                     ELSE HoraIngreso 
                                   END,
                idEstadoAtencion = CASE 
                                     WHEN idEstadoAtencion = 2 THEN 2 
                                     ELSE 1 
                                   END
            WHERE IdAtencion = @IdAtencion;
        END

        COMMIT TRANSACTION;

        SELECT 1 AS Exito, ISNULL(IdPaciente, 0) AS IdPaciente 
        FROM dbo.Atenciones WITH (NOLOCK) 
        WHERE IdAtencion = @IdAtencion;
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
