package agentclient_logs

import (
	"backend/database"
	"encoding/json"
	"fmt"

	"backend/models"
)

// AgentClientLogRepository handles database operations
type AgentClientLogRepository struct{}

// NewAgentClientLogRepository initializes the repository
func NewAgentClientLogRepository() *AgentClientLogRepository {
	repo := &AgentClientLogRepository{}
	repo.InitTable()
	return repo
}

// InitTable ensures the agent_client_logs table exists
func (r *AgentClientLogRepository) InitTable() {
	query := `
	CREATE TABLE IF NOT EXISTS agent_client_logs (
		id INT AUTO_INCREMENT PRIMARY KEY,
		agent_id INT NOT NULL,
		client_id VARCHAR(255) NOT NULL,
		action VARCHAR(50) NOT NULL,
		modified_fields JSON NOT NULL,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := database.DB.Exec(query)
	if err != nil {
		fmt.Println("❌ Error creating agent_client_logs table:", err)
	} else {
		fmt.Println("✅ agent_client_logs table checked/created!")
	}
}

// CreateAgentClientLog inserts a new agent-client log into the database
func (r *AgentClientLogRepository) CreateAgentClientLog(agentID int, clientID string, action string, modifiedFields map[string]interface{}) (models.AgentClientLog, error) {
	// No need to create timestamp manually, MySQL will do it for you
	logData := models.AgentClientLog{
		AgentID:        agentID,
		ClientID:       clientID,
		Action:         action,
		ModifiedFields: map[string]interface{}{"log_type": "client", "details": modifiedFields["details"]},
		// No need to pass Timestamp here, MySQL will fill it automatically
	}

	// Convert modified fields to JSON
	modifiedFieldsJSON, err := json.Marshal(logData.ModifiedFields)
	if err != nil {
		return models.AgentClientLog{}, fmt.Errorf("failed to convert modified fields to JSON: %v", err)
	}

	// Insert the log into the agent_client_logs table
	query := "INSERT INTO agent_client_logs (agent_id, client_id, action, modified_fields) VALUES (?, ?, ?, ?)"
	result, err := database.DB.Exec(query, agentID, clientID, action, modifiedFieldsJSON)
	if err != nil {
		return models.AgentClientLog{}, fmt.Errorf("failed to insert agent-client log: %v", err)
	}

	// Get the auto-generated log_id (the ID of the newly inserted log)
	logID, err := result.LastInsertId()
	if err != nil {
		return models.AgentClientLog{}, fmt.Errorf("failed to get the last inserted log_id: %v", err)
	}

	// Update the logData with the generated logID
	logData.ID = int(logID) // Type assertion since `LastInsertId()` returns int64

	return logData, nil
}

// GetClientTransactionLogsByClientID retrieves all client logs for a specific client
func (r *AgentClientLogRepository) GetClientTransactionLogsByClientID(clientID string) ([]models.AgentClientLog, error) {
	query := `SELECT id, agent_id, client_id, action, modified_fields, timestamp FROM agent_client_logs
		WHERE client_id = ? AND JSON_UNQUOTE(JSON_EXTRACT(modified_fields, '$.log_type')) = 'client'`
	rows, err := database.DB.Query(query, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve logs for client: %v", err)
	}
	defer rows.Close()

	var logs []models.AgentClientLog
	for rows.Next() {
		var log models.AgentClientLog
		var modifiedFieldsJSON string

		if err := rows.Scan(&log.ID, &log.AgentID, &log.ClientID, &log.Action, &modifiedFieldsJSON, &log.Timestamp); err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(modifiedFieldsJSON), &log.ModifiedFields)
		if err != nil {
			return nil, fmt.Errorf("failed to decode modified fields: %v", err)
		}

		logs = append(logs, log)
	}
	return logs, nil
}

// GetClientTransactionLogsByAgentID retrieves all client logs for a specific agent
func (r *AgentClientLogRepository) GetClientTransactionLogsByAgentID(agentID int) ([]models.AgentClientLog, error) {
	query := `SELECT id, agent_id, client_id, action, modified_fields, timestamp FROM agent_client_logs
		WHERE agent_id = ? AND JSON_UNQUOTE(JSON_EXTRACT(modified_fields, '$.log_type')) = 'client'`
	rows, err := database.DB.Query(query, agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve logs for agent: %v", err)
	}
	defer rows.Close()

	var logs []models.AgentClientLog
	for rows.Next() {
		var log models.AgentClientLog
		var modifiedFieldsJSON string

		if err := rows.Scan(&log.ID, &log.AgentID, &log.ClientID, &log.Action, &modifiedFieldsJSON, &log.Timestamp); err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(modifiedFieldsJSON), &log.ModifiedFields)
		if err != nil {
			return nil, fmt.Errorf("failed to decode modified fields: %v", err)
		}

		logs = append(logs, log)
	}
	return logs, nil
}

// GetAllClientTransactionLogs retrieves all client transaction logs
func (r *AgentClientLogRepository) GetAllClientTransactionLogs() ([]models.AgentClientLog, error) {
	query := `SELECT id, agent_id, client_id, action, modified_fields, timestamp FROM agent_client_logs
		WHERE JSON_UNQUOTE(JSON_EXTRACT(modified_fields, '$.log_type')) = 'client'`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all agent-client logs: %v", err)
	}
	defer rows.Close()

	var logs []models.AgentClientLog
	for rows.Next() {
		var log models.AgentClientLog
		var modifiedFieldsJSON string

		if err := rows.Scan(&log.ID, &log.AgentID, &log.ClientID, &log.Action, &modifiedFieldsJSON, &log.Timestamp); err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(modifiedFieldsJSON), &log.ModifiedFields)
		if err != nil {
			return nil, fmt.Errorf("failed to decode modified fields: %v", err)
		}

		logs = append(logs, log)
	}
	return logs, nil
}

// LogAccountChange inserts a new bank account log into the database
func (r *AgentClientLogRepository) LogAccountChange(agentID int, clientID string, action string, bankAccountInfo map[string]interface{}) error {
	// Log data for bank account
	logData := models.AgentClientLog{
		AgentID:        agentID,
		ClientID:       clientID,
		Action:         action,                                                                                    // "Create", "Update", "Delete"
		ModifiedFields: map[string]interface{}{"log_type": "bank_account", "details": bankAccountInfo["details"]}, // "log_type": "bank_account"
		// Don't manually set Timestamp, MySQL will handle it with CURRENT_TIMESTAMP
	}

	// Convert modified fields to JSON
	modifiedFieldsJSON, err := json.Marshal(logData.ModifiedFields)
	if err != nil {
		return fmt.Errorf("failed to convert modified fields to JSON: %v", err)
	}

	// Insert the log into the agent_client_logs table, without passing the timestamp
	query := `
		INSERT INTO agent_client_logs (agent_id, client_id, action, modified_fields)
		VALUES (?, ?, ?, ?)
	`

	_, err = database.DB.Exec(query, agentID, clientID, action, modifiedFieldsJSON)
	if err != nil {
		return fmt.Errorf("failed to create log: %v", err)
	}

	return nil
}

// GetAccountLogsByClientID retrieves all bank account logs for a specific client
func (r *AgentClientLogRepository) GetAccountLogsByClientID(clientID string) ([]models.AgentClientLog, error) {
	query := `
		SELECT id, agent_id, client_id, action, modified_fields, timestamp
		FROM agent_client_logs
		WHERE client_id = ? AND JSON_UNQUOTE(JSON_EXTRACT(modified_fields, '$.log_type')) = 'bank_account'
	`
	rows, err := database.DB.Query(query, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve bank account logs for client by client_id: %v", err)
	}
	defer rows.Close()

	var logs []models.AgentClientLog
	for rows.Next() {
		var log models.AgentClientLog
		var modifiedFieldsJSON string

		if err := rows.Scan(&log.ID, &log.AgentID, &log.ClientID, &log.Action, &modifiedFieldsJSON, &log.Timestamp); err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(modifiedFieldsJSON), &log.ModifiedFields)
		if err != nil {
			return nil, fmt.Errorf("failed to decode modified fields: %v", err)
		}

		logs = append(logs, log)
	}
	return logs, nil
}

// GetAccountLogsByAgentID retrieves all bank account logs for a specific agent
func (r *AgentClientLogRepository) GetAccountLogsByAgentID(agentID int) ([]models.AgentClientLog, error) {
	query := `
		SELECT id, agent_id, client_id, action, modified_fields, timestamp
		FROM agent_client_logs
		WHERE agent_id = ? AND JSON_UNQUOTE(JSON_EXTRACT(modified_fields, '$.log_type')) = 'bank_account'
	`
	rows, err := database.DB.Query(query, agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve bank account logs for agent by agent_id: %v", err)
	}
	defer rows.Close()

	var logs []models.AgentClientLog
	for rows.Next() {
		var log models.AgentClientLog
		var modifiedFieldsJSON string

		if err := rows.Scan(&log.ID, &log.AgentID, &log.ClientID, &log.Action, &modifiedFieldsJSON, &log.Timestamp); err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(modifiedFieldsJSON), &log.ModifiedFields)
		if err != nil {
			return nil, fmt.Errorf("failed to decode modified fields: %v", err)
		}

		logs = append(logs, log)
	}
	return logs, nil
}

// GetAllAccountLogs retrieves all bank account transaction logs
func (r *AgentClientLogRepository) GetAllAccountLogs() ([]models.AgentClientLog, error) {
	query := `
		SELECT id, agent_id, client_id, action, modified_fields, timestamp
		FROM agent_client_logs
		WHERE JSON_UNQUOTE(JSON_EXTRACT(modified_fields, '$.log_type')) = 'bank_account'
	`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all bank account transaction logs: %v", err)
	}
	defer rows.Close()

	var logs []models.AgentClientLog
	for rows.Next() {
		var log models.AgentClientLog
		var modifiedFieldsJSON string

		if err := rows.Scan(&log.ID, &log.AgentID, &log.ClientID, &log.Action, &modifiedFieldsJSON, &log.Timestamp); err != nil {
			return nil, err
		}

		err = json.Unmarshal([]byte(modifiedFieldsJSON), &log.ModifiedFields)
		if err != nil {
			return nil, fmt.Errorf("failed to decode modified fields: %v", err)
		}

		logs = append(logs, log)
	}
	return logs, nil
}

// GetClientAndAccountLogsByAgentID retrieves both client and account logs for a specific agent
func (r *AgentClientLogRepository) GetClientAndAccountLogsByAgentID(agentID int) ([]models.AgentClientLog, error) {
	query := `
		SELECT id, agent_id, client_id, action, modified_fields, timestamp
		FROM agent_client_logs
		WHERE agent_id = ? AND JSON_UNQUOTE(JSON_EXTRACT(modified_fields, '$.log_type')) IN ('client', 'bank_account')
	`
	rows, err := database.DB.Query(query, agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve client and account logs for agent by agent_id: %v", err)
	}
	defer rows.Close()

	var logs []models.AgentClientLog
	for rows.Next() {
		var log models.AgentClientLog
		var modifiedFieldsJSON string

		if err := rows.Scan(&log.ID, &log.AgentID, &log.ClientID, &log.Action, &modifiedFieldsJSON, &log.Timestamp); err != nil {
			return nil, err
		}

		// Decode the modified fields (details)
		err = json.Unmarshal([]byte(modifiedFieldsJSON), &log.ModifiedFields)
		if err != nil {
			return nil, fmt.Errorf("failed to decode modified fields: %v", err)
		}

		logs = append(logs, log)
	}
	return logs, nil
}

// GetClientAndAccountLogsByClientID retrieves both client and account logs for a specific client
func (r *AgentClientLogRepository) GetClientAndAccountLogsByClientID(clientID string) ([]models.AgentClientLog, error) {
	query := `
		SELECT id, agent_id, client_id, action, modified_fields, timestamp
		FROM agent_client_logs
		WHERE client_id = ? AND JSON_UNQUOTE(JSON_EXTRACT(modified_fields, '$.log_type')) IN ('client', 'bank_account')
	`
	rows, err := database.DB.Query(query, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve client and account logs for client by client_id: %v", err)
	}
	defer rows.Close()

	var logs []models.AgentClientLog
	for rows.Next() {
		var log models.AgentClientLog
		var modifiedFieldsJSON string

		if err := rows.Scan(&log.ID, &log.AgentID, &log.ClientID, &log.Action, &modifiedFieldsJSON, &log.Timestamp); err != nil {
			return nil, err
		}

		// Decode the modified fields (details)
		err = json.Unmarshal([]byte(modifiedFieldsJSON), &log.ModifiedFields)
		if err != nil {
			return nil, fmt.Errorf("failed to decode modified fields: %v", err)
		}

		logs = append(logs, log)
	}
	return logs, nil
}

// GetAllLogs retrieves all logs (both client and bank account logs)
func (r *AgentClientLogRepository) GetAllLogs() ([]models.AgentClientLog, error) {
	// Use SELECT to explicitly filter log types 'client' and 'bank_account'
	query := `
		SELECT id, agent_id, client_id, action, modified_fields, timestamp
		FROM agent_client_logs
		WHERE JSON_UNQUOTE(JSON_EXTRACT(modified_fields, '$.log_type')) IN ('client', 'bank_account')
	`

	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all logs: %v", err)
	}
	defer rows.Close()

	var logs []models.AgentClientLog
	for rows.Next() {
		var log models.AgentClientLog
		var modifiedFieldsJSON string
		var timestamp string // Store the timestamp as string

		// Scan the log data from the database
		if err := rows.Scan(&log.ID, &log.AgentID, &log.ClientID, &log.Action, &modifiedFieldsJSON, &timestamp); err != nil {
			return nil, err
		}

		log.Timestamp = timestamp // Store the timestamp as a string

		// Decode the modified fields (info)
		err = json.Unmarshal([]byte(modifiedFieldsJSON), &log.ModifiedFields)
		if err != nil {
			return nil, fmt.Errorf("failed to decode modified fields: %v", err)
		}

		logs = append(logs, log)
	}
	return logs, nil
}

// DeleteLog deletes an agent-client log by its ID
func (r *AgentClientLogRepository) DeleteLog(logID int) error {
	query := "DELETE FROM agent_client_logs WHERE id = ?"
	_, err := database.DB.Exec(query, logID)
	if err != nil {
		return fmt.Errorf("failed to delete agent-client log: %v", err)
	}
	return nil
}
