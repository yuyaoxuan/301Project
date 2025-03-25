package user

import (
	"fmt"
	"net/http"
	"strings"
)

// JWTAuthMiddleware ensures requests have a valid JWT token
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization Header", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			http.Error(w, "Invalid Token", http.StatusUnauthorized)
			return
		}

		// Set user ID and role in request context (optional)
		r.Header.Set("UserID", fmt.Sprintf("%d", claims.UserID))
		r.Header.Set("Role", claims.Role)

		next.ServeHTTP(w, r)
	})
}
