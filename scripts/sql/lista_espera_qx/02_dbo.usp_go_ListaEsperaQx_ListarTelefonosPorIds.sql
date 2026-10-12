CREATE OR ALTER PROCEDURE dbo.usp_go_ListaEsperaQx_ListarTelefonosPorIds
    @Ids VARCHAR(MAX)
AS
BEGIN
    SET NOCOUNT ON;

    SELECT q.Id, q.Telefono
    FROM dbo.ListaEsperaQx q WITH (NOLOCK)
    INNER JOIN STRING_SPLIT(@Ids, ',') s ON TRY_CAST(RTRIM(LTRIM(s.value)) AS INT) = q.Id
    WHERE ISNULL(q.Telefono, '') <> '';
END;
GO
