package agentClient

import (
	"encoding/json"
	"net/http"
)

// ✅ ONLY Admins can assign clients — role is verified from JWT token context
func AssignAgentsToUnassignedClientsHandler(w http.ResponseWriter, r *http.Request) {
	// ✅ Extract role from context (set by middleware)
	userCtx := r.Context().Value("user")
	if userCtx == nil {
		http.Error(w, "Unauthorized: missing token context", http.StatusUnauthorized)
		return
	}
	claims := userCtx.(map[string]interface{})
	role := claims["role"].(string)

	if role != "Admin" {
		http.Error(w, "Forbidden: only Admins can assign clients", http.StatusForbidden)
		return
	}

	// ✅ Continue with logic
	repo := NewAgentClientRepository()
	service := NewAgentClientService(repo)

	err := service.AssignAgentsToUnassignedClients()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Clients have been assigned to agents"))
}

// ✅ ONLY Admins can view unassigned clients
func GetUnassignedClientsHandler(w http.ResponseWriter, r *http.Request) {
	// ✅ Extract role from context (set by middleware)
	userCtx := r.Context().Value("user")
	if userCtx == nil {
		http.Error(w, "Unauthorized: missing token context", http.StatusUnauthorized)
		return
	}
	claims := userCtx.(map[string]interface{})
	role := claims["role"].(string)

	if role != "Admin" {
		http.Error(w, "Forbidden: only Admins can view unassigned clients", http.StatusForbidden)
		return
	}

	// ✅ Continue with logic
	repo := NewAgentClientRepository()
	service := NewAgentClientService(repo)

	clients, err := service.GetUnassignedClients()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(clients)
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}
