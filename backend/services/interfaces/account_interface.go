// interfaces/account_interface.go
package interfaces

import "backend/models"

// AccountServiceInterface defines the methods that an AccountService must implement
type AccountServiceInterface interface {
	GetAccountByClientId(clientID string) ([]models.Account, error)
	// Add other methods as needed
}