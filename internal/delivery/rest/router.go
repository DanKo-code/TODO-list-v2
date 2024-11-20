package rest

import (
	"net/http"
)

type Router struct {
	routes map[string]map[string]http.HandlerFunc
}

func NewRouter(handlers *Handlers) *Router {
	router := &Router{
		routes: make(map[string]map[string]http.HandlerFunc),
	}

	router.addRoute("POST", "/tasks", handlers.CreateTask)

	return router
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if handler, ok := r.routes[req.URL.Path]; ok {
		if handler, ok := handler[req.Method]; ok {
			handler(w, req)
			return
		}
	}
	http.NotFound(w, req)
}

func (r *Router) addRoute(method, path string, handler http.HandlerFunc) {
	if r.routes[path] == nil {
		r.routes[path] = make(map[string]http.HandlerFunc)
	}
	r.routes[path][method] = handler
}
