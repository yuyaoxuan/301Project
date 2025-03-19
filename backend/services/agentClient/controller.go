package agentClient

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// CreateUserHandler handles the HTTP request to create a user
func UpdateAgentToClientHandler(w http.ResponseWriter, r *http.Request) {
	var agent struct {
		NewID    int    `json:"new_id"`
	}

	vars := mux.Vars(r)
	ClientID := vars["client_id"]

	// Decode JSON request body into client struct
	err := json.NewDecoder(r.Body).Decode(&agent)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Initialize repository & service
	repo := NewAgentClientRepository()
	service := NewAgentClientService(repo)

	// ✅ Call service layer to update agent ID
	updatedClient, err := service.UpdateAgentToClient(ClientID, agent.NewID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send back the updated client as a response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedClient)
}

func AssignAgentsToUnassignedClientsHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize repository & service
	repo := NewAgentClientRepository()
	service := NewAgentClientService(repo)

	// ✅ Call service layer to assign agents
	err := service.AssignAgentsToUnassignedClients()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ✅ Success response
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Clients have been assigned to agents"))
}
