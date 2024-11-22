package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/DanKo-code/TODO-list/internal/dtos"
	internalErrors "github.com/DanKo-code/TODO-list/internal/errors"
	"github.com/DanKo-code/TODO-list/internal/models"
	"github.com/DanKo-code/TODO-list/pkg/logger"
	"strings"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(driver string, dsn string) (*TaskRepository, error) {

	db, err := sql.Open(driver, dsn)
	if err != nil {
		logger.ErrorLogger.Printf("Failed to connect database: %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.ErrorLogger.Printf("Failed to verify database: %v", err)
		return nil, err
	}

	logger.InfoLogger.Println("Database connected")

	return &TaskRepository{db: db}, nil
}

func (s *TaskRepository) Init(ctx context.Context) error {
	q := `CREATE TABLE IF NOT EXISTS tasks 
			(id uuid, title TEXT, description TEXT, due_date TEXT, overdue INTEGER, completed INTEGER)`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		logger.ErrorLogger.Println("Failed to init db: %v", err)
		return err
	}

	logger.InfoLogger.Println("Database initialized")

	return nil
}

func (s *TaskRepository) Close() {
	if err := s.db.Close(); err != nil {
		logger.ErrorLogger.Println("Failed to close db connection: %v", err)
	} else {
		logger.InfoLogger.Println("Database connection closed")
	}
}

func (s *TaskRepository) Save(ctx context.Context, task *models.Task) error {
	q := `INSERT INTO tasks (id, title, description, due_date, overdue, completed)
			VALUES ($1, $2, $3, $4, $5, $6);`

	_, err := s.db.ExecContext(ctx, q,
		task.Id,
		task.Title,
		task.Description,
		task.DueDate,
		task.Overdue,
		task.Completed,
	)
	if err != nil {
		logger.ErrorLogger.Printf("failed to save task: %v", err)
		return err
	}
	return nil
}

func (s *TaskRepository) GetAll(ctx context.Context) ([]*models.Task, error) {
	q := `SELECT id, title, description, due_date, overdue, completed FROM tasks`

	rows, err := s.db.QueryContext(ctx, q)
	if err != nil {
		logger.ErrorLogger.Printf("failed to fetch tasks: %v", err)
	}
	defer rows.Close()

	var tasks []*models.Task

	for rows.Next() {
		task := &models.Task{}

		err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.DueDate,
			&task.Overdue,
			&task.Completed,
		)
		if err != nil {
			logger.ErrorLogger.Printf("failed to scan task: %v", err)
			return nil, err
		}

		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		logger.ErrorLogger.Printf("rows iteration error: %v", err)
		return nil, err
	}

	return tasks, nil
}

func (s *TaskRepository) GetById(ctx context.Context, id string) (*models.Task, error) {
	q := `SELECT id, title, description, due_date, overdue, completed
		  FROM tasks
		  WHERE id = $1`

	task := &models.Task{}
	row := s.db.QueryRowContext(ctx, q, id)

	err := row.Scan(
		&task.Id,
		&task.Title,
		&task.Description,
		&task.DueDate,
		&task.Overdue,
		&task.Completed,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.ErrorLogger.Printf("task with id %s not found", id)
			return nil, internalErrors.TaskNotFound
		}

		logger.ErrorLogger.Printf("failed to fetch task: %v", err)
		return nil, err
	}

	return task, nil
}

func (s *TaskRepository) Update(ctx context.Context, id string, updateTaskCommand *dtos.UpdateTaskCommand) error {
	q := `UPDATE tasks SET `
	var args []interface{}
	var setClauses []string

	if updateTaskCommand.Title != "" {
		setClauses = append(setClauses, "title = ?")
		args = append(args, updateTaskCommand.Title)
	}
	if updateTaskCommand.Description != "" {
		setClauses = append(setClauses, "description = ?")
		args = append(args, updateTaskCommand.Description)
	}
	if updateTaskCommand.DueDate != "" {
		setClauses = append(setClauses, "due_date = ?")
		args = append(args, updateTaskCommand.DueDate)

		setClauses = append(setClauses, "overdue = ?")
		args = append(args, false)
	}

	if len(setClauses) == 0 {
		logger.ErrorLogger.Println("no fields to update")
		return fmt.Errorf("no fields to update")
	}

	q += strings.Join(setClauses, ", ")
	q += " WHERE id = ?"
	args = append(args, id)

	_, err := s.db.ExecContext(ctx, q, args...)
	if err != nil {
		logger.ErrorLogger.Printf("failed to update task: %v", err)
		return err
	}

	return nil
}

func (s *TaskRepository) DeleteById(ctx context.Context, id string) error {
	q := `DELETE FROM tasks WHERE id = $1`

	_, err := s.db.ExecContext(ctx, q, id)
	if err != nil {
		logger.ErrorLogger.Printf("failed to delete task: %v", err)
		return err
	}

	return nil
}

func (s *TaskRepository) ChangeCompletionStatus(ctx context.Context, id string, completionStatus bool) error {
	q := `UPDATE tasks SET completed = $1 WHERE id = $2`

	_, err := s.db.ExecContext(ctx, q, completionStatus, id)
	if err != nil {
		logger.ErrorLogger.Printf("failed to change completion status: %v", err)
		return err
	}

	return nil
}

func (s *TaskRepository) UpdateOverdueTasks(ctx context.Context) error {
	q := `UPDATE tasks 
		  SET overdue = TRUE 
		  WHERE due_date <= DATE('now') AND overdue = FALSE`

	_, err := s.db.ExecContext(ctx, q)
	if err != nil {
		return err
	}

	return nil
}
