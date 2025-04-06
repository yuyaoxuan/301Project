package account

import (
	"backend/models"
	"backend/services/agentClient"
	"backend/services/interfaces"
	"backend/services/observer"
	"fmt"
)

// UserService struct to interact with the repository layer
type AccountService struct {
	ObserverManager *observer.ObserverManager
	repo *AccountRepository
	AgentClientService *agentClient.AgentClientService
	ClientService      interfaces.ClientServiceInterface
}

// NewUserService initializes the user service
func NewAccountService(observerManager *observer.ObserverManager, repo *AccountRepository, agentClientService *agentClient.AgentClientService) *AccountService {
	return &AccountService{
		ObserverManager: observerManager, // Pass the ObserverManager here
		repo: repo, 
		AgentClientService: agentClientService,
		}
}

// SetClientService sets the client service
func (s *AccountService) SetClientService(clientService interfaces.ClientServiceInterface) {
	s.ClientService = clientService
}

func (s *AccountService) ClientExists(clientID string) (bool, error) {
	client, err := s.ClientService.GetClient(clientID)
	if err != nil {
		return false, fmt.Errorf("failed to get agent ID for client %q: %v", clientID, err)
	}

    return client.ClientID != "", nil
}

func (s *AccountService) GetAgentIDByClientID(clientID string) (int, error) {
	agentID, err := s.AgentClientService.GetAgentIDByClientID(clientID)
	if err != nil {
        return 0, fmt.Errorf("failed to check if client exists: %v", err)
    }
    return agentID, nil
}

// Add GetAccountByClientId method to implement the interface
func (s *AccountService) GetAccountByClientId(clientID string) ([]models.Account, error) {
	return s.repo.GetAccountByClientId(clientID)
}

func (s *AccountService) CreateAccount(account models.Account) (models.Account, error) {
	// Check if client_id exists before proceeding
	exists, err := s.ClientExists(account.ClientID)
	if err != nil {
		return models.Account{}, fmt.Errorf("failed to check client existence: %v", err)
	}
	if !exists {
		return models.Account{}, fmt.Errorf("client_id %q does not exist", account.ClientID)
	}

	// Validate account_status
	validAccountStatuses := map[string]bool{"Active": true, "Inactive": true, "Pending": true}
	if _, valid := validAccountStatuses[account.AccountStatus]; !valid {
		return models.Account{}, fmt.Errorf("invalid account status: %s. Valid options are: 'Active', 'Inactive', 'Pending'", account.AccountStatus)
	}

	// Validate initial_deposit
	if account.InitialDeposit <= 0 {
		return models.Account{}, fmt.Errorf("initial deposit must be greater than 0")
	}

	// Missing values 
	if account.Currency == "" || account.BranchID == "" {
		return models.Account{}, fmt.Errorf("missing required fields")
	}

	// Call repository function to insert account
	createdAccount, err := s.repo.CreateAccount(account)
	if err != nil {
		return models.Account{}, fmt.Errorf("failed to create account: %v", err)
	}

	// get agent info 
	agentID, err_agentID := s.AgentClientService.GetAgentIDByClientID(account.ClientID)
	if err_agentID != nil {
	    return models.Account{}, fmt.Errorf("failed to check agent id existence: %v", err_agentID)
	}

	// Notify observers after client update
	if s.ObserverManager != nil {
		s.ObserverManager.NotifyAccountCreate(agentID, account.ClientID, &account)
	}

	return createdAccount, nil
}


func (s *AccountService) DeleteAccount(AccountID int) (error) {
	// Check if account_id exists before proceeding
	account, err := s.repo.GetAccountByID(AccountID)
	if err != nil {
		return fmt.Errorf("failed to check account id existence: %v", err)
	}

	// If the account is not found, we return an error
	if account.AccountID == 0 {
		return fmt.Errorf("account with ID %v does not exist", AccountID)
	}

	// get agent info 
	agentID, err_agentID := s.AgentClientService.GetAgentIDByClientID(account.ClientID)
	if err_agentID != nil {
	    return fmt.Errorf("failed to check agent id existence: %v", err_agentID)
	}

	// Notify observers after client update
	if s.ObserverManager != nil {
		fmt.Println(account.AccountID)
		s.ObserverManager.NotifyAccountDelete(agentID, account.ClientID, &account)
	}

	// Call repository function to delete account
	err = s.repo.DeleteAccount(AccountID)
	if err != nil {
		return fmt.Errorf("failed to delete account: %v", err)
	}
	
	return nil
}


