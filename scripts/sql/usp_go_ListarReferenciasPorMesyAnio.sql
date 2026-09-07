-- usp_go_ListarReferenciasPorMesyAnio
-- Dashboard de referencias y contrarreferencias (gráfico de barras)
-- Uso: EXEC usp_go_ListarReferenciasPorMesyAnio @mes, @año, @estado, @ups, @opcion
-- Ejemplo: EXEC usp_go_ListarReferenciasPorMesyAnio 1, 2026, 0, '', 1
create proc usp_go_ListarReferenciasPorMesyAnio --'2026-01-01','2026-02-28',1
@mes int,
@año int,
@estado int,
@ups varchar(10),
@opcion int
as

set language spanish

if @opcion = 1
begin

selecT 
month(FechaReferencia),
DATENAME ( month , FechaReferencia )  as Mes,
count(IdReferencia) as Cantidad
from Referencia
where year(FechaReferencia)=@año
and (month(FechaReferencia)=@mes or @mes = 0)
group by DATENAME ( month , FechaReferencia ),month(FechaReferencia)
order by month(FechaReferencia)

end

if @opcion = 2
begin

selecT 
es.descripcion,
count(es.idestado),
month(FechaReferencia),
DATENAME ( month , FechaReferencia )  as Mes,
count(IdReferencia) as Cantidad
from Referencia ref inner join EstadosRefconReferencia es
on ref.IdEstadoReferencia=es.IdEstado
where year(FechaReferencia)=@año
and (month(FechaReferencia)=@mes or @mes = 0)
and es.IdEstado in (3,5,7)
group by DATENAME ( month , FechaReferencia ),month(FechaReferencia),
es.idestado,es.descripcion
order by month(FechaReferencia)

end

if @opcion = 3
begin

selecT 
es.descripcion,
count(es.idestado),
month(FechaReferencia),
DATENAME ( month , FechaReferencia )  as Mes,
count(IdReferencia) as Cantidad
from Referencia ref inner join EstadosRefconReferencia es
on ref.IdEstadoReferencia=es.IdEstado
where year(FechaReferencia)=@año
and (month(FechaReferencia)=@mes or @mes = 0)
and es.IdEstado =@estado
group by DATENAME ( month , FechaReferencia ),month(FechaReferencia),
es.idestado,es.descripcion
order by month(FechaReferencia)
end

if @opcion = 4
begin

selecT 
es.descripcion,
count(es.idestado),
month(FechaReferencia),
DATENAME ( month , FechaReferencia )  as Mes,
count(IdReferencia) as Cantidad
from Referencia ref inner join EstadosRefconReferencia es
on ref.IdEstadoReferencia=es.IdEstado
where year(FechaReferencia)=@año
and (month(FechaReferencia)=@mes or @mes = 0)
and ref.ServicioDestino=@ups
and es.IdEstado =@estado
group by DATENAME ( month , FechaReferencia ),month(FechaReferencia),
es.idestado,es.descripcion
order by month(FechaReferencia)
end