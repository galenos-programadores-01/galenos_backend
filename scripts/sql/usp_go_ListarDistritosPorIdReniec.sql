CREATE PROCEDURE usp_go_ListarDistritosPorIdReniec
	@IdReniec int
AS
	SELECT IdDistrito, Nombre, IdReniec, IdProvincia
	FROM Distritos
	WHERE IdReniec = @IdReniec