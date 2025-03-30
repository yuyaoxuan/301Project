package account

import (
	"backend/models"
	"backend/services/agentClient"
	"backend/services/client"
	"backend/services/observer"
	"fmt"
)

// UserService struct to interact with the repository layer
type AccountService struct {
	ObserverManager *observer.ObserverManager
	repo *AccountRepository
	AgentClientService *agentClient.AgentClientService
	ClientService *client.ClientService
}

// NewUserService initializes the user service
func NewAccountService(observerManager *observer.ObserverManager, repo *AccountRepository, agentClientService *agentClient.AgentClientService, clientService *client.ClientService) *AccountService {
	return &AccountService{
		ObserverManager: observerManager, // Pass the ObserverManager here
		repo: repo, 
		AgentClientService: agentClientService,
		ClientService: clientService,
		}
}

func (s *AccountService) ClientExists(clientID string) (bool, error) {
	client, err := s.ClientService.GetClient(clientID)
	if err != nil {
		return false, fmt.Errorf("failed to get agent ID for client %q: %v", clientID, err)
	}

    return client.ClientID != "", nil
}

func (s *AccountService)  GetAgentIDByClientID(clientID string) (int, error) {
	agentID, err := s.AgentClientService.GetAgentIDByClientID(clientID)
	if err != nil {
        return 0, fmt.Errorf("failed to check if client exists: %v", err)
    }
    return agentID, nil
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

	// Call repository function to insert account
	createdAccount, err := s.repo.CreateAccount(account)
	if err != nil {
		return models.Account{}, fmt.Errorf("failed to create account: %v", err)
	}

	// Notify observers after client update
	if s.ObserverManager != nil {
		s.ObserverManager.NotifyAccountCreate(account.AccountID, account.ClientID, &account)
	}

	return createdAccount, nil
}


func (s *AccountService) DeleteAccount(AccountID int) (error) {
	// Check if account_id exists before proceeding
	account, err := s.repo.GetAccountByID(AccountID)
	if err != nil {
		return fmt.Errorf("failed to check account id existence: %d", err)
	}

	// If the account is not found, we return an error
	if account.AccountID == 0 {
		return fmt.Errorf("account with ID %d does not exist", AccountID)
	}

	// get agent info 
	agentID, err_agentID := s.AgentClientService.GetAgentIDByClientID(account.ClientID)
	if err_agentID != nil {
	    return fmt.Errorf("failed to check agent id existence: %d", err_agentID)
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


