CREATE OR ALTER PROCEDURE dbo.usp_go_Pacientes_ListarPaginado
    @Offset INT = 0,
    @PageSize INT = 20
AS
BEGIN
    SET NOCOUNT ON;

    SELECT COUNT(1) AS TotalRegistros
    FROM dbo.Pacientes WITH (NOLOCK)
    WHERE ISNULL(NroDocumento, '') <> '';

    SELECT 
        IdPaciente, 
        NroDocumento, 
        NroHistoriaClinica, 
        ApellidoPaterno,
        ApellidoMaterno, 
        PrimerNombre, 
        SegundoNombre, 
        TercerNombre, 
        FechaNacimiento, 
        IdTipoSexo
    FROM dbo.Pacientes WITH (NOLOCK)
    WHERE ISNULL(NroDocumento, '') <> '' 
    ORDER BY NroDocumento
    OFFSET @Offset ROWS FETCH NEXT @PageSize ROWS ONLY;
END;
GO
