package interfaces

import "backend/models"

// AgentClientServiceInterface defines the methods that the AgentClientService must implement
type AgentClientServiceInterface interface {
    GetUnassignedClients() ([]models.AgentClient, error)
}
