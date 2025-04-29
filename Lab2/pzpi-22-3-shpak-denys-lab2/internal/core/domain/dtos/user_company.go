package dtos // import "github.com/zatte/iam/internal/core/domain/dtos"

// UserCompanyDTO is a DTO for UserCompany
type UserCompanyDTO struct {
	// Id is the unique identifier of the UserCompany
	// Example: 1
	ID uint `json:"id"`

	// UserID is the unique identifier of the User
	// Example: 1
	UserID uint `json:"user_id"`

	// CompanyID is the unique identifier of the Company
	// Example: 1
	CompanyID uint `json:"company_id"`

	// Role is the role of the User in the Company
	// Example: admin
	Role string `json:"role"`

	// User is the User of the UserCompany
	User UserDTO `json:"user,omitempty"`

	// Company is the Company of the UserCompany
	Company CompanyDTO `json:"company,omitempty"`
}
