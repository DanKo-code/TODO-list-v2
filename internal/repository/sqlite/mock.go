package sqlite

import (
	"context"
	"github.com/DanKo-code/TODO-list/internal/dtos"
	"github.com/DanKo-code/TODO-list/internal/models"
)

type MockTaskRepository struct {
	CloseFunc                  func()
	SaveFunc                   func(ctx context.Context, task *models.Task) error
	GetAllFunc                 func(ctx context.Context) ([]*models.Task, error)
	GetByIdFunc                func(ctx context.Context, id string) (*models.Task, error)
	UpdateFunc                 func(ctx context.Context, id string, updateTaskCommand *dtos.UpdateTaskCommand) error
	DeleteByIdFunc             func(ctx context.Context, id string) error
	ChangeCompletionStatusFunc func(ctx context.Context, id string, completionStatus bool) error
	UpdateOverdueTasksFunc     func(ctx context.Context) error
}

func (m MockTaskRepository) Close() {
}

func (m MockTaskRepository) Save(ctx context.Context, task *models.Task) error {
	return m.SaveFunc(ctx, task)
}

func (m MockTaskRepository) GetAll(ctx context.Context) ([]*models.Task, error) {
	return m.GetAllFunc(ctx)
}

func (m MockTaskRepository) GetById(ctx context.Context, id string) (*models.Task, error) {
	return m.GetByIdFunc(ctx, id)
}

func (m MockTaskRepository) Update(ctx context.Context, id string, updateTaskCommand *dtos.UpdateTaskCommand) error {
	return m.UpdateFunc(ctx, id, updateTaskCommand)
}

func (m MockTaskRepository) DeleteById(ctx context.Context, id string) error {
	return m.DeleteByIdFunc(ctx, id)
}

func (m MockTaskRepository) ChangeCompletionStatus(ctx context.Context, id string, completionStatus bool) error {
	return m.ChangeCompletionStatusFunc(ctx, id, completionStatus)
}

func (m MockTaskRepository) UpdateOverdueTasks(ctx context.Context) error {
	return m.UpdateOverdueTasksFunc(ctx)
}
