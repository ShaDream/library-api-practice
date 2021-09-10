package handler

import (
	"context"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("handle \"%s\" request", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

func GetTimeoutMiddleware(duration time.Duration) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Context()
			timeoutContext, cancel := context.WithTimeout(r.Context(), duration)
			defer cancel()
			newR := r.WithContext(timeoutContext)
			next.ServeHTTP(w, newR)
		})
	}
}
