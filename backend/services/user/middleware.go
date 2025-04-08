package user

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

var (
	cognitoIssuer   = "https://cognito-idp.ap-southeast-1.amazonaws.com/ap-southeast-1_ZTmaj2omi"
	cognitoClientID = "2g45kkfdu9ba88uhksnv3c86uu"
	provider        *oidc.Provider
	verifier        *oidc.IDTokenVerifier
)

func init() {
	var err error
	provider, err = oidc.NewProvider(context.Background(), cognitoIssuer)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize OIDC provider: %v", err))
	}
	verifier = provider.Verifier(&oidc.Config{ClientID: cognitoClientID})
}

func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		rawIDToken := strings.TrimPrefix(authHeader, "Bearer ")

		idToken, err := verifier.Verify(r.Context(), rawIDToken)
		if err != nil {
			http.Error(w, "Failed to verify ID token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		var claims map[string]interface{}
		if err := idToken.Claims(&claims); err != nil {
			http.Error(w, "Failed to parse claims: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Extract role from Cognito groups
		var role string
		if groups, ok := claims["cognito:groups"].([]interface{}); ok {
			for _, group := range groups {
				if strGroup, ok := group.(string); ok && (strGroup == "Admin" || strGroup == "Agent") {
					role = strGroup
					break
				}
			}
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), "user", map[string]interface{}{
			"id":   claims["sub"],
			"role": role,
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
