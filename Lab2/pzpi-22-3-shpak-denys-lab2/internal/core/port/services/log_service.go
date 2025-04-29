package services

import "wayra/internal/core/domain/models"

type LogService interface {
	Service[models.Log]
	LogAction(userID *uint, actionType, description string, success bool)
}
