package routes

import (
	"net/http"

	"backend/services/user" // Import user service

	"github.com/gorilla/mux"
)

// SetupRoutes initializes the router and returns it
func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Health Check Route
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is running!"))
	}).Methods("GET")

	// Public Routes (No Authentication Needed)
	r.HandleFunc("/api/users/authenticate", user.AuthenticateUserHandler).Methods("POST") // Login
	r.HandleFunc("/api/users", user.CreateUserHandler).Methods("POST")                     // Register User
	r.HandleFunc("/api/users/reset-password", user.ResetPasswordHandler).Methods("POST")
	// Protected Routes (Require JWT)
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(user.JWTAuthMiddleware) // Apply JWT Middleware

	protected.HandleFunc("api//users/{userId}", user.DisableUserHandler).Methods("DELETE")  // Disable User
	protected.HandleFunc("api//users/{userId}", user.UpdateUserHandler).Methods("PUT")      // Update User
	

	return r
}
