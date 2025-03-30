package account

import (
	"backend/models"
	"backend/services/observer"
	"fmt"
)

// UserService struct to interact with the repository layer
type AccountService struct {
	repo *AccountRepository
	ObserverManager *observer.ObserverManager
}

// NewUserService initializes the user service
func NewAccountService(repo *AccountRepository,  observerManager *observer.ObserverManager) *AccountService {
	return &AccountService{
		repo: repo, 
		ObserverManager: observerManager, // Pass the ObserverManager here
		}
}

func (s *AccountService) CreateAccount(account models.Account) (models.Account, error) {
	// Check if client_id exists before proceeding
	exists, err := s.repo.ClientExists(account.ClientID)
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
	exists, err := s.repo.AccountExists(AccountID)
	if err != nil {
		return fmt.Errorf("failed to check account id existence: %d", err)
	}
	if !exists {
		return fmt.Errorf("clientid %q does not exist", AccountID)
	}

	// Call repository function to insert account
	err = s.repo.DeleteAccount(AccountID)
	if err != nil {
		return fmt.Errorf("failed to create account: %v", err)
	}

	return nil
}
