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
	// Check if client_id exists before proceeding
	exists, err := s.repo.ClientExists(account.ClientID)
	if err != nil {
		return Account{}, fmt.Errorf("failed to check client existence: %v", err)
	}
	if !exists {
		return Account{}, fmt.Errorf("client_id %q does not exist", account.ClientID)
	}

	// Call repository function to insert account
	createdAccount, err := s.repo.CreateAccount(account)
	if err != nil {
		return Account{}, fmt.Errorf("failed to create account: %v", err)
	}

	return createdAccount, nil
}


func (s *Clientservice) DeleteAccount(ClientID string) (error) {
	// Check if account_id exists before proceeding
	exists, err := s.repo.ClientExists(ClientID)
	if err != nil {
		return fmt.Errorf("failed to check account id existence: %d", err)
	}
	if !exists {
		return fmt.Errorf("clientid %q does not exist", ClientID)
	}

	// Call repository function to insert account
	err = s.repo.DeleteAccount(ClientID)
	if err != nil {
		return fmt.Errorf("failed to create account: %v", err)
	}

	return nil
}
