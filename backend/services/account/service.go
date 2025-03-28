package account

import (
	"fmt"
)

// UserService struct to interact with the repository layer
type AccountService struct {
	repo *AccountRepository
	// ObserverManager *observer.ObserverManager
}

// NewUserService initializes the user service
func NewAccountService(repo *AccountRepository) *AccountService {
	return &AccountService{repo: repo}
}

func (s *AccountService) CreateAccount(account Account) (Account, error) {
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

	// Notify observers about the account creation (this will trigger AccountObserver)
	// s.ObserverManager.NotifyAccountCreate(account.AccountID, account.ClientID, &account)

	// Print success message
	// fmt.Println("Account created and observer notified.")

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
