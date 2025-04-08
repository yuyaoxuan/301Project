package middleware

import (
	"context"
	"net/http"
	"strings"
	"backend/services/user"
	"backend/services/auth"
)

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		rawIDToken := strings.TrimPrefix(authHeader, "Bearer ")

		idToken, err := auth.Verifier().Verify(r.Context(), rawIDToken)
		if err != nil {
			http.Error(w, "Failed to verify ID token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		var claims map[string]interface{}
		if err := idToken.Claims(&claims); err != nil {
			http.Error(w, "Failed to parse claims: "+err.Error(), http.StatusInternalServerError)
			return
		}

		email, ok := claims["email"].(string)
		if !ok || email == "" {
			http.Error(w, "Missing email in token claims", http.StatusUnauthorized)
			return
		}

		role := "Agent"
		if groups, ok := claims["cognito:groups"].([]interface{}); ok {
			for _, group := range groups {
				if strGroup, ok := group.(string); ok && (strGroup == "Admin" || strGroup == "Agent") {
					role = strGroup
					break
				}
			}
		}

		// ✅ DB sync now handled by user package
		userID, err := user.SyncOrInsertUser(email, role)
		if err != nil {
			http.Error(w, "Failed to sync user: "+err.Error(), http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), "user", map[string]interface{}{
			"id":    userID,
			"email": email,
			"role":  role,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
