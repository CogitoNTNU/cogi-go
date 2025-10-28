package userRepository

const (
	qGetUser      = `SELECT * FROM USERS WHERE user_id = :userId`
	qGetUsers     = `SELECT * FROM USERS`
	qInsertUser   = `INSERT INTO USERS (user_id, first_name, last_name, description, nickname, email, phone, gender, github_url, linkedin_url, kaggle_url, huggingface_url, password, avatar, image_permission, food_preference) VALUES (:userId, :firstName, :lastName, :description, :nickname, :email, :phone, :gender, :github_url, :linkedinUrl, :kaggleUrl, :huggingfaceUrl, :password, :avatar, :imagePermission, :foodPreference)`
	qDeleteUser   = `DELETE FROM USERS WHERE user_id = :userId`
	qGetUserByEmail = `SELECT * FROM USERS WHERE email = :email`
)

