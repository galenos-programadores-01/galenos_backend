SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

-- ====================================================================================================
-- Procedimiento : dbo.usp_go_AtencionesTriaje_ObtenerVigentePorAtencion
-- Módulo        : Consulta Externa (CE) - Triaje
-- Hospital      : Hospital Nacional Sergio E. Bernales (HNSEB)
-- Propósito     : Recupera los signos vitales vigentes del último triaje realizado a una atención.
-- ====================================================================================================

CREATE OR ALTER PROCEDURE dbo.usp_go_AtencionesTriaje_ObtenerVigentePorAtencion
    @IdAtencion INT
AS
BEGIN
    SET NOCOUNT ON;

    BEGIN TRY
        SELECT TOP 1
            ISNULL(PresionArterial, '')        AS PresionArterial,
            ISNULL(TRY_CAST(FrecCardiaca AS INT), 0)         AS FrecuenciaCardiaca,
            ISNULL(TRY_CAST(FrecRespiratoria AS INT), 0)     AS FrecuenciaRespiratoria,
            ISNULL(TRY_CAST(Temperatura AS DECIMAL(5,2)), 0.0) AS Temperatura,
            ISNULL(TRY_CAST(SAT02 AS INT), 0)                AS Saturacion,
            ISNULL(TRY_CAST(Peso AS DECIMAL(5,2)), 0.0)      AS Peso,
            ISNULL(TRY_CAST(Talla AS DECIMAL(5,2)), 0.0)     AS Talla,
            ISNULL(TRY_CAST(IMC AS DECIMAL(5,2)), 0.0)       AS IMC,
            Fecha                              AS FechaRegistro,
            ISNULL(Observacion, '')            AS Observacion
        FROM dbo.AtencionesTriaje WITH (NOLOCK)
        WHERE IdAtencion = @IdAtencion
          AND UltimoTriaje = 1
        ORDER BY IdTriaje DESC;
    END TRY
    BEGIN CATCH
        DECLARE @ErrorMessage NVARCHAR(4000) = ERROR_MESSAGE();
        DECLARE @ErrorSeverity INT = ERROR_SEVERITY();
        DECLARE @ErrorState INT = ERROR_STATE();
        RAISERROR (@ErrorMessage, @ErrorSeverity, @ErrorState);
    END CATCH
END;
GO
