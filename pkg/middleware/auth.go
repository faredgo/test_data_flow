package middleware

import (
	"context"
	"net/http"
	"strings"
	"test_data_flow/configs"
	"test_data_flow/pkg/jwt"
	"test_data_flow/pkg/res"
)

type contextKey string

const (
	ContextKeyID    contextKey = "userID"
	ContextKeyLogin contextKey = "userLogin"
)

func IsAuthed(next http.Handler, cfg *configs.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authedHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authedHeader, "Bearer ") {
			res.ReturnError(w, ErrAuthHeader, http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(authedHeader, "Bearer ")

		data, err := jwt.NewJWT(cfg.Auth.Secret).Parse(token)
		if err != nil {
			res.ReturnError(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, ContextKeyID, data.ID)
		ctx = context.WithValue(ctx, ContextKeyLogin, data.Login)

		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
