CREATE OR ALTER PROCEDURE [dbo].[usp_go_RecetaDetalleAgregar]
    @Mensaje VARCHAR(150) OUTPUT,  
    @idReceta int,
    @IdProducto int,
    @Cantidad int,
    @Precio money,
    @SaldoEnRegistroReceta int,
    @idDosisRecetada int,
    @observaciones varchar(300),
    @IdViaAdministracion int,
    @CodigoDiagnostico varchar(20),
    @Justificacion text,
    @idUNIDDosisReceta int,
    @idFrecuencia int,
    @DescripcionadicionalReceta varchar(300),
    @Duracion int,
    @idCuentaAtencionProxCita INT
AS
BEGIN
    SET NOCOUNT ON;

    IF NOT EXISTS (
        SELECT 1
        FROM [dbo].[RecetaDetalle] WITH (NOLOCK)
        WHERE [idReceta] = @idReceta AND [idItem] = @IdProducto
    )
    BEGIN
        DECLARE @Total MONEY = @Cantidad * @Precio;
        DECLARE @xCodigoDiagnostico VARCHAR(20) = NULL;
        DECLARE @xIdDiagnostico INT = NULL;

        IF ISNUMERIC(@CodigoDiagnostico) = 1 AND CHARINDEX('.', @CodigoDiagnostico) = 0
        BEGIN
            SET @xIdDiagnostico = TRY_CAST(@CodigoDiagnostico AS INT);
            SELECT TOP 1 @xCodigoDiagnostico = RTRIM(LTRIM(d.CodigoCIE10))
            FROM [dbo].[Diagnosticos] d WITH (NOLOCK)
            WHERE d.IdDiagnostico = @xIdDiagnostico;
        END
        ELSE
        BEGIN
            SET @xCodigoDiagnostico = RTRIM(LTRIM(@CodigoDiagnostico));
            SELECT TOP 1 @xIdDiagnostico = d.IdDiagnostico
            FROM [dbo].[Diagnosticos] d WITH (NOLOCK)
            WHERE RTRIM(LTRIM(d.CodigoCIE10)) = @xCodigoDiagnostico;
        END

        INSERT INTO [dbo].[RecetaDetalle] (
            [idReceta],
            [idItem],
            [CantidadPedida],
            [Precio],
            [Total],
            [SaldoEnRegistroReceta],
            [SaldoEnDespachoReceta],
            [CantidadDespachada],
            [idDosisRecetada],
            [idEstadoDetalle],
            [MotivoAnulacionMedico],
            [observaciones],
            [IdViaAdministracion],
            [FechaProgramada],
            [CodigoDiagnostico],
            [Justificacion],
            [idUNIDDosisReceta],
            [idFrecuencia],
            [DescripcionadicionalReceta],
            [Duracion],
            [Aprobado],
            [CantidadAprobada],
            [IdDiagnostico],
            [idCuentaAtencionProxCita]
        )
        VALUES (
            @idReceta,
            @IdProducto,
            @Cantidad,
            @Precio,
            @Total,
            @SaldoEnRegistroReceta,
            NULL,
            NULL,
            @idDosisRecetada,
            1,
            NULL,
            @observaciones,
            @IdViaAdministracion,
            DATEADD(DAY, 7, GETDATE()),
            @xCodigoDiagnostico,
            @Justificacion,
            @idUNIDDosisReceta,
            @idFrecuencia,
            @DescripcionadicionalReceta,
            @Duracion,
            0,
            0,
            @xIdDiagnostico,
            @idCuentaAtencionProxCita
        );

        SET @Mensaje = 'OK;' + CAST(@idReceta AS VARCHAR);
    END
    ELSE
    BEGIN
        SET @Mensaje = 'El producto Ingresado ya Existe dentro de la lista ingresada';
    END
END;
GO
