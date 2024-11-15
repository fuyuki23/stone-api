package web

import (
	"net/http"

	"github.com/google/uuid"
)

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-Id")
		if requestID == "" {
			requestUUID, err := uuid.NewV7()
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
			r.Header.Set("X-Request-ID", requestUUID.String())
			w.Header().Set("X-Request-ID", requestUUID.String())
		}

		next.ServeHTTP(w, r)
	})
}
