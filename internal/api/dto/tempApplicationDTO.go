package dto

type CreateTempApplicationRequest struct {
	FirstName       string
	LastName        string
	Email           string
	PhoneNumber     string
	Projects        []string
	ApplicationText string
}
