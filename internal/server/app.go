package server

import (
	"context"
	"errors"
	"github.com/DanKo-code/TODO-list/internal/background"
	"github.com/DanKo-code/TODO-list/internal/background/task_background"
	"github.com/DanKo-code/TODO-list/internal/delivery/rest"
	"github.com/DanKo-code/TODO-list/internal/repository"
	sqliteRep "github.com/DanKo-code/TODO-list/internal/repository/sqlite"
	"github.com/DanKo-code/TODO-list/internal/usecase/task_usecase"
	"github.com/DanKo-code/TODO-list/pkg/logger"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	interval = 1 * time.Minute / 3
)

type App struct {
	server *http.Server
	tRep   repository.TaskRepository
	tc     background.TaskChecker
}

func NewApp(appAddress, driver, dsn string) (*App, error) {
	tRep, err := sqliteRep.NewTaskRepository(driver, dsn)
	if err != nil {
		return nil, err
	}

	err = tRep.Init(context.TODO())
	if err != nil {
		return nil, err
	}

	taskUseCase := task_usecase.NewTaskUseCase(tRep)

	handlers := rest.NewHandlers(taskUseCase)

	router := rest.NewRouter(handlers)

	server := &http.Server{
		Addr:    appAddress,
		Handler: router,
	}

	tc := task_background.NewTaskChecker(taskUseCase)

	return &App{
		server: server,
		tRep:   tRep,
		tc:     tc,
	}, nil
}

func (a *App) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		a.tRep.Close()
		cancel()
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	stopChecker := make(chan struct{})

	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.FatalLogger.Fatalf("Failed to listen and serve: %v", err)
		}
	}()

	logger.InfoLogger.Printf("Server started on address %s", a.server.Addr)

	go a.tc.StartOverdueStatusChecker(context.TODO(), interval, stopChecker)

	<-quit

	close(stopChecker)

	return a.server.Shutdown(ctx)
}
