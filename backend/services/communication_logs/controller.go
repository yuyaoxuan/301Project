package communicationlogs

import (
	"backend/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CreateCommunicationLogHandler handles email logging
func CreateCommunicationLogHandler(w http.ResponseWriter, r *http.Request) {
	var logData models.AgentClientLog // Now we expect an AgentClientLog instead of CommunicationLog

	err := json.NewDecoder(r.Body).Decode(&logData)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Initialize the repository and service
	repo := NewCommunicationLogRepository()
	service := NewCommunicationLogService(repo)

	// Pass the whole AgentClientLog to the service to process and send the email
	err = service.LogCommunication(logData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Communication log created and email sent successfully"))
}

// GetCommunicationLogByLogIDHandler retrieves a specific communication log by log ID
func GetCommunicationLogByLogIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	logID, err := strconv.Atoi(vars["logID"])
	if err != nil {
		http.Error(w, "Invalid log ID", http.StatusBadRequest)
		return
	}

	repo := NewCommunicationLogRepository()
	service := NewCommunicationLogService(repo)

	log, err := service.GetCommunicationLogByLogID(logID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(log)
}
