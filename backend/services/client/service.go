package client

import (
	"backend/models"
	"backend/services/observer"
	"fmt"
	"regexp"
	"time"
)

// Predefined gender options
var validGenders = map[string]bool{
	"Male":              true,
	"Female":            true,
	"Non-binary":        true,
	"Prefer not to say": true,
}

// ClientService struct to interact with the repository layer
type ClientService struct {
	repo *ClientRepository
	ObserverManager *observer.ObserverManager
}

// NewClientService initializes the client service
func NewClientService(repo *ClientRepository) *ClientService {
	return &ClientService{
        repo: repo,
        ObserverManager: &observer.ObserverManager{}, // Initialize with an empty ObserverManager
    }
}

// CreateUser processes user creation request
func (s *ClientService) CreateClient(client models.Client, AgentID int) (models.Client, error) {
	// ✅ Check if agent exists
	exists, err := s.repo.AgentExists(AgentID)
	if err != nil {
		return models.Client{}, fmt.Errorf("failed to check agent existence: %v", err)
	}

	if !exists {
		return models.Client{}, fmt.Errorf("agent's id not found")
	}

	if err := validateClient(client); err != nil {
		return models.Client{}, err
	}

	// Check for existing email
	if exists, err := s.repo.EmailExists(client.Email); err != nil || exists {
		return models.Client{}, fmt.Errorf("email address already exists")
	}

	// Check for existing phone
	if exists, err := s.repo.PhoneExists(client.Phone); err != nil || exists {
		return models.Client{}, fmt.Errorf("phone number already exists")
	}

	// Call repository function to insert client
	createdClient, err := s.repo.CreateClient(client, AgentID)
	if err != nil {
		return models.Client{}, fmt.Errorf("failed to create client: %v", err)
	}

	// Add nil check
    if s.ObserverManager != nil {
		// Prepare the client data to be sent to the observer manager
		// `client` object here contains all the client details, which will be sent for logging/notification
        s.ObserverManager.NotifyClientCreate(AgentID, createdClient.ClientID, &createdClient)
    }
	
	return createdClient, nil
}

// GetClient retrieves a client by ID
func (s *ClientService) GetClient(clientID string) (models.Client, error) {
	if clientID == "" {
		return models.Client{}, fmt.Errorf("client ID cannot be empty")
	}

	client, err := s.repo.GetClientByID(clientID)
	if err != nil {
		return models.Client{}, fmt.Errorf("failed to retrieve client: %v", err)
	}

	return client, nil
}

// UpdateClient updates client information
func (s *ClientService) UpdateClient(client models.Client) (models.Client, error) {
	if err := validateClient(client); err != nil {
		return models.Client{}, err
	}

	updatedClient, err := s.repo.UpdateClient(client)
	if err != nil {
		return models.Client{}, fmt.Errorf("failed to update client: %v", err)
	}

	return updatedClient, nil
}

// DeleteClient removes a client profile
func (s *ClientService) DeleteClient(clientID string) error {
	if clientID == "" {
		return fmt.Errorf("client ID cannot be empty")
	}

	err := s.repo.DeleteClient(clientID)
	if err != nil {
		return fmt.Errorf("failed to delete client: %v", err)
	}

	return nil
}

// VerifyClient verifies a client's identity
func (s *ClientService) VerifyClient(clientID string, nric string) error {
	// Validate inputs
	if clientID == "" {
		return fmt.Errorf("client ID cannot be empty")
	}

	if nric == "" {
		return fmt.Errorf("NRIC cannot be empty")
	}

	// In a real application, you would validate the NRIC against external systems

	err := s.repo.VerifyClient(clientID)
	if err != nil {
		return fmt.Errorf("failed to verify client: %v", err)
	}

	return nil
}

// Helper function to validate client data
func validateClient(client models.Client) error {
	// Name validations
	if err := validateName(client.FirstName, "first name"); err != nil {
		return err
	}
	if err := validateName(client.LastName, "last name"); err != nil {
		return err
	}

	// Date of Birth validation
	dob, err := time.Parse("2006-01-02", client.DOB)
	if err != nil {
		return fmt.Errorf("invalid date format, use YYYY-MM-DD")
	}
	if dob.After(time.Now().AddDate(-18, 0, 0)) {
		return fmt.Errorf("age must be at least 18 years")
	}
	if dob.Before(time.Now().AddDate(-100, 0, 0)) {
		return fmt.Errorf("age must be under 100 years")
	}

	// Gender validation
	if !validGenders[client.Gender] {
		return fmt.Errorf("invalid gender, must be one of: Male, Female, Non-binary, Prefer not to say")
	}

	// Email validation
	if !isValidEmail(client.Email) {
		return fmt.Errorf("invalid email format")
	}

	// Phone validation
	if !isValidPhone(client.Phone) {
		return fmt.Errorf("phone must start with + and contain 10-15 digits")
	}

	// Address validation
	if len(client.Address) < 5 || len(client.Address) > 100 {
		return fmt.Errorf("address must be between 5 and 100 characters")
	}

	// Location validations
	if err := validateLocation(client.City, "city"); err != nil {
		return err
	}
	if err := validateLocation(client.State, "state"); err != nil {
		return err
	}
	
	// Country validation
	if len(client.Country) < 2 || len(client.Country) > 50 {
		return fmt.Errorf("country must be between 2 and 50 characters")
	}

	// Postal code validation
	if len(client.PostalCode) < 4 || len(client.PostalCode) > 10 {
		return fmt.Errorf("postal code must be between 4 and 10 characters")
	}

	return nil
}

// Validation helper functions
func validateName(name, field string) error {
	if len(name) < 2 || len(name) > 50 {
		return fmt.Errorf("%s must be 2-50 characters", field)
	}
	if match, _ := regexp.MatchString(`^[a-zA-Z ]+$`, name); !match {
		return fmt.Errorf("%s can only contain letters and spaces", field)
	}
	return nil
}

func validateLocation(value, field string) error {
	if len(value) < 2 || len(value) > 50 {
		return fmt.Errorf("%s must be 2-50 characters", field)
	}
	if match, _ := regexp.MatchString(`^[a-zA-Z \-]+$`, value); !match {
		return fmt.Errorf("%s contains invalid characters", field)
	}
	return nil
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func isValidPhone(phone string) bool {
	phoneRegex := regexp.MustCompile(`^\+\d{10,15}$`)
	return phoneRegex.MatchString(phone)
}
