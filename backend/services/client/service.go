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
func (s *Clientservice) CreateClient(client Client) (Client, error) {
	// Call repository function to insert client
	createdClient, err := s.repo.CreateClient(client)
	if err != nil {
		return Client{}, fmt.Errorf("failed to create client: %v", err)
	}

	return createdClient, nil
}

func (s *Clientservice) CreateAccount(account Account) (Account, error) {
	// Call repository function to insert account
	createdAccount, err := s.repo.CreateAccount(account)
	if err != nil {
		return Account{}, fmt.Errorf("failed to create account: %v", err)
	}

	return createdAccount, nil
}