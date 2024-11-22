package task_usecase

import (
	"context"
	"github.com/DanKo-code/TODO-list/internal/dtos"
	"github.com/DanKo-code/TODO-list/internal/models"
	"github.com/DanKo-code/TODO-list/internal/repository/sqlite"
	"testing"
	"time"
)

func TestCreateTaskUseCase(t *testing.T) {
	test := []struct {
		name         string
		param        dtos.CreateTaskCommand
		mockSaveFunc func(ctx context.Context, task *models.Task) error
		result       models.Task
	}{
		{
			name: "success",
			param: dtos.CreateTaskCommand{
				Title:       "Test Task",
				Description: "Test Description",
				DueDate:     "2024-11-22",
			},
			mockSaveFunc: func(ctx context.Context, task *models.Task) error {
				return nil
			},
			result: models.Task{
				Title:       "Test Task",
				Description: "Test Description",
				DueDate:     "2024-11-22",
				Overdue:     false,
				Completed:   false,
			},
		},
		{
			name: "void due date",
			param: dtos.CreateTaskCommand{
				Title:       "Test Task",
				Description: "Test Description",
			},
			mockSaveFunc: func(ctx context.Context, task *models.Task) error {
				return nil
			},
			result: models.Task{
				Title:       "Test Task",
				Description: "Test Description",
				DueDate:     time.Now().Add(24 * time.Hour).Format("2006-01-02"),
				Overdue:     false,
				Completed:   false,
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mockRepository := &sqlite.MockTaskRepository{
				SaveFunc: tt.mockSaveFunc,
			}

			ntuc := NewTaskUseCase(mockRepository)

			task, err := ntuc.CreateTask(ctx, &tt.param)
			if err != nil {
				return
			}

			if !(task.Title == tt.param.Title &&
				task.Description == tt.param.Description &&
				task.DueDate == tt.param.DueDate) {
				t.Errorf("expected the same fields from cmd: %v but on result model: %v", task, tt.result)
			}
		})
	}
}

func TestGetTasksUseCase(t *testing.T) {
	test := []struct {
		name           string
		mockGetAllFunc func(ctx context.Context) ([]*models.Task, error)
		result         []*models.Task
	}{
		{
			name: "success",
			mockGetAllFunc: func(ctx context.Context) ([]*models.Task, error) {
				return []*models.Task{
					{Id: "a495465c-d177-48e1-8954-516bba76d541", Title: "Test Task", Description: "This is a test task", DueDate: "2024-11-22", Overdue: false, Completed: false},
				}, nil
			},
			result: []*models.Task{
				{Id: "a495465c-d177-48e1-8954-516bba76d541", Title: "Test Task", Description: "This is a test task", DueDate: "2024-11-22", Overdue: false, Completed: false},
			},
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mockRepository := &sqlite.MockTaskRepository{
				GetAllFunc: tt.mockGetAllFunc,
			}

			ntuc := NewTaskUseCase(mockRepository)

			tasks, err := ntuc.GetTasks(ctx)
			if err != nil {
				return
			}

			if len(tasks) != len(tt.result) {
				t.Errorf("expected: %v but got: %v", len(tt.result), len(tasks))
			}
		})
	}
}

func TestUpdateUseCase(t *testing.T) {
	test := []struct {
		name            string
		id              string
		param           *dtos.CreateTaskCommand
		mockGetByIdFunc func(ctx context.Context, id string) (*models.Task, error)
		mockUpdate      func(ctx context.Context, id string, updateTaskCommand *dtos.UpdateTaskCommand) error
		result          *models.Task
	}{
		{
			name: "success",
			id:   "a495465c-d177-48e1-8954-516bba76d541",
			param: &dtos.CreateTaskCommand{
				Title:       "Test Task!",
				Description: "Test Description!",
			},
			mockGetByIdFunc: func(ctx context.Context, id string) (*models.Task, error) {
				return &models.Task{
					Id: "a495465c-d177-48e1-8954-516bba76d541", Title: "Test Task", Description: "This is a test task", DueDate: "2024-11-22", Overdue: false, Completed: false,
				}, nil
			},
			mockUpdate: func(ctx context.Context, id string, updateTaskCommand *dtos.UpdateTaskCommand) error {
				return nil
			},
			result: &models.Task{
				Id: "a495465c-d177-48e1-8954-516bba76d541", Title: "Test Task!", Description: "This is a test task1", DueDate: "2024-11-22", Overdue: false, Completed: false,
			},
		},
	}

	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()

			mockRepository := &sqlite.MockTaskRepository{
				GetByIdFunc: tt.mockGetByIdFunc,
				UpdateFunc:  tt.mockUpdate,
			}

			ntuc := NewTaskUseCase(mockRepository)

			task, err := ntuc.UpdateTask(ctx, tt.id, &dtos.UpdateTaskCommand{
				Title:       tt.param.Title,
				Description: tt.param.Description,
			})
			if err != nil {
				return
			}

			if !(task.Title == tt.param.Title &&
				task.Description == tt.param.Description) {
				t.Errorf("expected the same fields from cmd: %v but on result model: %v", task, tt.result)
			}
		})
	}
}
