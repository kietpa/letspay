package middleware

import (
	"context"
	"letspay/common/constants"
	"letspay/controller"
	"letspay/model"
	"letspay/tool/util"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			controller.RespondWithError(w, http.StatusUnauthorized, model.Error{
				Code:    http.StatusUnauthorized,
				Message: constants.MISSING_AUTH_HEADER_MESSAGE,
			})
			return
		}

		tokenString := authHeader[len("Bearer "):]

		// Parse token
		claims := &util.Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return util.JWTKey, nil
		})

		if err != nil || !token.Valid {
			controller.RespondWithError(w, http.StatusUnauthorized, model.Error{
				Code:    http.StatusUnauthorized,
				Message: constants.INVALID_TOKEN_MESSAGE,
			})
			return
		}

		ctx := context.WithValue(r.Context(), constants.USER_ID, claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
