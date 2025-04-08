package user

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
	"strconv"
	"github.com/coreos/go-oidc/v3/oidc"
)

var oauthConfig = oauth2.Config{
	ClientID:     "2g45kkfdu9ba88uhksnv3c86uu",
	ClientSecret: "1o647osm9clfsud2eturdb95hduj97cb61qgucvbp171s2jedbak",
	RedirectURL:  "http://localhost:8080/api/auth/callback",
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://ap-southeast-1ztmaj2omi.auth.ap-southeast-1.amazoncognito.com/oauth2/authorize",
		TokenURL: "https://ap-southeast-1ztmaj2omi.auth.ap-southeast-1.amazoncognito.com/oauth2/token",
		AuthStyle:   oauth2.AuthStyleInParams,
	},
	Scopes:      []string{"openid", "email", "profile"},
}


// AuthenticateUserHandler initiates Cognito login
func AuthenticateUserHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, oauthConfig.AuthCodeURL("state"), http.StatusFound)
}

// AuthCallbackHandler handles Cognito callback
func AuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Authentication failed: "+err.Error(), http.StatusUnauthorized)
		return
	}	
		 // Extract ID token claims to get user information
		 rawIDToken, ok := token.Extra("id_token").(string)
		 if !ok {
			 http.Error(w, "Missing ID token in response", http.StatusInternalServerError)
			 return
		 }
	 
		 // Parse the ID token to get claims
		 idToken, err := provider.Verifier(&oidc.Config{ClientID: cognitoClientID}).Verify(context.Background(), rawIDToken)
		 if err != nil {
			 http.Error(w, "Failed to verify ID token: "+err.Error(), http.StatusInternalServerError)
			 return
		 }

	// Get user claims
    var claims struct {
        Email string   `json:"email"`
        Groups []string `json:"cognito:groups"`
    }
    if err := idToken.Claims(&claims); err != nil {
        http.Error(w, "Failed to parse claims: "+err.Error(), http.StatusInternalServerError)
        return
    }


	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token":  token.AccessToken,
		"id_token":      token.Extra("id_token").(string),
		"refresh_token": token.RefreshToken,
		"expires_in":    strconv.FormatInt(int64(token.Expiry.Sub(time.Now()).Seconds()), 10),
		"user_info": map[string]interface{}{
			"email":  claims.Email,
			"groups": claims.Groups,
	}})
}

// CreateUserHandler handles user creation requests
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var newUser struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Role      string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	user, err := NewUserService(NewUserRepository()).CreateUser(
		newUser.FirstName,
		newUser.LastName,
		newUser.Email,
		newUser.Role,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// DisableUserHandler handles user deactivation
func DisableUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userId"]

	err := NewUserService(NewUserRepository()).DisableUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User disabled successfully"))
}

// UpdateUserHandler handles user updates
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userId"]
	var updatedUser User

	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := NewUserService(NewUserRepository()).UpdateUser(userID, updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User updated successfully"))
}

// ResetPasswordHandler handles password resets
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email       string `json:"email"`
		NewPassword string `json:"new_password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := NewUserService(NewUserRepository()).ResetPassword(req.Email, req.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password reset successfully"))
}
