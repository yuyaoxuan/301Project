package client

import (
	"fmt"
	"strings"
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


// GetClient retrieves a client by ID
func (s *Clientservice) GetClient(clientID string) (Client, error) {
    // Validate client ID
    if clientID == "" {
        return Client{}, fmt.Errorf("client ID cannot be empty")
    }
    
    // Call repository to get client
    client, err := s.repo.GetClientByID(clientID)
    if err != nil {
        return Client{}, fmt.Errorf("failed to retrieve client: %v", err)
    }
    
    return client, nil
}

// UpdateClient updates client information
func (s *Clientservice) UpdateClient(client Client) (Client, error) {
    // Validate client data
    if err := validateClient(client); err != nil {
        return Client{}, err
    }
    
    // Call repository to update client
    updatedClient, err := s.repo.UpdateClient(client)
    if err != nil {
        return Client{}, fmt.Errorf("failed to update client: %v", err)
    }
    
    return updatedClient, nil
}

// DeleteClient removes a client profile
func (s *Clientservice) DeleteClient(clientID string) error {
    // Validate client ID
    if clientID == "" {
        return fmt.Errorf("client ID cannot be empty")
    }
    
    // Call repository to delete client
    err := s.repo.DeleteClient(clientID)
    if err != nil {
        return fmt.Errorf("failed to delete client: %v", err)
    }
    
    return nil
}

// VerifyClient verifies a client's identity
func (s *Clientservice) VerifyClient(clientID string, nric string) error {
    // Validate inputs
    if clientID == "" {
        return fmt.Errorf("client ID cannot be empty")
    }
    
    if nric == "" {
        return fmt.Errorf("NRIC cannot be empty")
    }
    
    // In a real application, you would validate the NRIC against external systems
    
    // Call repository to verify client
    err := s.repo.VerifyClient(clientID)
    if err != nil {
        return fmt.Errorf("failed to verify client: %v", err)
    }
    
    return nil
}

// Helper function to validate client data
func validateClient(client Client) error {
    // First name and last name validation
    if len(client.FirstName) < 2 || len(client.FirstName) > 50 {
        return fmt.Errorf("first name must be between 2 and 50 characters")
    }
    
    if len(client.LastName) < 2 || len(client.LastName) > 50 {
        return fmt.Errorf("last name must be between 2 and 50 characters")
    }
    
    // Email validation (basic check)
    if client.Email == "" || !strings.Contains(client.Email, "@") {
        return fmt.Errorf("invalid email format")
    }
    
    // Phone validation (basic check)
    if len(client.Phone) < 10 || len(client.Phone) > 15 {
        return fmt.Errorf("phone number must be between 10 and 15 digits")
    }
    
    // Address validation
    if len(client.Address) < 5 || len(client.Address) > 100 {
        return fmt.Errorf("address must be between 5 and 100 characters")
    }
    
    // Other validations as per Appendix 2
    // ...
    
    return nil
}
