package agentclient_logs

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CreateAgentClientLogHandler handles log creation requests for agent-client logs
func CreateAgentClientLogHandler(w http.ResponseWriter, r *http.Request) {
	var logData AgentClientLog
	err := json.NewDecoder(r.Body).Decode(&logData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Ensure client_id is a string
	if logData.ClientID == "" {
		http.Error(w, "Invalid client_id, must be a non-empty string", http.StatusBadRequest)
		return
	}

	repo := NewAgentClientLogRepository()
	service := NewAgentClientLogService(repo)

	// Call the service to log the action
	err = service.LogAgentClientAction(logData.AgentID, logData.ClientID, logData.Action, logData.ModifiedFields)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Agent-client log created successfully"))
}

// GetAgentClientLogsByClientHandler retrieves transaction logs for a specific client
func GetAgentClientLogsByClientHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientID := vars["clientID"]
	repo := NewAgentClientLogRepository()
	service := NewAgentClientLogService(repo)

	logs, err := service.GetAgentClientLogs(clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// GetAgentClientLogsByAgentHandler retrieves agent-client logs for a specific agent
func GetAgentClientLogsByAgentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentID := vars["agentID"]
	agentIDInt, err := strconv.Atoi(agentID)
	if err != nil {
		http.Error(w, "Invalid agent ID, must be an integer", http.StatusBadRequest)
		return
	}

	repo := NewAgentClientLogRepository()
	service := NewAgentClientLogService(repo)

	logs, err := service.GetAgentClientLogsByAgent(agentIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// GetAllAgentClientLogsHandler retrieves all agent-client logs
func GetAllAgentClientLogsHandler(w http.ResponseWriter, r *http.Request) {
	repo := NewAgentClientLogRepository()
	service := NewAgentClientLogService(repo)

	logs, err := service.GetAllAgentClientLogs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// CreateBankAccountLogHandler handles the creation of bank account logs
func CreateBankAccountLogHandler(w http.ResponseWriter, r *http.Request) {
	var logData AgentClientLog

	// Decode the request body to get the log data
	err := json.NewDecoder(r.Body).Decode(&logData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Ensure client_id is not empty
	if logData.ClientID == "" {
		http.Error(w, "Invalid client_id, must be a non-empty string", http.StatusBadRequest)
		return
	}

	repo := NewAgentClientLogRepository()
	service := NewAgentClientLogService(repo)

	// Call the repository function to insert the log (without timestamp)
	err = service.LogAccountChange(logData.AgentID, logData.ClientID, logData.Action, logData.ModifiedFields)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Bank account log created successfully"))
}

// GetAccountLogsByClientHandler retrieves bank account logs for a specific client
func GetAccountLogsByClientHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientID := vars["clientID"]
	repo := NewAgentClientLogRepository()
	service := NewAgentClientLogService(repo)

	logs, err := service.GetAccountLogsByClientID(clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// GetAccountLogsByAgentHandler retrieves bank account logs for a specific agent
func GetAccountLogsByAgentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentID := vars["agentID"]
	agentIDInt, err := strconv.Atoi(agentID)
	if err != nil {
		http.Error(w, "Invalid agent ID, must be an integer", http.StatusBadRequest)
		return
	}

	repo := NewAgentClientLogRepository()
	service := NewAgentClientLogService(repo)

	logs, err := service.GetAccountLogsByAgentID(agentIDInt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// GetAllAccountLogsHandler retrieves all bank account transaction logs
func GetAllAccountLogsHandler(w http.ResponseWriter, r *http.Request) {
	repo := NewAgentClientLogRepository()
	service := NewAgentClientLogService(repo)

	logs, err := service.GetAllAccountLogs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// GetAllLogsHandler retrieves all logs from the database (client and bank account logs)
func GetAllLogsHandler(w http.ResponseWriter, r *http.Request) {
	repo := NewAgentClientLogRepository()
	service := NewAgentClientLogService(repo)

	logs, err := service.GetAllLogs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// DeleteLogHandler handles deleting any log by its ID
func DeleteLogHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	logID, err := strconv.Atoi(vars["logID"]) // Extract the log ID from URL path
	if err != nil {
		http.Error(w, "Invalid log ID", http.StatusBadRequest)
		return
	}

	repo := NewAgentClientLogRepository()
	service := NewAgentClientLogService(repo)

	// Call the service to delete the log
	err = service.DeleteLog(logID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Log deleted successfully"))
}
