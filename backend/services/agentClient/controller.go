package agentClient

import (
	"net/http"
)


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
