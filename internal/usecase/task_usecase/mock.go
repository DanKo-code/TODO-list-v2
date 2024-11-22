package task_usecase

import (
	"context"
	"github.com/DanKo-code/TODO-list/internal/dtos"
	"github.com/DanKo-code/TODO-list/internal/models"
)

type MockTaskUseCase struct {
	CreateTaskFunc                 func(ctx context.Context, cmd *dtos.CreateTaskCommand) (*models.Task, error)
	GetTaskFunc                    func(ctx context.Context) ([]*models.Task, error)
	UpdateTaskFunc                 func(ctx context.Context, id string, updateTaskCommand *dtos.UpdateTaskCommand) (*models.Task, error)
	DeleteTaskFunc                 func(ctx context.Context, id string) error
	ChangeTaskCompletionStatusFunc func(ctx context.Context, id string, completionStatus bool) (*models.Task, error)
	UpdateOverdueTasksFunc         func(ctx context.Context) error
	Called                         bool
}

func (m *MockTaskUseCase) CreateTask(ctx context.Context, cmd *dtos.CreateTaskCommand) (*models.Task, error) {
	return m.CreateTaskFunc(ctx, cmd)
}

func (m *MockTaskUseCase) GetTasks(ctx context.Context) ([]*models.Task, error) {
	return m.GetTaskFunc(ctx)
}

func (m *MockTaskUseCase) UpdateTask(ctx context.Context, id string, updateTaskCommand *dtos.UpdateTaskCommand) (*models.Task, error) {
	return m.UpdateTaskFunc(ctx, id, updateTaskCommand)
}

func (m *MockTaskUseCase) DeleteTask(ctx context.Context, id string) error {
	return m.DeleteTaskFunc(ctx, id)
}

func (m *MockTaskUseCase) ChangeTaskCompletionStatus(ctx context.Context, id string, completionStatus bool) (*models.Task, error) {
	return m.ChangeTaskCompletionStatusFunc(ctx, id, completionStatus)
}

func (m *MockTaskUseCase) UpdateOverdueTasks(ctx context.Context) error {
	m.Called = true
	return m.UpdateOverdueTasksFunc(ctx)
}
