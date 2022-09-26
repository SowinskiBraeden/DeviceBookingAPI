package api

import (
	"net/http"
)

// Middleware adds some basic header authentication around accessing the routes
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do Authentication stuff?
		next.ServeHTTP(w, r)
	})
}
