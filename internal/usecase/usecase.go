package usecase

import (
	"context"
	"github.com/DanKo-code/TODO-list/internal/dtos"
	"github.com/DanKo-code/TODO-list/internal/models"
)

type TaskUseCase interface {
	CreateTask(ctx context.Context, cmd *dtos.CreateTaskCommand) (*models.Task, error)
	GetTasks(ctx context.Context) ([]*models.Task, error)
	UpdateTask(ctx context.Context, id string, updateTaskCommand *dtos.UpdateTaskCommand) (*models.Task, error)
}
