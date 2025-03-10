package routes

import (
	"net/http"

	"backend/services/user" // import user service

	"github.com/gorilla/mux"
)

// SetupRoutes initializes the router and returns it
func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is running!"))
	}).Methods("GET")

	// USER ROUTES
	r.HandleFunc("/users", user.CreateUserHandler).Methods("POST")              // Create User
	r.HandleFunc("/api/users/{userId}", user.DisableUserHandler).Methods("DELETE")  // Disable User
	r.HandleFunc("/api/users/{userId}", user.UpdateUserHandler).Methods("PUT")      // Update User
	r.HandleFunc("/api/users/authenticate", user.AuthenticateUserHandler).Methods("POST") // Authenticate User (Login)
	r.HandleFunc("/api/users/reset-password", user.ResetPasswordHandler).Methods("POST")  // Reset Password


	// CLIENT ROUTES

	return r
}
