package utils

import (
	"fmt"
	"net/http"
	"posts/domain"
	"strconv"
	"strings"
)

// Helper to generate placeholders for IN clause
func Placeholders(count int) string {
	placeholders := make([]string, count)
	for i := range placeholders {
		placeholders[i] = fmt.Sprintf("$%d", i+1)
	}
	return strings.Join(placeholders, ", ")
}

// Helper to convert int slice to interface{} slice for query arguments
func ToInterface(ids []int) []interface{} {
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	return args
}

// parsePagination extracts limit and offset from query parameters with defaults
func ParsePagination(r *http.Request) domain.Pagination {
	query := r.URL.Query()
	page, err := strconv.Atoi(query.Get("page"))
	if err != nil || page <= 0 {
		page = 1
	}
	limit, err := strconv.Atoi(query.Get("limit"))
	if err != nil || limit <= 0 {
		limit = 10
	}

	return domain.Pagination{
		Limit:  limit,
		Offset: (page - 1) * limit,
	}
}
