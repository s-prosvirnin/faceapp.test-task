package api

import (
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type Middleware struct {
	service Service
}

func NewMiddleware(service Service) Middleware {
	return Middleware{service: service}
}

// AuthRequest check auth token for access
func (m Middleware) AuthRequest(nextHandler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			accessToken := r.Header.Get("Authorization")
			splitToken := strings.Split(accessToken, "Bearer ")
			accessToken = splitToken[1]
			err := m.service.CheckAuth(accessToken)

			if err != nil {
				writeErrorResponse(w, errors.Wrap(err, "service.CheckAuth"))

				return
			}

			nextHandler.ServeHTTP(w, r)
		},
	)
}

func (m Middleware) MutateResponseHeaders(handler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-type", "application/json")
			handler.ServeHTTP(w, r)
		},
	)
}
