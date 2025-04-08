package agentClient

import (
	"encoding/json"
	"net/http"
)

//KAI ZHE ONLY FOR ADMIN TOKEN
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

func GetUnassignedClientsHandler(w http.ResponseWriter, r *http.Request) {
	// Initialize repository & service
	repo := NewAgentClientRepository()
	service := NewAgentClientService(repo)

	// ✅ Call service layer to get unassigned clients
	clients, err := service.GetUnassignedClients()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ✅ Success response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	// Marshal clients into JSON and write the response
	err = json.NewEncoder(w).Encode(clients)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}