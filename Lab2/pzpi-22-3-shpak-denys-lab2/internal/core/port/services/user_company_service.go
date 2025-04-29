package services // import "wayra/internal/core/port/services"

import (
	"wayra/internal/core/domain/models"
)

// UserCompanyService is an interface to define the methods that the UserCompanyService should implement
type UserCompanyService interface {
	Service[models.UserCompany]
	UserBelongsToCompany(userID, companyID uint) bool
}
