IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'appointments')
BEGIN
    CREATE TABLE dbo.appointments (
        id VARCHAR(50) PRIMARY KEY,
        patient_id VARCHAR(50) NOT NULL,
        doctor_id VARCHAR(50) NOT NULL,
        starts_at DATETIME2 NOT NULL,
        ends_at DATETIME2 NOT NULL,
        status VARCHAR(50) NOT NULL,
        reason NVARCHAR(500) NULL,
        created_at DATETIME2 NOT NULL,
        updated_at DATETIME2 NOT NULL
    );
END;
GO

CREATE OR ALTER PROCEDURE dbo.usp_go_Appointments_Crear
    @ID VARCHAR(50),
    @PatientID VARCHAR(50),
    @DoctorID VARCHAR(50),
    @StartsAt DATETIME2,
    @EndsAt DATETIME2,
    @Status VARCHAR(50),
    @Reason NVARCHAR(500),
    @CreatedAt DATETIME2,
    @UpdatedAt DATETIME2
AS
BEGIN
    SET NOCOUNT ON;

    -- Validar conflicto de horario
    IF EXISTS (
        SELECT 1
        FROM dbo.appointments WITH (UPDLOCK, HOLDLOCK)
        WHERE doctor_id = @DoctorID
          AND status <> 'CANCELLED'
          AND starts_at < @EndsAt
          AND ends_at > @StartsAt
    )
    BEGIN
        SELECT -1 AS CodigoResultado, 'El médico no se encuentra disponible en el horario seleccionado' AS Mensaje;
        RETURN;
    END;

    INSERT INTO dbo.appointments
        (id, patient_id, doctor_id, starts_at, ends_at, status, reason, created_at, updated_at)
    VALUES
        (@ID, @PatientID, @DoctorID, @StartsAt, @EndsAt, @Status, @Reason, @CreatedAt, @UpdatedAt);

    SELECT 1 AS CodigoResultado, 'Cita registrada exitosamente' AS Mensaje;
END;
GO
