package agentclient_logs

import (
	"backend/models"
	"fmt"
)

type CommunicationNotifier interface {
	NotifyCommunication(agentID int, clientID string, log models.AgentClientLog)
}

// AgentClientLogService handles log operations
type AgentClientLogService struct {
	repo     *AgentClientLogRepository
	notifier CommunicationNotifier // only this interface, not the whole ObserverManager
}

// NewAgentClientLogService initializes the service
func NewAgentClientLogService(repo *AgentClientLogRepository, notifier CommunicationNotifier) *AgentClientLogService {
	return &AgentClientLogService{repo: repo, notifier: notifier}
}

// LogAgentClientAction processes and stores agent-client logs
func (s *AgentClientLogService) LogAgentClientAction(agentID int, clientID string, action string, modifiedFields map[string]interface{}) (models.AgentClientLog, error) {
	if action == "" {
		return models.AgentClientLog{}, fmt.Errorf("missing action type")
	}

	// Pass the correct types to the repository
	log, err := s.repo.CreateAgentClientLog(agentID, clientID, action, modifiedFields)
	if err != nil {
		return log, err
	}

	s.notifier.NotifyCommunication(agentID, clientID, log)
	return log, nil
}

// GetAgentClientLogs retrieves logs for a specific client
func (s *AgentClientLogService) GetAgentClientLogs(clientID string) ([]models.AgentClientLog, error) {
	return s.repo.GetClientTransactionLogsByClientID(clientID)
}

// GetAgentClientLogsByAgent retrieves logs for a specific agent
func (s *AgentClientLogService) GetAgentClientLogsByAgent(agentID int) ([]models.AgentClientLog, error) {
	return s.repo.GetClientTransactionLogsByAgentID(agentID)
}

// GetAllAgentClientLogs retrieves all logs
func (s *AgentClientLogService) GetAllAgentClientLogs() ([]models.AgentClientLog, error) {
	return s.repo.GetAllClientTransactionLogs()
}

// LogAccountChange inserts a new bank account log into the database
func (s *AgentClientLogService) LogAccountChange(agentID int, clientID string, action string, bankAccountInfo map[string]interface{}) error {
	return s.repo.LogAccountChange(agentID, clientID, action, bankAccountInfo)
}

// GetAccountLogsByClientID retrieves all bank account logs for a specific client
func (s *AgentClientLogService) GetAccountLogsByClientID(clientID string) ([]models.AgentClientLog, error) {
	return s.repo.GetAccountLogsByClientID(clientID)
}

// GetAccountLogsByAgentID retrieves all bank account logs for a specific agent
func (s *AgentClientLogService) GetAccountLogsByAgentID(agentID int) ([]models.AgentClientLog, error) {
	return s.repo.GetAccountLogsByAgentID(agentID)
}

// GetAllAccountLogs retrieves all bank account transaction logs
func (s *AgentClientLogService) GetAllAccountLogs() ([]models.AgentClientLog, error) {
	return s.repo.GetAllAccountLogs()
}

// GetClientAndAccountLogsByAgentID retreieves client and account logs by agent ID
func (s *AgentClientLogService) GetClientAndAccountLogsByAgentID(agentID int) ([]models.AgentClientLog, error) {
	return s.repo.GetClientAndAccountLogsByAgentID(agentID)
}

// GetClientAndAccountLogsByClientID retreives client and account logs by client ID
func (s *AgentClientLogService) GetClientAndAccountLogsByClientID(clientID string) ([]models.AgentClientLog, error) {
	return s.repo.GetClientAndAccountLogsByClientID(clientID)
}

// GetAllLogs retrieves all logs (client and bank account logs)
func (s *AgentClientLogService) GetAllLogs() ([]models.AgentClientLog, error) {
	return s.repo.GetAllLogs()
}

// DeleteLog deletes any log by its ID (either client or bank account)
func (s *AgentClientLogService) DeleteLog(logID int) error {
	// Call the repository to delete the log (it could be of any type)
	return s.repo.DeleteLog(logID)
}
