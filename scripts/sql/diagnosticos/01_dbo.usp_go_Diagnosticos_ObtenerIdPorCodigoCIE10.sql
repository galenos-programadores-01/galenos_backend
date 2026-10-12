CREATE OR ALTER PROCEDURE dbo.usp_go_Diagnosticos_ObtenerIdPorCodigoCIE10
    @CodigoCIE10 VARCHAR(20)
AS
BEGIN
    SET NOCOUNT ON;

    DECLARE @IdDiagnostico INT = 0;
    SET @CodigoCIE10 = RTRIM(LTRIM(@CodigoCIE10));

    SELECT TOP 1 @IdDiagnostico = IdDiagnostico
    FROM dbo.Diagnosticos WITH (NOLOCK)
    WHERE ISNULL(CodigoCIE10, CodigoCIE2004) = @CodigoCIE10;

    IF @IdDiagnostico = 0
    BEGIN
        SELECT TOP 1 @IdDiagnostico = IdDiagnostico
        FROM dbo.Diagnosticos WITH (NOLOCK)
        WHERE REPLACE(ISNULL(CodigoCIE10, CodigoCIE2004), '.', '') = REPLACE(@CodigoCIE10, '.', '');
    END

    SELECT ISNULL(@IdDiagnostico, 0) AS IdDiagnostico;
END;
GO
