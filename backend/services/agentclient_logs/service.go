package agentclient_logs

import (
	"fmt"
)

// AgentClientLogService handles log operations
type AgentClientLogService struct {
	repo *AgentClientLogRepository
}

// NewAgentClientLogService initializes the service
func NewAgentClientLogService(repo *AgentClientLogRepository) *AgentClientLogService {
	return &AgentClientLogService{repo: repo}
}

// LogAgentClientAction processes and stores agent-client logs
func (s *AgentClientLogService) LogAgentClientAction(agentID int, clientID string, action string, modifiedFields map[string]interface{}) error {
	if action == "" {
		return fmt.Errorf("missing action type")
	}

	// Pass the correct types to the repository
	return s.repo.CreateAgentClientLog(agentID, clientID, action, modifiedFields)
}

// GetAgentClientLogs retrieves logs for a specific client
func (s *AgentClientLogService) GetAgentClientLogs(clientID string) ([]AgentClientLog, error) {
	return s.repo.GetAgentClientLogs(clientID)
}

// GetAgentClientLogsByAgent retrieves logs for a specific agent
func (s *AgentClientLogService) GetAgentClientLogsByAgent(agentID string) ([]AgentClientLog, error) {
	return s.repo.GetAgentClientLogsByAgent(agentID)
}

// GetAllAgentClientLogs retrieves all logs
func (s *AgentClientLogService) GetAllAgentClientLogs() ([]AgentClientLog, error) {
	return s.repo.GetAllAgentClientLogs()
}

// DeleteAgentClientLog deletes an agent-client log by its ID
func (s *AgentClientLogService) DeleteAgentClientLog(logID int) error {
	return s.repo.DeleteAgentClientLog(logID)
}
