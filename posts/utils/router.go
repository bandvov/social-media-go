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

// GET adds a handler specifically for GET requests.
func (r *Router) GET(pattern string, handler func(w http.ResponseWriter, req *http.Request)) {
	r.mux.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, req)
	})
}

// POST adds a handler specifically for POST requests.
func (r *Router) POST(pattern string, handler func(w http.ResponseWriter, req *http.Request)) {
	r.mux.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, req)
	})
}

// POST adds a handler specifically for PUT requests.
func (r *Router) PUT(pattern string, handler func(w http.ResponseWriter, req *http.Request)) {
	r.mux.HandleFunc(pattern, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodPut {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}
		handler(w, req)
	})
}

// ServeHTTP makes Router implement http.Handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}
