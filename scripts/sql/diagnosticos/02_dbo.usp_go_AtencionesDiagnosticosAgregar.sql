SET ANSI_NULLS ON
GO
SET QUOTED_IDENTIFIER ON
GO

-- ====================================================================================================
-- Procedimiento : dbo.usp_go_AtencionesDiagnosticosAgregar
-- Módulo        : Consulta Externa (CE) - Diagnósticos de Atención
-- Hospital      : Hospital Nacional Sergio E. Bernales (HNSEB)
-- Propósito     : Registra un diagnóstico para la atención con correlativo y subclasificación HIS.
--                 Mapea correctamente labConfHIS_1, labConfHIS_2, labConfHIS_3 (varchar) y sus IDs con TRY_CONVERT.
-- ====================================================================================================

CREATE OR ALTER PROCEDURE [dbo].[usp_go_AtencionesDiagnosticosAgregar]  
(  
    @IdAtencionDiagnostico        VARCHAR(100) OUTPUT,  
    @IdClasificacionDx            INT,
    @IdSubclasificacionDiagnostico INT,  
    @IdDiagnostico                INT,  
    @IdAtencion                   INT,  
    @LabConfHIS                   VARCHAR(3),  
    @LabConfHIS_1                 VARCHAR(3),  
    @LabConfHIS_2                 VARCHAR(3),  
    @LabConfHIS_3                 VARCHAR(3),  
    @GrupoHIS                     INT,  
    @SubGrupoHIS                  INT,  
    @IdEpisodio                   INT,  
    @IdEvolucion                  INT,  
    @IdPrimeraAtencion            INT,  
    @IdPaciente                   INT  
)  
AS  
BEGIN  
    SET NOCOUNT ON;  

    DECLARE  
        @DescripcionSubclasificacion VARCHAR(50),  
        @NumeroOrdenCorrelativo      INT,  
        @IdAtencionDiagnosticoGenerado INT;  

    BEGIN TRY  
        SET @IdAtencionDiagnostico = NULL;  

        IF ISNULL(@IdAtencion, 0) = 0  
        BEGIN  
            SET @IdAtencionDiagnostico = '-1';  
            RETURN;  
        END;  

        IF ISNULL(@IdDiagnostico, 0) = 0  
        BEGIN  
            SET @IdAtencionDiagnostico = '-1';  
            RETURN;  
        END;  

        SELECT TOP (1)  
            @IdAtencionDiagnostico = CAST(AD.IdAtencionDiagnostico AS VARCHAR(100))  
        FROM dbo.AtencionesDiagnosticos AS AD  
        WHERE AD.IdAtencion = @IdAtencion  
          AND AD.IdDiagnostico = @IdDiagnostico;  

        IF @IdAtencionDiagnostico IS NOT NULL  
        BEGIN  
            SET @IdAtencionDiagnostico = '-1';  
            RETURN;  
        END;  

        SELECT  
            @NumeroOrdenCorrelativo = ISNULL(MAX(AD.NumOrdenCorelativo), 0) + 1  
        FROM dbo.AtencionesDiagnosticos AS AD WITH (UPDLOCK, HOLDLOCK)  
        WHERE AD.IdAtencion = @IdAtencion;  

        INSERT INTO dbo.AtencionesDiagnosticos  
        (  
            IdSubclasificacionDx,  
            IdClasificacionDx,  
            IdDiagnostico,  
            IdAtencion,  
            IdPrimeraAtencion,  
            labConfHIS,  
            labConfHIS_1,  
            labConfHIS_2,  
            labConfHIS_3,  
            id_labConfHIS_1,  
            id_labConfHIS_2,  
            id_labConfHIS_3,  
            GrupoHIS,  
            SubGrupoHIS,  
            Idepisodio,  
            NumOrdenCorelativo  
        )  
        VALUES  
        (  
            @IdSubclasificacionDiagnostico,  
            @IdClasificacionDx,  
            @IdDiagnostico,  
            @IdAtencion,  
            @IdPrimeraAtencion,  
            @LabConfHIS,  
            @LabConfHIS_1,  
            @LabConfHIS_2,  
            @LabConfHIS_3,  
            TRY_CONVERT(INT, @LabConfHIS_1),  
            TRY_CONVERT(INT, @LabConfHIS_2),  
            TRY_CONVERT(INT, @LabConfHIS_3),  
            @GrupoHIS,  
            @SubGrupoHIS,  
            @IdEpisodio,  
            @NumeroOrdenCorrelativo  
        );  

        SET @IdAtencionDiagnosticoGenerado = CONVERT(INT, SCOPE_IDENTITY());  

        SELECT TOP (1)  
            @DescripcionSubclasificacion = SC.Descripcion  
        FROM dbo.Atenciones AS A  
        INNER JOIN dbo.AtencionesDiagnosticos AS AD  
            ON AD.IdAtencion = A.IdAtencion  
        INNER JOIN dbo.SubclasificacionDiagnosticos AS SC  
            ON SC.IdSubclasificacionDx = AD.IdSubclasificacionDx  
        WHERE A.IdPaciente = @IdPaciente  
          AND A.IdAtencion <> @IdAtencion  
          AND AD.IdDiagnostico = @IdDiagnostico  
        ORDER BY A.IdAtencion DESC;  

        SET @IdAtencionDiagnostico =  
            CAST(@IdAtencionDiagnosticoGenerado AS VARCHAR(50))  
            + ';'  
            + ISNULL(@DescripcionSubclasificacion, '');  
    END TRY  
    BEGIN CATCH  
        SET @IdAtencionDiagnostico = 'ERROR;' + ERROR_MESSAGE();  
        RETURN;  
    END CATCH  
END;
GO
