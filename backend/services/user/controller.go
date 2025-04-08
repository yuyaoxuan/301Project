package user

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"backend/services/auth"
	"github.com/gorilla/mux"
	"golang.org/x/oauth2"
)

var oauthConfig = oauth2.Config{
	ClientID:     auth.ClientID,
	ClientSecret: auth.ClientSecret,
	RedirectURL:  "http://localhost:8080/api/auth/callback",
	Endpoint: oauth2.Endpoint{
		AuthURL:   auth.AuthURL,
		TokenURL:  auth.TokenURL,
		AuthStyle: oauth2.AuthStyleInParams,
	},
	Scopes: []string{"openid", "email", "profile"},
}


// AuthenticateUserHandler redirects user to Cognito login
func AuthenticateUserHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, oauthConfig.AuthCodeURL("state"), http.StatusFound)
}

// AuthCallbackHandler handles Cognito login callback
func AuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Authentication failed: "+err.Error(), http.StatusUnauthorized)
		return
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		http.Error(w, "Missing ID token", http.StatusInternalServerError)
		return
	}

	idToken, err := auth.Verifier().Verify(r.Context(), rawIDToken)
	if err != nil {
		http.Error(w, "Failed to verify ID token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var claims struct {
		Email  string   `json:"email"`
		Groups []string `json:"cognito:groups"`
	}
	if err := idToken.Claims(&claims); err != nil {
		http.Error(w, "Failed to parse claims: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"access_token":  token.AccessToken,
		"id_token":      rawIDToken,
		"refresh_token": token.RefreshToken,
		"expires_in":    strconv.FormatInt(int64(token.Expiry.Sub(time.Now()).Seconds()), 10),
		"user_info": map[string]interface{}{
			"email":  claims.Email,
			"groups": claims.Groups,
		},
	})
}

// CreateUserHandler registers in Cognito, then stores user metadata in DB
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
		Role      string `json:"role"` // Admin or Agent
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if input.Role != "Admin" && input.Role != "Agent" {
		http.Error(w, "Role must be either Admin or Agent", http.StatusBadRequest)
		return
	}

	// Register in Cognito
	err := auth.RegisterUserInCognito(input.Email, input.Password, input.Role)
	if err != nil {
		http.Error(w, "Cognito registration failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Save metadata in local DB
	user, err := NewUserService(NewUserRepository()).CreateUser(
		input.FirstName,
		input.LastName,
		input.Email,
		input.Role,
	)
	if err != nil {
		http.Error(w, "Failed to save user metadata: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// DisableUserHandler uses JWT to restrict access
func DisableUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userId"]

	userCtx := r.Context().Value("user").(map[string]interface{})
	requesterID := userCtx["id"].(int)
	requesterRole := userCtx["role"].(string)

	err := NewUserService(NewUserRepository()).DisableUser(userID, requesterID, requesterRole)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User disabled successfully"))
}

// UpdateUserHandler for Admins to update metadata
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userId"]

	userCtx := r.Context().Value("user").(map[string]interface{})
	requesterRole := userCtx["role"].(string)
	if requesterRole != "Admin" {
		http.Error(w, "Only admins can update users", http.StatusForbidden)
		return
	}

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
