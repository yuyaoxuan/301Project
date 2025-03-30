package account

import (
	"backend/database"
	"backend/models"
	"backend/services/agentClient"
	"backend/services/client"
	"backend/services/observer"
	"database/sql"

	"fmt"
	"log"
)

// UserRepository struct for interacting with database
type AccountRepository struct{
	ObserverManager *observer.ObserverManager
	ClientRepository *client.ClientRepository
	AgentClientRepository *agentClient.AgentClientRepository
}

// NewClientRepository initializes a new ClientRepository
func NewAccountRepository(observerManager *observer.ObserverManager, clientRepository *client.ClientRepository, agentClientRepository *agentClient.AgentClientRepository) *AccountRepository {
	repo := &AccountRepository{ObserverManager: observerManager, ClientRepository: clientRepository, AgentClientRepository: agentClientRepository}
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
	// Check if account exists
	account, err := r.GetAccountByID(accountID)
	if err != nil {
		return err
	}

	// If the account is not found, we return an error
	if account.AccountID == 0 {
		return fmt.Errorf("account with ID %d does not exist", accountID)
	}

	// Proceed with deleting the account
	query := `DELETE FROM account WHERE account_id = ?`

	_, err = database.DB.Exec(query, accountID)
	if err != nil {
		return fmt.Errorf("failed to delete account: %v", err)
	}
	
	return nil
}

// GetAccountByID retrieves an account by accountID
func (r *AccountRepository) GetAccountByID(account_id int) (models.Account, error) {
	query := `SELECT * FROM account WHERE account_id = ?`

	var account models.Account
	err := database.DB.QueryRow(query, account_id).Scan(
		&account.AccountID,
		&account.ClientID,
		&account.AccountType,
		&account.AccountStatus,
		&account.OpeningDate,
		&account.InitialDeposit,
		&account.Currency,
		&account.BranchID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Account{}, fmt.Errorf("account with ID %d does not exist", account_id)
		}
		return models.Account{}, fmt.Errorf("failed to retrieve account: %v", err)
	}

	return account, nil
}

