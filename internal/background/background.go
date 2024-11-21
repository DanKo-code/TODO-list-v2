package background

import (
	"context"
	"time"
)

type TaskChecker interface {
	StartOverdueStatusChecker(ctx context.Context, interval time.Duration, stopChan <-chan struct{})
}
