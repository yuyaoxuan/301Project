package client

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CreateUserHandler handles the HTTP request to create a user
func CreateClientHandler(w http.ResponseWriter, r *http.Request) {
	var client Client

	vars := mux.Vars(r)
	AgentID, stringToInt_err := strconv.Atoi(vars["agent_id"])

	if stringToInt_err != nil {
		http.Error(w, "Unable to convert string to int", http.StatusInternalServerError)
		return
	 }

	// Decode JSON request body into client struct
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Initialize repository & service
	repo := NewClientRepository() // ✅ This ensures the table exists before inserting data
	service := NewClientService(repo)

	// Call service layer to create client 
	createdClient, err := service.CreateClient(client, AgentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send back the created user as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(createdClient)
}