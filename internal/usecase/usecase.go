package usecase

import (
	"context"
	"github.com/DanKo-code/TODO-list/internal/dtos"
	"github.com/DanKo-code/TODO-list/internal/models"
)

type TaskUseCase interface {
	CreateTask(ctx context.Context, cmd *dtos.CreateTaskCommand) (*models.Task, error)
}
