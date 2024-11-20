package repository

import (
	"context"
	"github.com/DanKo-code/TODO-list/internal/models"
)

type TaskRepository interface {
	Close()
	Save(ctx context.Context, task *models.Task) error
	/*GetAll(ctx context.Context) ([]*models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	DeleteById(ctx context.Context, id string) (*models.Task, error)
	MarkAsCompleted(ctx context.Context, id string) (*models.Task, error)*/

}
