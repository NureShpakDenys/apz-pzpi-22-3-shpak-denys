package service // import "wayra/internal/core/service"

import (
	"context"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port"
	"wayra/internal/core/port/services"

	"gorm.io/gorm"
)

// UserCompanyService is a struct to manage the user company service
type UserCompanyService struct {
	*GenericService[models.UserCompany] // Embedding the generic service
}

// NewUserCompanyService creates a new user company service
// repo: Repository of the user company
// returns: A new user company service
func NewUserCompanyService(repo port.Repository[models.UserCompany]) services.UserCompanyService {
	return &UserCompanyService{
		GenericService: NewGenericService(repo),
	}
}

// UserBelongsToCompany checks if a user belongs to a company
// userID: ID of the user
// companyID: ID of the company
// returns: True if the user belongs to the company, false otherwise
func (s *UserCompanyService) UserBelongsToCompany(userID, companyID uint) bool {
	userCompany, err := s.getByUserAndCompany(context.Background(), userID, companyID)
	if err != nil || userCompany == nil {
		return false
	}
	return true
}

// getByUserAndCompany gets a user company by user and company
// ctx: Context of the request
// userID: ID of the user
// companyID: ID of the company
// returns: The user company, or an error if it does not exist
func (s *UserCompanyService) getByUserAndCompany(ctx context.Context, userID, companyID uint) (*models.UserCompany, error) {
	userCompanies, err := s.Where(ctx, &models.UserCompany{
		UserID:    userID,
		CompanyID: companyID,
	})
	if err != nil {
		return nil, err
	}

	if len(userCompanies) == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &userCompanies[0], nil
}
