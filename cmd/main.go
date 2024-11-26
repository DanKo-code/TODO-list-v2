package main

import (
	"github.com/DanKo-code/TODO-list/internal/server"
	"github.com/DanKo-code/TODO-list/pkg/logger"
	"os"
)

func main() {

	//Test workflows

	app, err := server.NewApp(os.Getenv("APP_ADDRESS"), os.Getenv("DB_DRIVER"), os.Getenv("DB_NAME"))
	if err != nil {
		logger.FatalLogger.Fatal("Failed to initialize app")
	}

	if err := app.Run(); err != nil {
		logger.FatalLogger.Fatalf("Server Shutdown Failed:%+v", err)
	}

	logger.InfoLogger.Println("Server Shutdown Success")
}
