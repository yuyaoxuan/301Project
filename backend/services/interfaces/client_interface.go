package interfaces

import "backend/models"

// ClientServiceInterface defines the methods that a ClientService must implement
type ClientServiceInterface interface {
	GetClient(clientID string) (models.Client, error)
	// Add other methods as needed
}

