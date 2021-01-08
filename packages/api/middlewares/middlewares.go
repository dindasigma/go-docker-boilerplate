package middlewares

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/dindasigma/go-docker-boilerplate/packages/api/auth"
	"github.com/dindasigma/go-docker-boilerplate/packages/api/utils/responses"
)

func SetJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func SetAuthentication(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := auth.TokenValid(r)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(w, r)
	}
}

func SetTimeout(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
		defer cancel()

		// This gives you a copy of the request with a the request context
		// changed to the new context with the 5-second timeout created
		// above.
		r = r.WithContext(ctx)
		next(w, r)
	}
}
