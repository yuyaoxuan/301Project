package agentclient_logs

import (
	"backend/database"
	"encoding/json"
	"fmt"
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
func (r *AgentClientLogRepository) CreateAgentClientLog(agentID int, clientID string, action string, modifiedFields map[string]interface{}) error {
	modifiedFieldsJSON, err := json.Marshal(modifiedFields)
	if err != nil {
		return fmt.Errorf("failed to convert modified fields to JSON: %v", err)
	}

	query := "INSERT INTO agent_client_logs (agent_id, client_id, action, modified_fields) VALUES (?, ?, ?, ?)"
	_, err = database.DB.Exec(query, agentID, clientID, action, string(modifiedFieldsJSON))
	if err != nil {
		return fmt.Errorf("failed to insert agent-client log: %v", err)
	}
	return nil
}

// GetAgentClientLogs retrieves all agent-client logs for a specific client
func (r *AgentClientLogRepository) GetAgentClientLogs(clientID string) ([]AgentClientLog, error) {
	query := "SELECT id, agent_id, client_id, action, modified_fields, timestamp FROM agent_client_logs WHERE client_id = ?"
	rows, err := database.DB.Query(query, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve logs for client: %v", err)
	}
	defer rows.Close()

	var logs []AgentClientLog
	for rows.Next() {
		var log AgentClientLog
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

// GetAgentClientLogsByAgent retrieves all agent-client logs for a specific agent
func (r *AgentClientLogRepository) GetAgentClientLogsByAgent(agentID string) ([]AgentClientLog, error) {
	query := "SELECT id, agent_id, client_id, action, modified_fields, timestamp FROM agent_client_logs WHERE agent_id = ?"
	rows, err := database.DB.Query(query, agentID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve logs for agent: %v", err)
	}
	defer rows.Close()

	var logs []AgentClientLog
	for rows.Next() {
		var log AgentClientLog
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

// GetAllAgentClientLogs retrieves all agent-client logs
func (r *AgentClientLogRepository) GetAllAgentClientLogs() ([]AgentClientLog, error) {
	query := "SELECT id, agent_id, client_id, action, modified_fields, timestamp FROM agent_client_logs"
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all agent-client logs: %v", err)
	}
	defer rows.Close()

	var logs []AgentClientLog
	for rows.Next() {
		var log AgentClientLog
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

// DeleteAgentClientLog deletes an agent-client log by its ID
func (r *AgentClientLogRepository) DeleteAgentClientLog(logID int) error {
	query := "DELETE FROM agent_client_logs WHERE id = ?"
	_, err := database.DB.Exec(query, logID)
	if err != nil {
		return fmt.Errorf("failed to delete agent-client log: %v", err)
	}
	return nil
}
