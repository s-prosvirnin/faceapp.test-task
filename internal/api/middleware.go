package api

import (
	"context"
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

func (m Middleware) AuthRequest(nextHandler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			accessToken := r.Header.Get("Authorization")
			splitToken := strings.Split(accessToken, "Bearer ")
			accessToken = splitToken[1]
			// @todo: убрать team_id из запросов, а брать его по токену
			teamId, err := m.service.GetTeamIdByAuthToken(accessToken)

			if err != nil {
				writeErrorResponse(w, errors.Wrap(err, "service.GetTeamIdByAuthToken"))

				return
			}
			ctx := context.WithValue(r.Context(), teamIdContextKey, teamId)

			nextHandler.ServeHTTP(w, r.WithContext(ctx))
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
