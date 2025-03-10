package user
//First layer interface API calls and response first comes in
// Controller - > Service - > Repo
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
// To Disable User 
func DisableUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	repo := NewUserRepository()
	service := NewUserService(repo)

	err := service.DisableUser(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User disabled successfully"))
}

// To update user
func UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	var updatedUser User
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	repo := NewUserRepository()
	service := NewUserService(repo)

	err = service.UpdateUser(userID, updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User updated successfully"))
}

//Authenticate User 
func AuthenticateUserHandler(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	repo := NewUserRepository()
	service := NewUserService(repo)

	token, err := service.AuthenticateUser(credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

//Reset password
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var resetRequest struct {
		Email       string `json:"email"`
		NewPassword string `json:"newPassword"`
	}

	err := json.NewDecoder(r.Body).Decode(&resetRequest)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	repo := NewUserRepository()
	service := NewUserService(repo)

	err = service.ResetPassword(resetRequest.Email, resetRequest.NewPassword)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Password reset successfully"))
}
