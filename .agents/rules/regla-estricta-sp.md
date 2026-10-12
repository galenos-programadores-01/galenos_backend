# Regla Estricta de Persistencia: 0 Consultas SQL Directas y Uso Exclusivo de Stored Procedures

1. **PROHIBIDO EL CÓDIGO SQL DIRECTO (0 SELECT, 0 UPDATE, 0 INSERT, 0 DELETE)**:
   - Jamás escribir consultas `SELECT`, `UPDATE`, `INSERT` o `DELETE` inline en el código de los repositorios de Go.
   - Toda interacción con la base de datos SQL Server debe realizarse únicamente invocando Stored Procedures mediante la sintaxis `EXEC dbo.<nombre_sp> ...`.

2. **PROHIBIDO CREAR O INVENTAR STORED PROCEDURES EN LA BASE DE DATOS**:
   - La IA nunca debe ejecutar `CREATE PROCEDURE` en la base de datos del hospital.
   - Solo se deben consumir los Stored Procedures existentes en el motor SQL Server.

3. **PROCEDIMIENTO PARA SPS FALTANTES**:
   - Si se requiere una consulta o métrica agregada para la cual no existe un Stored Procedure en la base de datos:
     1. NO intentar reemplazarlo con una consulta SQL directa.
     2. NO crear el SP en la base de datos.
     3. Definir la llamada esperada en el repositorio (`EXEC dbo.usp_go_...`) y documentar detalladamente la especificación técnica del Stored Procedure faltante (parámetros de entrada, tipos de datos y columnas esperadas de salida) para que el usuario o el equipo de DBA lo cree de acuerdo con las políticas institucionales.
