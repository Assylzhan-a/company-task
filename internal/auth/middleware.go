package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"strings"

	"github.com/assylzhan-a/company-task/pkg/config"
	"github.com/assylzhan-a/company-task/pkg/errors"
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			errors.RespondWithError(w, errors.NewUnauthorizedError("Authorization header is required"))
			return
		}

		bearerToken := strings.Split(authHeader, " ")
		if len(bearerToken) != 2 || bearerToken[0] != "Bearer" {
			errors.RespondWithError(w, errors.NewUnauthorizedError("Invalid authorization header format"))
			return
		}

		claims := jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(bearerToken[1], claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Load().JWTSecret), nil
		})

		if err != nil {
			errors.RespondWithError(w, errors.NewUnauthorizedError("Invalid token"))
			return
		}

		if !token.Valid {
			errors.RespondWithError(w, errors.NewUnauthorizedError("Invalid token"))
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims["user_id"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
