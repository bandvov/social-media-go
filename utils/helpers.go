package utils

import (
	"fmt"
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

// Helper to convert int64 slice to interface{} slice for query arguments
func ToInterface(ids []int64) []interface{} {
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}
	return args
}
