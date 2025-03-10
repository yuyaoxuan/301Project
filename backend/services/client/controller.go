package client

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// CreateUserHandler handles the HTTP request to create a user
func CreateClientHandler(w http.ResponseWriter, r *http.Request) {
	var client Client

	// Decode JSON request body into client struct
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// TO-DO: Generate a unique client ID (e.g., UUID)

	// Initialize repository & service
	repo := NewClientRepository() // ✅ This ensures the table exists before inserting data
	service := NewClientService(repo)

	// Call service layer to create client 
	createdClient, err := service.CreateClient(client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send back the created user as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdClient)
}

// CreateUserHandler handles the HTTP request to create a user
func CreateAccountHandler(w http.ResponseWriter, r *http.Request) {
	var account Account

	// Decode JSON request body into account struct
	err := json.NewDecoder(r.Body).Decode(&account)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Initialize repository & service TODO: do we need seperate for accounts?
	repo := NewClientRepository()  // ✅ This ensures the table exists before inserting data
	service := NewClientService(repo)

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
	ClientID := vars["client_id"]

	// Initialize repository & service TODO: do we need seperate for accounts?
	repo := NewClientRepository()  // ✅ This ensures the table exists before inserting data
	service := NewClientService(repo)

	// Call service layer to create account
	err := service.DeleteAccount(ClientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}