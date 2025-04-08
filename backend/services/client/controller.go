package client

import (
	"backend/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

//KAI ZHE DO CHECK AGENTID and OR ADMINID with TOKEN
func CreateClientHandler(service *ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var client models.Client
		vars := mux.Vars(r)
		AgentID, stringToInt_err := strconv.Atoi(vars["agent_id"])

		if stringToInt_err != nil {
			http.Error(w, "Unable to convert agent_id to int", http.StatusInternalServerError)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&client)
		if err != nil {
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

func GetClientHandler(clientService *ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		clientID := vars["clientId"]

		client, err := clientService.GetClient(clientID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(client)
	}
}

func UpdateClientHandler(clientService *ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		clientID := vars["clientId"]
		agentID, err := strconv.Atoi(vars["agent_id"])
		if err != nil {
			http.Error(w, "Invalid agent ID", http.StatusBadRequest)
			return
		}

		var client models.Client
		if err := json.NewDecoder(r.Body).Decode(&client); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}
		client.ClientID = clientID

		updatedClient, err := clientService.UpdateClient(client, agentID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedClient)
	}
}

func DeleteClientHandler(clientService *ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		clientID := vars["clientId"]

		err := clientService.DeleteClient(clientID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Client deleted successfully"})
	}
}

func VerifyClientHandler(clientService *ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		clientID := vars["clientId"]

		var requestBody struct {
			NRIC string `json:"nric"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		err := clientService.VerifyClient(clientID, requestBody.NRIC)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Client verified successfully"})
	}
}

func GetAllClientsHandler(clientService *ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clients, err := clientService.GetAllClients()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clients)
	}
}

//KAI ZHE DO CHECK AGENTID and OR ADMINID with TOKEN
func GetClientsByAgentHandler(clientService *ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		agentID, err := strconv.Atoi(vars["agentId"])
		if err != nil {
			http.Error(w, "Invalid agent ID", http.StatusBadRequest)
			return
		}

		clients, err := clientService.GetClientsByAgentID(agentID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(clients)
	}
}

func GetUnassignedClientsHandler(clientService *ClientService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Call the service layer to get unassigned clients
		clients, err := clientService.GetUnassignedClients()
		if err != nil {
			// If there is an error fetching the clients, return a 500 Internal Server Error
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Step 2: Set response header to JSON
		w.Header().Set("Content-Type", "application/json")
		// Step 3: Encode the clients into JSON and send the response
		err = json.NewEncoder(w).Encode(clients)
		if err != nil {
			// If there is an error encoding the response, return a 500 Internal Server Error
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
		}
	}
}


