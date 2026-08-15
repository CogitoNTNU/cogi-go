package dto

// NOTE: keep validation in sync with what BOTH frontends actually send.
// The old site's Verv form legitimately submits an empty ApplicationText,
// and neither frontend guarantees RFC-valid email format — so only the
// identity fields are required, and email format is not enforced here.
type CreateTempApplicationRequest struct {
	FirstName       string   `binding:"required"`
	LastName        string   `binding:"required"`
	Email           string   `binding:"required"`
	PhoneNumber     string   `binding:"required"`
	Projects        []string `binding:"required,min=1"`
	ApplicationText string
}
