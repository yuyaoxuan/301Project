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


// GetClientHandler handles the HTTP request to get a client
func GetClientHandler(w http.ResponseWriter, r *http.Request) {
    // Extract client ID from URL path
    vars := mux.Vars(r)
    clientID := vars["clientId"]
    
    // Initialize repository & service
    repo := NewClientRepository()
    service := NewClientService(repo)
    
    // Call service to get client
    client, err := service.GetClient(clientID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    
    // Return client as JSON response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(client)
}

// UpdateClientHandler handles the HTTP request to update a client
func UpdateClientHandler(w http.ResponseWriter, r *http.Request) {
    // Extract client ID from URL path
    vars := mux.Vars(r)
    clientID := vars["clientId"]
    
    // Decode JSON request body
    var client Client
    err := json.NewDecoder(r.Body).Decode(&client)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }
    
    // Ensure client ID in URL matches client ID in body
    client.ClientID = clientID
    
    // Initialize repository & service
    repo := NewClientRepository()
    service := NewClientService(repo)
    
    // Call service to update client
    updatedClient, err := service.UpdateClient(client)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Return updated client as JSON response
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(updatedClient)
}

// DeleteClientHandler handles the HTTP request to delete a client
func DeleteClientHandler(w http.ResponseWriter, r *http.Request) {
    // Extract client ID from URL path
    vars := mux.Vars(r)
    clientID := vars["clientId"]
    
    // Initialize repository & service
    repo := NewClientRepository()
    service := NewClientService(repo)
    
    // Call service to delete client
    err := service.DeleteClient(clientID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusNotFound)
        return
    }
    
    // Return success message
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Client deleted successfully"})
}

// VerifyClientHandler handles the HTTP request to verify a client
func VerifyClientHandler(w http.ResponseWriter, r *http.Request) {
    // Extract client ID from URL path
    vars := mux.Vars(r)
    clientID := vars["clientId"]
    
    // Decode JSON request body to get NRIC
    var requestBody struct {
        NRIC string `json:"nric"`
    }
    
    err := json.NewDecoder(r.Body).Decode(&requestBody)
    if err != nil {
        http.Error(w, "Invalid request payload", http.StatusBadRequest)
        return
    }
    
    // Initialize repository & service
    repo := NewClientRepository()
    service := NewClientService(repo)
    
    // Call service to verify client
    err = service.VerifyClient(clientID, requestBody.NRIC)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // Return success message
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"message": "Client verified successfully"})
}
