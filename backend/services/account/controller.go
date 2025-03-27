package account

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CreateUserHandler handles the HTTP request to create a user
func CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var account Account

	// Decode JSON request body into account struct
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Initialize repository & service 
	repo := NewAccountRepository()  // ✅ This ensures the table exists before inserting data
	service := NewAccountService(repo)

	// Call service layer to create account
	createdAccount, err := service.CreateAccount(account)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send back the created user as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdAccount)
}

// CreateUserHandler handles the HTTP request to create a user
func DeleteAccountHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	AccountID, stringToInt_err := strconv.Atoi(vars["account_id"])

	if stringToInt_err != nil {
		http.Error(w, "Unable to convert string to int", http.StatusInternalServerError)
		return
	 }

	// Initialize repository & service 
	repo := NewAccountRepository()  // ✅ This ensures the table exists before inserting data
	service := NewAccountService(repo)

	// Call service layer to create account
	err := service.DeleteAccount(AccountID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}