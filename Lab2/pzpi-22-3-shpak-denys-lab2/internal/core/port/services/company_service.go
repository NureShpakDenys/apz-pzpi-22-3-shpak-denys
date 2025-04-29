package services // import "wayra/internal/core/port/services"

import (
	"wayra/internal/core/domain/models"
)

// CompanyService is an interface to represent the company service
type CompanyService interface {
	Service[models.Company]
	GetAll() ([]models.Company, error)
}
