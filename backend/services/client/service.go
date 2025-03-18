package client

import (
	"fmt"
)

// UserService struct to interact with the repository layer
type Clientservice struct {
	repo *ClientRepository
}

// NewUserService initializes the user service
func NewClientService(repo *ClientRepository) *Clientservice {
	return &Clientservice{repo: repo}
}

// CreateUser processes user creation request
func (s *Clientservice) CreateClient(client Client, AgentID int) (Client, error) {
	// ✅ Check if agent exists
	exists, err := s.repo.AgentExists(AgentID)
	if err != nil {
		return Client{}, fmt.Errorf("failed to check agent existence: %v", err)
	}

	if !exists {
		return Client{}, fmt.Errorf("agent's id not found")
	}

	// Call repository function to insert client
	createdClient, err := s.repo.CreateClient(client, AgentID)
	if err != nil {
		return Client{}, fmt.Errorf("failed to create client: %v", err)
	}

	return createdClient, nil
}