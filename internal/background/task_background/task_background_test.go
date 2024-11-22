package task_background

import (
	"context"
	"errors"
	"github.com/DanKo-code/TODO-list/internal/usecase/task_usecase"
	"testing"
	"time"
)

func TestStartOverdueStatusChecker(t *testing.T) {
	t.Run("successful update and stop", func(t *testing.T) {
		mockUseCase := &task_usecase.MockTaskUseCase{
			UpdateOverdueTasksFunc: func(ctx context.Context) error {
				return nil
			},
			Called: false,
		}

		tc := NewTaskChecker(mockUseCase)
		ctx := context.Background()

		stopChan := make(chan struct{})
		interval := 10 * time.Millisecond

		go func() {
			time.Sleep(50 * time.Millisecond)
			close(stopChan)
		}()

		tc.StartOverdueStatusChecker(ctx, interval, stopChan)

		if !mockUseCase.Called {
			t.Error("Expected UpdateOverdueTasks to be called")
		}
	})

	t.Run("error during update", func(t *testing.T) {
		mockUseCase := &task_usecase.MockTaskUseCase{
			UpdateOverdueTasksFunc: func(ctx context.Context) error {
				return errors.New("update error")
			},
			Called: false,
		}

		tc := NewTaskChecker(mockUseCase)
		ctx := context.Background()

		stopChan := make(chan struct{})
		interval := 10 * time.Millisecond

		go func() {
			time.Sleep(50 * time.Millisecond)
			close(stopChan)
		}()

		tc.StartOverdueStatusChecker(ctx, interval, stopChan)

		if mockUseCase.UpdateOverdueTasksFunc == nil {
			t.Error("Expected UpdateOverdueTasks to be called")
		}
	})
}
