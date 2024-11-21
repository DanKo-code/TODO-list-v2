package rest

import (
	"context"
	"net/http"
	"strings"
)

type Router struct {
	routes map[string]map[string]http.HandlerFunc
}

func NewRouter(handlers *Handlers) *Router {
	router := &Router{
		routes: make(map[string]map[string]http.HandlerFunc),
	}

	router.addRoute("POST", "/tasks", handlers.CreateTask)
	router.addRoute("GET", "/tasks", handlers.GetTasks)
	router.addRoute("PUT", "/tasks/{id}", handlers.UpdateTask)

	return router
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	for routPath, methods := range r.routes {
		params, match := matchRoute(routPath, req.URL.Path)

		if match {
			if handler, ok := methods[req.Method]; ok {
				ctx := req.Context()

				for key, value := range params {
					ctx = context.WithValue(ctx, key, value)
				}
				handler(w, req.WithContext(ctx))
				return
			}

			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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

func matchRoute(routePath, requestPath string) (map[string]string, bool) {
	routeParts := strings.Split(routePath, "/")
	requestParts := strings.Split(requestPath, "/")

	if len(routeParts) != len(requestParts) {
		return nil, false
	}

	params := make(map[string]string)

	for i, routePart := range routeParts {
		if strings.HasPrefix(routePart, "{") && strings.HasSuffix(routePart, "}") {
			paramName := strings.Trim(routePart, "{}")
			params[paramName] = requestParts[i]
		} else if routePart != requestParts[i] {
			return nil, false
		}
	}

	return params, true
}
