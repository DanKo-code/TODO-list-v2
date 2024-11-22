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

func (tuc *TaskUseCase) GetTasks(ctx context.Context) ([]*models.Task, error) {

	tasks, err := tuc.taskRep.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (tuc *TaskUseCase) UpdateTask(ctx context.Context, id string, updateTaskCommand *dtos.UpdateTaskCommand) (*models.Task, error) {

	task, err := tuc.taskRep.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	err = tuc.taskRep.Update(ctx, id, updateTaskCommand)
	if err != nil {
		return nil, err
	}

	updatedTask := createUpdateTaskRes(task, updateTaskCommand)

	return updatedTask, nil
}

func createUpdateTaskRes(task *models.Task, updateTaskCommand *dtos.UpdateTaskCommand) *models.Task {
	updatedTask := &models.Task{
		Id:        task.Id,
		Overdue:   task.Overdue,
		Completed: task.Completed,
	}
	if updateTaskCommand.Title == "" {
		updatedTask.Title = task.Title
	} else {
		updatedTask.Title = updateTaskCommand.Title
	}

	if updateTaskCommand.Description == "" {
		updatedTask.Description = task.Description
	} else {
		updatedTask.Description = updateTaskCommand.Description
	}

	if updateTaskCommand.DueDate == "" {
		updatedTask.DueDate = task.DueDate
	} else {
		updatedTask.DueDate = updateTaskCommand.DueDate
		updatedTask.Overdue = false
	}

	return updatedTask
}

func (tuc *TaskUseCase) DeleteTask(ctx context.Context, id string) error {

	_, err := tuc.taskRep.GetById(ctx, id)
	if err != nil {
		return err
	}

	err = tuc.taskRep.DeleteById(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (tuc *TaskUseCase) ChangeTaskCompletionStatus(ctx context.Context, id string, completionStatus bool) (*models.Task, error) {
	task, err := tuc.taskRep.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	err = tuc.taskRep.ChangeCompletionStatus(ctx, id, completionStatus)
	if err != nil {
		return nil, err
	}

	task.Completed = completionStatus

	return task, nil
}

func (tuc *TaskUseCase) UpdateOverdueTasks(ctx context.Context) error {
	err := tuc.taskRep.UpdateOverdueTasks(ctx)
	if err != nil {
		return err
	}

	return nil
}
