package utils

import (
	"net/http"
)

// Router wraps the default ServeMux to add method-based routing.
type Router struct {
	mux *http.ServeMux
}

// NewRouter creates a new Router instance.
func NewRouter() *Router {
	return &Router{
		mux: http.NewServeMux(),
	}
}

// Handle registers a handler for a specific method and path.
func (r *Router) Handle(method string, path string, handler http.HandlerFunc) {
	r.mux.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, req)
	})
}

// ServeHTTP makes Router implement http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
