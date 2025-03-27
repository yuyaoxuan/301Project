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
	return s.repo.GetClientTransactionLogsByClientID(clientID)
}

// GetAgentClientLogsByAgent retrieves logs for a specific agent
func (s *AgentClientLogService) GetAgentClientLogsByAgent(agentID int) ([]AgentClientLog, error) {
	return s.repo.GetClientTransactionLogsByAgentID(agentID)
}

// GetAllAgentClientLogs retrieves all logs
func (s *AgentClientLogService) GetAllAgentClientLogs() ([]AgentClientLog, error) {
	return s.repo.GetAllClientTransactionLogs()
}

// LogAccountChange inserts a new bank account log into the database
func (s *AgentClientLogService) LogAccountChange(agentID int, clientID string, action string, bankAccountInfo map[string]interface{}) error {
	return s.repo.LogAccountChange(agentID, clientID, action, bankAccountInfo)
}

// GetAccountLogsByClientID retrieves all bank account logs for a specific client
func (s *AgentClientLogService) GetAccountLogsByClientID(clientID string) ([]AgentClientLog, error) {
	return s.repo.GetAccountLogsByClientID(clientID)
}

// GetAccountLogsByAgentID retrieves all bank account logs for a specific agent
func (s *AgentClientLogService) GetAccountLogsByAgentID(agentID int) ([]AgentClientLog, error) {
	return s.repo.GetAccountLogsByAgentID(agentID)
}

// GetAllAccountLogs retrieves all bank account transaction logs
func (s *AgentClientLogService) GetAllAccountLogs() ([]AgentClientLog, error) {
	return s.repo.GetAllAccountLogs()
}

// GetAllLogs retrieves all logs (client and bank account logs)
func (s *AgentClientLogService) GetAllLogs() ([]AgentClientLog, error) {
	return s.repo.GetAllLogs()
}

// DeleteLog deletes any log by its ID (either client or bank account)
func (s *AgentClientLogService) DeleteLog(logID int) error {
	// Call the repository to delete the log (it could be of any type)
	return s.repo.DeleteLog(logID)
}
