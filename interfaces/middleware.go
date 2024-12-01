package interfaces

import (
	"net/http"
)

func AdminOnlyMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assuming admin status is part of context
		isAdmin := r.Context().Value("is_admin").(bool)
		if !isAdmin {
			http.Error(w, "forbidden: admin access required", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
