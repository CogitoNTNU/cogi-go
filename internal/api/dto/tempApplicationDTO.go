package dto

type CreateTempApplicationRequest struct {
	FirstName       string   `binding:"required"`
	LastName        string   `binding:"required"`
	Email           string   `binding:"required,email"`
	PhoneNumber     string   `binding:"required"`
	Projects        []string `binding:"required,min=1"`
	ApplicationText string   `binding:"required"`
}
