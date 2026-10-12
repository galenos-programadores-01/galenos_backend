CREATE OR ALTER PROCEDURE dbo.usp_go_Appointments_ObtenerPorId
    @ID VARCHAR(50)
AS
BEGIN
    SET NOCOUNT ON;

    SELECT 
        id, 
        patient_id, 
        doctor_id, 
        starts_at, 
        ends_at, 
        status, 
        reason, 
        created_at, 
        updated_at
    FROM dbo.appointments WITH (NOLOCK)
    WHERE id = @ID;
END;
GO
