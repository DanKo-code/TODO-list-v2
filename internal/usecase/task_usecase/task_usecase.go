package task_usecase

import (
	"context"
	"github.com/DanKo-code/TODO-list/internal/dtos"
	"github.com/DanKo-code/TODO-list/internal/models"
	"github.com/DanKo-code/TODO-list/internal/repository"
	"github.com/DanKo-code/TODO-list/pkg/helper"
	"time"
)

type TaskUseCase struct {
	taskRep repository.TaskRepository
}

func NewTaskUseCase(taskRep repository.TaskRepository) *TaskUseCase {
	return &TaskUseCase{
		taskRep: taskRep,
	}
}

func (tuc *TaskUseCase) CreateTask(ctx context.Context, cmd *dtos.CreateTaskCommand) (*models.Task, error) {

	taskId, _ := helper.GenerateUUID()

	if cmd.DueDate == "" {
		cmd.DueDate = time.Now().Add(24 * time.Hour).Format("2006-01-02")
	}

	task := &models.Task{
		Id:          taskId,
		Title:       cmd.Title,
		Description: cmd.Description,
		DueDate:     cmd.DueDate,
		Overdue:     false,
		Completed:   false,
	}

	err := tuc.taskRep.Save(ctx, task)
	if err != nil {
		return nil, err
	}

	return task, nil
}
