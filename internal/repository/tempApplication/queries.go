package tempApplicationRepository

const (
	qGetTempApplication        = `SELECT * FROM temp_applications WHERE temp_application_id = :tempApplicationId`
	qGetTempApplicationByEmail = `SELECT * FROM temp_applications WHERE email = :email`
	qGetAllTempApplications    = `SELECT * FROM temp_applications`
	qInsertTempApplication     = `INSERT INTO temp_applications (temp_application_id, first_name, last_name, email, phone_number, projects, application_text, created_at) VALUES (:temp_application_id, :first_name, :last_name, :email, :phone_number, :projects, :application_text, :created_at)`
	qDeleteTempApplication     = `DELETE FROM temp_applications WHERE temp_application_id = :tempApplicationId`
)
