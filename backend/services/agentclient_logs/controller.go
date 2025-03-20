package agentclient_logs

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CreateAgentClientLogHandler handles log creation requests
func CreateAgentClientLogHandler(w http.ResponseWriter, r *http.Request) {
	var logData AgentClientLog

	// Step 1: Decode the request body to get client data
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

	// Step 2: Call the log service to record the transaction with the proper data types
	repo := NewAgentClientLogRepository()
	service := NewAgentClientLogService(repo)

	// Call the service to log the action
	err = service.LogAgentClientAction(logData.AgentID, logData.ClientID, logData.Action, logData.ModifiedFields)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Step 3: Send success response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Agent-client log created successfully"))
}

// DeleteAgentClientLogHandler handles deleting a transaction log by ID
func DeleteAgentClientLogHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	logID, err := strconv.Atoi(vars["logID"])
	if err != nil {
		http.Error(w, "Invalid log ID", http.StatusBadRequest)
		return
	}

	repo := NewAgentClientLogRepository()
	service := NewAgentClientLogService(repo)

	err = service.DeleteAgentClientLog(logID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Agent-client log deleted successfully"))
}

// GetAgentClientLogsByClientHandler retrieves transaction logs for a specific client
func GetAgentClientLogsByClientHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientID := vars["clientID"] // Extract client_id from URL

	// Create service instance
	repo := NewAgentClientLogRepository()
	service := NewAgentClientLogService(repo)

	// Get agent-client logs for client
	logs, err := service.GetAgentClientLogs(clientID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the logs as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// GetAgentClientLogsByAgentHandler retrieves agent-client logs for a specific agent
func GetAgentClientLogsByAgentHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	agentID := vars["agentID"] // Extract agent_id from URL

	// Create service instance
	repo := NewAgentClientLogRepository()
	service := NewAgentClientLogService(repo)

	// Get agent-client logs for agent
	logs, err := service.GetAgentClientLogsByAgent(agentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the logs as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}

// GetAllAgentClientLogsHandler retrieves all agent-client logs
func GetAllAgentClientLogsHandler(w http.ResponseWriter, r *http.Request) {
	// Create service instance
	repo := NewAgentClientLogRepository()
	service := NewAgentClientLogService(repo)

	// Get all agent-client logs
	logs, err := service.GetAllAgentClientLogs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return the logs as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(logs)
}
