package service // import "wayra/internal/core/service"

import (
	"context"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port"
	"wayra/internal/core/port/services"
)

// CompanyService is a struct that defines the CompanyService
type CompanyService struct {
	*GenericService[models.Company] // Embedding the GenericService
}

// NewCompanyService is a function that returns a new CompanyService
// port: port.Repository[models.Company] - The repository that will be used by the service
// returns: *CompanyService - The service that will be used to interact with the repository
func NewCompanyService(repo port.Repository[models.Company]) services.CompanyService {
	return &CompanyService{
		GenericService: NewGenericService(repo),
	}
}

// GetAll is a method that returns all the companies in the repository
// returns: ([]models.Company, error) - The list of companies and an error if any
func (cs *CompanyService) GetAll() ([]models.Company, error) {
	companies, err := cs.Repository.GetAll(context.Background())
	if err != nil {
		return nil, err
	}
	return companies, nil
}
