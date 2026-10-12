CREATE OR ALTER PROCEDURE dbo.usp_go_Recetas_ListarCabecerasPorCuenta
    @IdCuentaAtencion INT
AS
BEGIN
    SET NOCOUNT ON;

    SELECT 
        rc.idReceta, 
        CONVERT(varchar(19), rc.FechaReceta, 120) AS FechaReceta,
        RTRIM(ISNULL(em.ApellidoPaterno,'')) + ' ' + RTRIM(ISNULL(em.ApellidoMaterno,'')) + ' ' + RTRIM(ISNULL(em.Nombres,'')) AS Medico
    FROM dbo.RecetaCabecera rc WITH (NOLOCK)
    LEFT JOIN dbo.Medicos m WITH (NOLOCK) ON m.IdMedico = rc.idMedicoReceta
    LEFT JOIN dbo.Empleados em WITH (NOLOCK) ON em.IdEmpleado = m.IdEmpleado
    WHERE rc.idCuentaAtencion = @IdCuentaAtencion;
END;
GO
