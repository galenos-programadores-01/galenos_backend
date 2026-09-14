CREATE PROCEDURE usp_go_ListarEstablecimientos
AS
	SELECT Codigo, Nombre
	FROM Establecimientos