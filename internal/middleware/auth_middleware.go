package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/JuanPidarraga/talkus-backend/internal/service"
)

type AuthMiddleware struct {
	authService *service.AuthService
}

func NewAuthMiddleware(authService *service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func (middleware *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		AuthHeader := r.Header.Get("Authorization")

		if AuthHeader == "" {
			http.Error(w, "❌ authorization header is empty", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(AuthHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.Error(w, "❌ authorization header format is invalid", http.StatusUnauthorized)
			return
		}

		decodedToken, err := middleware.authService.VerifyIDToken(r.Context(), tokenParts[1])
		if err != nil {
			http.Error(w, "❌ invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), AuthUserKey, decodedToken)

		next.ServeHTTP(w, r.WithContext(ctx))

	})
}
