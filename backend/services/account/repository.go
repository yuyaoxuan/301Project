package account

import (
	"backend/database"
	"backend/models"
	"backend/services/agentClient"
	"backend/services/observer"
	"database/sql"
	"time"

	"fmt"
	"log"
)

// UserRepository struct for interacting with database
type AccountRepository struct{
	ObserverManager *observer.ObserverManager
	AgentClientRepository *agentClient.AgentClientRepository
}

// NewClientRepository initializes a new ClientRepository
func NewAccountRepository(observerManager *observer.ObserverManager, agentClientRepository *agentClient.AgentClientRepository) *AccountRepository {
	repo := &AccountRepository{ObserverManager: observerManager, AgentClientRepository: agentClientRepository}
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
	is_active BOOLEAN NOT NULL DEFAULT TRUE
);`

	_, err := database.DB.Exec(query)
	if err != nil {
		log.Fatal("❌ Error creating account table:", err)
	}

	fmt.Println("✅ Account tables checked/created!")
}


func (r *AccountRepository) CreateAccount(account models.Account) (models.Account, error) {

	// Check if opening_date is empty, if so, set it to today's date
	if account.OpeningDate == "" {
		account.OpeningDate = time.Now().Format("2006-01-02") // Format the date as yyyy-mm-dd
	}

	// SQL query to insert a new account while omitting 'is_active' and using today's date for 'opening_date'
	query := `
	INSERT INTO account 
	(client_id, account_type, account_status, opening_date, initial_deposit, currency, branch_id) 
	VALUES (?, ?, ?, ?, ?, ?, ?)` 
	
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
	account.IsActive = true

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

	// Proceed with updating the is_active field to false (soft delete)
	query := `UPDATE account SET is_active = FALSE WHERE account_id = ?`

	_, err = database.DB.Exec(query, accountID)
	if err != nil {
		return fmt.Errorf("failed to soft delete account: %v", err)
	}
	
	return nil
}

// GetAccountByID retrieves an account by accountID
func (r *AccountRepository) GetAccountByID(account_id int) (models.Account, error) {
	// Query updated to fetch only active accounts
	query := `SELECT * FROM account WHERE account_id = ? AND is_active = TRUE`

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
		&account.IsActive,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Account{}, fmt.Errorf("account with ID %d does not exist", account_id)
		}
		return models.Account{}, fmt.Errorf("failed to retrieve account: %v", err)
	}

	return account, nil
}

func (r *AccountRepository) GetAccountByClientId(client_id string) ([]models.Account, error) {
	// Query updated to fetch only active accounts
	query := `SELECT * FROM account WHERE client_id = ? AND is_active = TRUE`

	// Prepare a slice to store all the accounts
    var accounts []models.Account

    // Execute the query to fetch all accounts for the client
    rows, err := database.DB.Query(query, client_id)
    if err != nil {
        return nil, fmt.Errorf("failed to execute query: %v", err)
    }
    defer rows.Close()

    // Iterate through the rows and scan each one into the accounts slice
    for rows.Next() {
        var account models.Account
        if err := rows.Scan(
            &account.AccountID,
            &account.ClientID,
            &account.AccountType,
            &account.AccountStatus,
            &account.OpeningDate,
            &account.InitialDeposit,
            &account.Currency,
            &account.BranchID,
            &account.IsActive,
        ); err != nil {
            return nil, fmt.Errorf("failed to scan account: %v", err)
        }
        accounts = append(accounts, account)
    }

    // Check for any error encountered during iteration
    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("error encountered during rows iteration: %v", err)
    }

    // If no accounts were found, return an empty slice with no error
    if len(accounts) == 0 {
        return nil, nil // Return nil to indicate no accounts were found
    }

	return accounts, nil
}

// GetAllAccounts retrieves all active accounts from the database
func (r *AccountRepository) GetAllAccounts() ([]models.Account, error) {
	query := `SELECT * FROM account WHERE is_active = TRUE`
	
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve accounts: %v", err)
	}
	defer rows.Close()
	
	var accounts []models.Account
	for rows.Next() {
		var account models.Account
		err := rows.Scan(
			&account.AccountID,
			&account.ClientID,
			&account.AccountType,
			&account.AccountStatus,
			&account.OpeningDate,
			&account.InitialDeposit,
			&account.Currency,
			&account.BranchID,
			&account.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning account row: %v", err)
		}
		accounts = append(accounts, account)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating account rows: %v", err)
	}
	
	return accounts, nil
}

// GetAccountsByClientID retrieves all active accounts for a specific client
func (r *AccountRepository) GetAccountsByClientID(clientID string) ([]models.Account, error) {
	query := `SELECT * FROM account WHERE client_id = ? AND is_active = TRUE`
	
	rows, err := database.DB.Query(query, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve accounts for client %s: %v", clientID, err)
	}
	defer rows.Close()
	
	var accounts []models.Account
	for rows.Next() {
		var account models.Account
		err := rows.Scan(
			&account.AccountID,
			&account.ClientID,
			&account.AccountType,
			&account.AccountStatus,
			&account.OpeningDate,
			&account.InitialDeposit,
			&account.Currency,
			&account.BranchID,
			&account.IsActive,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning account row: %v", err)
		}
		accounts = append(accounts, account)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating account rows: %v", err)
	}
	
	return accounts, nil
}