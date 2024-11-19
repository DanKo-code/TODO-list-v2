package cmd

import "os"

const (
	ServerAddressEnv = "SERVER_ADDRESS"
)

func main() {
	app := server.NewApp()

	if err := app.Run(os.Getenv("APP_PORT")); err != nil {
		logrusCustom.Logger.Fatalf("Error when running server: %s", err.Error())

	}
}
