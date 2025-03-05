package user

import (
	"encoding/json"
	"net/http"
)

// CreateUserHandler handles the HTTP request to create a user
func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	// Decode JSON request body into user struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Initialize repository & service
	repo := NewUserRepository() // ✅ This ensures the table exists before inserting data
	service := NewUserService(repo)

	// Call service layer to create user
	createdUser, err := service.CreateUser(user.FirstName, user.LastName, user.Email, user.Role)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send back the created user as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdUser)
}
