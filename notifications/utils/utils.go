package utils

import (
	"net/http"
	"strconv"
)

type contextKey string

var UserIDKey contextKey = "userID"

// parsePagination extracts limit and offset from query parameters with defaults
func ParsePagination(r *http.Request) (int, int) {
	query := r.URL.Query()
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page < 0 {
		page = 1
	}
	return limit, page
}
