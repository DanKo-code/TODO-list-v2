package server

import "net/http"

type App struct {
	server *http.Server
	taskUC
}
