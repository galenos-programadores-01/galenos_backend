-- usp_go_ListarUPs
-- Lista las Unidades Productoras de Servicios (UPS) activas
create proc usp_go_ListarUPs
as
selecT distinct su.Codigo, su.Descripcion from Servicios S inner join SuSalud_ups su
on s.codigoServicioSuSalud=su.Codigo
where s.IdTipoServicio=1 and s.idEstado=1