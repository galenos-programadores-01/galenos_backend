CREATE OR ALTER PROCEDURE dbo.usp_go_AtencionSintomas_Insertar
    @IdRegAtencion     INT,
    @Sistema           VARCHAR(100),
    @Sintoma           VARCHAR(250),
    @IdUsuarioAuditoria INT = 1
AS
BEGIN
    SET NOCOUNT ON;

    INSERT INTO dbo.Tab_Atencion_Sintomas
    (
        IdRegAtencion,
        Sistema,
        Sintoma,
        EsPredefinido,
        SG_LOG_CRE_FECH,
        SG_LOG_CRE_USUARIO
    )
    VALUES
    (
        @IdRegAtencion,
        ISNULL(NULLIF(RTRIM(LTRIM(@Sistema)), ''), 'General'),
        RTRIM(LTRIM(@Sintoma)),
        1,
        GETDATE(),
        ISNULL(@IdUsuarioAuditoria, 1)
    );
END;
GO
