package auth

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			http.Error(w, "Unauthorized - No token provided", http.StatusUnauthorized)
			return
		}

		tokenString := cookie.Value

		token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, http.ErrNoCookie
			}

			secretKey := os.Getenv("JWT_SECRET_KEY")
			if secretKey == "" {
				return nil, http.ErrNoCookie
			}

			return []byte(secretKey), nil
		})

		claims, ok := token.Claims.(*CustomClaims)
		if !ok || !token.Valid {
			http.Error(w, "Unauthorized - Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
