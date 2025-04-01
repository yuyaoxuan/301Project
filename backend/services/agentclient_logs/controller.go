package agentclient_logs

import (
	"backend/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CreateAgentClientLogHandler handles log creation requests for agent-client logs
func CreateAgentClientLogHandler(service *AgentClientLogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var logData models.AgentClientLog
		err := json.NewDecoder(r.Body).Decode(&logData)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if logData.ClientID == "" {
			http.Error(w, "Invalid client_id, must be a non-empty string", http.StatusBadRequest)
			return
		}

		// Call the service
		createdLog, err := service.LogAgentClientAction(logData.AgentID, logData.ClientID, logData.Action, logData.ModifiedFields)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdLog)
	}
}

// GetAgentClientLogsByClientHandler retrieves transaction logs for a specific client
func GetAgentClientLogsByClientHandler(service *AgentClientLogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		clientID := vars["clientID"]

		logs, err := service.GetAgentClientLogs(clientID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logs)
	}
}

// GetAgentClientLogsByAgentHandler retrieves agent-client logs for a specific agent
func GetAgentClientLogsByAgentHandler(service *AgentClientLogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		agentID := vars["agentID"]
		agentIDInt, err := strconv.Atoi(agentID)
		if err != nil {
			http.Error(w, "Invalid agent ID, must be an integer", http.StatusBadRequest)
			return
		}

		logs, err := service.GetAgentClientLogsByAgent(agentIDInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logs)
	}
}

// GetAllAgentClientLogsHandler retrieves all agent-client logs
func GetAllAgentClientLogsHandler(service *AgentClientLogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logs, err := service.GetAllAgentClientLogs()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logs)
	}
}

// CreateBankAccountLogHandler handles the creation of bank account logs
func CreateBankAccountLogHandler(service *AgentClientLogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var logData models.AgentClientLog
		err := json.NewDecoder(r.Body).Decode(&logData)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		if logData.ClientID == "" {
			http.Error(w, "Invalid client_id, must be a non-empty string", http.StatusBadRequest)
			return
		}

		err = service.LogAccountChange(logData.AgentID, logData.ClientID, logData.Action, logData.ModifiedFields)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Bank account log created successfully"))
	}
}

// GetAccountLogsByClientHandler retrieves bank account logs for a specific client
func GetAccountLogsByClientHandler(service *AgentClientLogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		clientID := vars["clientID"]

		logs, err := service.GetAccountLogsByClientID(clientID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logs)
	}
}

// GetAccountLogsByAgentHandler retrieves bank account logs for a specific agent
func GetAccountLogsByAgentHandler(service *AgentClientLogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		agentID := vars["agentID"]
		agentIDInt, err := strconv.Atoi(agentID)
		if err != nil {
			http.Error(w, "Invalid agent ID, must be an integer", http.StatusBadRequest)
			return
		}

		logs, err := service.GetAccountLogsByAgentID(agentIDInt)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logs)
	}
}

// GetAllAccountLogsHandler retrieves all bank account transaction logs
func GetAllAccountLogsHandler(service *AgentClientLogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logs, err := service.GetAllAccountLogs()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logs)
	}
}

// GetClientAndAccountLogsByAgentIDHandler retrieves both client and account logs for a specific agent
func GetClientAndAccountLogsByAgentIDHandler(service *AgentClientLogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		agentID, err := strconv.Atoi(vars["agentID"])
		if err != nil {
			http.Error(w, "Invalid agent ID", http.StatusBadRequest)
			return
		}

		logs, err := service.GetClientAndAccountLogsByAgentID(agentID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logs)
	}
}

// GetClientAndAccountLogsByClientIDHandler retrieves both client and account logs for a specific client
func GetClientAndAccountLogsByClientIDHandler(service *AgentClientLogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		clientID := vars["clientID"]

		logs, err := service.GetClientAndAccountLogsByClientID(clientID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logs)
	}
}

// GetAllLogsHandler retrieves all logs from the database (client and bank account logs)
func GetAllLogsHandler(service *AgentClientLogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logs, err := service.GetAllLogs()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(logs)
	}
}

// DeleteLogHandler handles deleting any log by its ID
func DeleteLogHandler(service *AgentClientLogService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		logID, err := strconv.Atoi(vars["logID"])
		if err != nil {
			http.Error(w, "Invalid log ID", http.StatusBadRequest)
			return
		}

		err = service.DeleteLog(logID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Log deleted successfully"))
	}
}
