package sqlite

import (
	"context"
	"database/sql"
	"github.com/DanKo-code/TODO-list/internal/models"
	"github.com/DanKo-code/TODO-list/pkg/logger"
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
