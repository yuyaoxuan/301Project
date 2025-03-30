package account

import (
	"backend/database"
	"backend/models"
	"backend/services/client"
	"backend/services/observer"

	"fmt"
	"log"
)

// UserRepository struct for interacting with database
type AccountRepository struct{
	ObserverManager *observer.ObserverManager
	clientService *client.ClientService

}

// NewClientRepository initializes a new ClientRepository
func NewAccountRepository(observerManager *observer.ObserverManager, clientService *client.ClientService) *AccountRepository {
	repo := &AccountRepository{ObserverManager: observerManager, clientService: clientService}
	repo.InitAccountTables() // ✅ Ensure tables exist when the repository is created
	return repo
}

// InitClientTables creates the client and account tables if they don't exist
func (r *AccountRepository) InitAccountTables() {
	// Create account table
	query := `
	CREATE TABLE IF NOT EXISTS account (
    account_id INT AUTO_INCREMENT PRIMARY KEY,
    client_id VARCHAR(50) NOT NULL,
    account_type ENUM('Savings', 'Checking', 'Business') NOT NULL DEFAULT 'Checking',
    account_status ENUM('Active', 'Inactive', 'Pending') NOT NULL DEFAULT 'Inactive',
    opening_date VARCHAR(50),
    initial_deposit FLOAT NOT NULL,
    currency VARCHAR(50) NOT NULL,
    branch_id VARCHAR(50) NOT NULL,
    FOREIGN KEY (client_id) REFERENCES client(client_id) ON DELETE CASCADE ON UPDATE CASCADE
);`

	_, err := database.DB.Exec(query)
	if err != nil {
		log.Fatal("❌ Error creating account table:", err)
	}

	fmt.Println("✅ Account tables checked/created!")
}


func (r *AccountRepository) CreateAccount(account models.Account) (models.Account, error) {

	clientExists, err := r.ClientExists(account.ClientID)
	if err != nil {
		return models.Account{}, fmt.Errorf("error checking if client exists: %v", err)
	}
	if !clientExists {
		return models.Account{}, fmt.Errorf("client with ID %s does not exist", account.ClientID)
	}

	query := `
	INSERT INTO account 
	(client_id, account_type, account_status, opening_date, initial_deposit, currency, branch_id) 
	VALUES (?, ?, ?, COALESCE(NULLIF(?, ''), CURDATE()), ?, ?, ?)`

	result, err := database.DB.Exec(query,
		account.ClientID, account.AccountType, account.AccountStatus,
		account.OpeningDate, account.InitialDeposit, account.Currency,
		account.BranchID,
	)
	if err != nil {
		return models.Account{}, fmt.Errorf("failed to insert account: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return models.Account{}, fmt.Errorf("failed to retrieve inserted account ID: %v", err)
	}

	account.AccountID = int(id)
	return account, nil
}

func (r *AccountRepository) DeleteAccount(accountID int) (error) {
	query := `DELETE FROM account WHERE account_id = ?`
	_, err := database.DB.Exec(query, accountID)

	if err != nil {
		return fmt.Errorf("failed to delete account: %v", err)
	}
	
	return nil
}

func (r *AccountRepository) ClientExists(clientID string) (bool, error) {
	client, err := r.clientService.GetClient(clientID)
	if err != nil {
        return false, fmt.Errorf("failed to check if client exists: %v", err)
    }
    return client.ClientID != "", nil
}


func (r *AccountRepository) AccountExists(account_id int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM account WHERE account_id = ?)`
	err := database.DB.QueryRow(query, account_id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check account id existence: %v", err)
	}
	return exists, nil
}
