package communicationlogs

import (
	"backend/database"
	"backend/models"
	"database/sql"
	"fmt"
)

// CommunicationLogRepository handles database operations
type CommunicationLogRepository struct{}

// NewCommunicationLogRepository initializes the repository
func NewCommunicationLogRepository() *CommunicationLogRepository {
	repo := &CommunicationLogRepository{}
	repo.InitTable()
	return repo
}

// InitTable ensures the communication_logs table exists
func (r *CommunicationLogRepository) InitTable() {
	query := `
	CREATE TABLE IF NOT EXISTS communication_logs (
		id INT AUTO_INCREMENT PRIMARY KEY,
		log_id INT NOT NULL,
		client_id VARCHAR(255) NOT NULL,  -- client_id is now a string
		agent_id INT NOT NULL,
		email_subject VARCHAR(255) NOT NULL,
		email_status ENUM('Sent', 'Failed') NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := database.DB.Exec(query)
	if err != nil {
		fmt.Println("Error creating communication_logs table:", err)
	} else {
		fmt.Println("communication_logs table checked/created!")
	}
}

// InsertCommunicationLog inserts a new communication log
func (r *CommunicationLogRepository) InsertCommunicationLog(logID int, clientID string, agentID int, emailSubject, emailStatus string) error {
	query := `
		INSERT INTO communication_logs 
		(log_id, client_id, agent_id, email_subject, email_status) 
		VALUES (?, ?, ?, ?, ?)`
	_, err := database.DB.Exec(query, logID, clientID, agentID, emailSubject, emailStatus)
	if err != nil {
		return fmt.Errorf("failed to insert communication log: %v", err)
	}
	return nil
}

// GetCommunicationLogByLogID retrieves a specific communication log by log ID
func (r *CommunicationLogRepository) GetCommunicationLogByLogID(logID int) (models.CommunicationLog, error) {
	query := "SELECT id, log_id, client_id, agent_id, email_subject, email_status, timestamp FROM communication_logs WHERE log_id = ?"
	row := database.DB.QueryRow(query, logID) // Use QueryRow for single result

	var log models.CommunicationLog
	if err := row.Scan(&log.ID, &log.LogID, &log.ClientID, &log.AgentID, &log.EmailSubject, &log.EmailStatus, &log.Timestamp); err != nil {
		if err == sql.ErrNoRows {
			return log, fmt.Errorf("no communication log found with log_id %d", logID)
		}
		return log, fmt.Errorf("failed to retrieve communication log: %v", err)
	}

	return log, nil
}
