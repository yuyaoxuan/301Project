package client

import (
	"backend/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Helper to validate access for agent/admin roles
func IsClientOwnedByAgent(agentID int, clientID string, role string, service *ClientService) bool {
	if role == "Admin" {
		return true
	}
	isOwned, err := service.IsClientOwnedByAgent(clientID, agentID)
	if err != nil {
		return false
	}
	return isOwned
}


// CreateClientHandler handles the creation of a client
func CreateClientHandler(service *ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var client models.Client
		vars := mux.Vars(r)
		AgentID, err := strconv.Atoi(vars["agent_id"])
		if err != nil {
			http.Error(w, "Unable to convert agent_id to int", http.StatusInternalServerError)
			return
		}

		// Role check
		userCtx := r.Context().Value("user").(map[string]interface{})
		tokenAgentID := userCtx["id"].(int)
		role := userCtx["role"].(string)
		if !IsClientOwnedByAgent(tokenAgentID, "", role, service) && tokenAgentID != AgentID {
			http.Error(w, "Unauthorized: not your agent ID", http.StatusForbidden)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		createdClient, err := service.CreateClient(client, AgentID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdClient)
	}
}

func GetClientHandler(service *ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		clientID := vars["clientId"]

		userCtx := r.Context().Value("user").(map[string]interface{})
		agentID := userCtx["id"].(int)
		role := userCtx["role"].(string)
		if !IsClientOwnedByAgent(agentID, clientID, role, service) {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}

		client, err := service.GetClient(clientID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(client)
	}
}

func UpdateClientHandler(service *ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		clientID := vars["clientId"]

		userCtx := r.Context().Value("user").(map[string]interface{})
		agentID := userCtx["id"].(int)
		role := userCtx["role"].(string)

		if !IsClientOwnedByAgent(agentID, clientID, role, service) {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}

		var client models.Client
		if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		client.ClientID = clientID

		updatedClient, err := service.UpdateClient(client, agentID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedClient)
	}
}

func DeleteClientHandler(service *ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		clientID := vars["clientId"]

		userCtx := r.Context().Value("user").(map[string]interface{})
		agentID := userCtx["id"].(int)
		role := userCtx["role"].(string)

		if !IsClientOwnedByAgent(agentID, clientID, role, service) {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}

		err := service.DeleteClient(clientID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Client deleted successfully"})
	}
}

func VerifyClientHandler(service *ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		clientID := vars["clientId"]

		userCtx := r.Context().Value("user").(map[string]interface{})
		agentID := userCtx["id"].(int)
		role := userCtx["role"].(string)

		if !IsClientOwnedByAgent(agentID, clientID, role, service) {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}

		var requestBody struct {
			NRIC string `json:"nric"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		err := service.VerifyClient(clientID, requestBody.NRIC)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Client verified successfully"})
	}
}

func GetAllClientsHandler(service *ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userCtx := r.Context().Value("user").(map[string]interface{})
		role := userCtx["role"].(string)

		if role != "Admin" {
			http.Error(w, "Only admins can access all clients", http.StatusForbidden)
			return
		}

		clients, err := service.GetAllClients()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clients)
	}
}

func GetClientsByAgentHandler(service *ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		requestedAgentID, err := strconv.Atoi(vars["agentId"])
		if err != nil {
			http.Error(w, "Invalid agent ID", http.StatusBadRequest)
			return
		}

		userCtx := r.Context().Value("user").(map[string]interface{})
		requesterID := userCtx["id"].(int)
		role := userCtx["role"].(string)

		if role != "Admin" && requesterID != requestedAgentID {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}

		clients, err := service.GetClientsByAgentID(requestedAgentID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clients)
	}
}

func GetUnassignedClientsHandler(service *ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userCtx := r.Context().Value("user").(map[string]interface{})
		role := userCtx["role"].(string)

		if role != "Admin" {
			http.Error(w, "Only admins can view unassigned clients", http.StatusForbidden)
			return
		}

		clients, err := service.GetUnassignedClients()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clients)
	}
}
