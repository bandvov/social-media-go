package utils

import (
	"net/http"
	"strconv"
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

// parsePagination extracts limit and offset from query parameters with defaults
func ParsePagination(r *http.Request) (int, int) {
	query := r.URL.Query()
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	offset, err := strconv.Atoi(query.Get("offset"))
	if err != nil || offset < 0 {
		offset = 0
	}
	return limit, offset
}
