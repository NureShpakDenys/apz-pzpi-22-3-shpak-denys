package service

import (
	"context"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port"
	"wayra/internal/core/port/services"
)

type LogService struct {
	*GenericService[models.Log]
}

func NewLogService(repo port.Repository[models.Log]) services.LogService {
	return &LogService{
		GenericService: NewGenericService(repo),
	}
}

func (ls *LogService) LogAction(userID *uint, actionType, description string, success bool) {
	ls.Repository.Add(context.Background(), &models.Log{
		UserID:      userID,
		ActionType:  actionType,
		Description: description,
		Success:     success,
	})
}
