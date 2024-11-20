package server

import (
	"context"
	"errors"
	"fmt"
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

type App struct {
	server *http.Server
	tRep   repository.TaskRepository
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

	return &App{
		server: server,
		tRep:   tRep,
	}, nil
}

func (a *App) Run() error {

	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.FatalLogger.Fatalf("Failed to listen and serve: %v", err)
		}
	}()

	logger.InfoLogger.Println(fmt.Sprintf("Server started on address %s", a.server.Addr))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit
	logger.InfoLogger.Println("Server stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		a.tRep.Close()
		cancel()
	}()

	return a.server.Shutdown(ctx)
}
