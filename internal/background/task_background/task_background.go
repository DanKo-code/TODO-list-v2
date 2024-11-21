package task_background

import (
	"context"
	"github.com/DanKo-code/TODO-list/internal/usecase"
	"github.com/DanKo-code/TODO-list/pkg/logger"
	"time"
)

type TaskChecker struct {
	usecase usecase.TaskUseCase
}

func NewTaskChecker(useCase usecase.TaskUseCase) *TaskChecker {
	return &TaskChecker{
		usecase: useCase,
	}
}

func (tc *TaskChecker) StartOverdueStatusChecker(ctx context.Context, interval time.Duration, stopChan <-chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := tc.usecase.UpdateOverdueTasks(ctx); err != nil {
				logger.ErrorLogger.Printf("Error updating overdue tasks: %v", err)
				return
			}
			logger.InfoLogger.Println("Updated overdue tasks")
		case <-stopChan:
			logger.InfoLogger.Println("Stopping overdue task checker")
			return
		}
	}
}
