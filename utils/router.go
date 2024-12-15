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
func (r *Router) HandleFunc(pattern string, handler func(w http.ResponseWriter, req *http.Request)) {
	r.mux.HandleFunc(pattern, handler)
}

// ServeHTTP makes Router implement http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
