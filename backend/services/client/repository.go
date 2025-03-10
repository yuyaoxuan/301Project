package client

import (
	"backend/database"
	"fmt"
	"log"
	"github.com/google/uuid"
)

// User struct represents a user in the system
type Client struct {
	ClientID   string `json:"client_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	DOB        string `json:"dob"`
	Gender     string `json:"gender"`
	Email      string `json:"email"`
	Phone      string `json:"phone"`
	Address    string `json:"address"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
	
}


type Account struct {
	AccountID     int     `json:"account_id"`
	ClientID      string  `json:"client_id"`
	AccountType   string  `json:"account_type"`
	AccountStatus string  `json:"account_status"`
	OpeningDate   string  `json:"opening_date"`
	InitialDeposit float64 `json:"initial_deposit"`
	Currency      string  `json:"currency"`
	BranchID      string  `json:"branch_id"`
}


// UserRepository struct for interacting with database
type ClientRepository struct{}

// NewClientRepository initializes a new ClientRepository
func NewClientRepository() *ClientRepository {
	repo := &ClientRepository{}
	repo.InitClientTables() // ✅ Ensure tables exist when the repository is created
	return repo
}

// InitClientTables creates the client and account tables if they don't exist
func (r *ClientRepository) InitClientTables() {
	// Create client table
	clientTable := `
	CREATE TABLE IF NOT EXISTS client (
		client_id VARCHAR(50) PRIMARY KEY,
		first_name CHAR(50) NOT NULL,
		last_name CHAR(50) NOT NULL,
		dob DATE NOT NULL,
		gender VARCHAR(10) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		phone VARCHAR(15) UNIQUE NOT NULL,
		address VARCHAR(100) NOT NULL,
		city VARCHAR(50) NOT NULL,
		state VARCHAR(50) NOT NULL,
		country VARCHAR(50) NOT NULL,
		postal_code VARCHAR(10) NOT NULL
	);`
	_, err := database.DB.Exec(clientTable)
	if err != nil {
		log.Fatal("❌ Error creating client table:", err)
	}

	// Create account table
	accountTable := `
	CREATE TABLE IF NOT EXISTS account (
		account_id INT AUTO_INCREMENT PRIMARY KEY,
		client_id VARCHAR(50) NOT NULL,
		account_type VARCHAR(50) NOT NULL,
		account_status VARCHAR(50) NOT NULL,
		opening_date DATE NOT NULL,
		initial_deposit FLOAT NOT NULL,
		currency VARCHAR(50) NOT NULL,
		branch_id VARCHAR(50) NOT NULL,
		FOREIGN KEY (client_id) REFERENCES client(client_id)
	);`
	_, err = database.DB.Exec(accountTable)
	if err != nil {
		log.Fatal("❌ Error creating account table:", err)
	}

	fmt.Println("✅ Client and account tables checked/created!")
}

// CreateAccount inserts a new account into the database
func (r *ClientRepository) CreateClient(client Client) (Client, error) {

	// Generate a unique client ID if one isn't provided
	if client.ClientID == "" {
		client.ClientID = uuid.New().String()
	}

	query := `
	INSERT INTO client 
	(client_id, first_name, last_name, dob, gender, email, phone, address, city, state, country, postal_code) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	
	_, err := database.DB.Exec(query,
		client.ClientID, client.FirstName, client.LastName, client.DOB, client.Gender,
		client.Email, client.Phone, client.Address, client.City, client.State,
		client.Country, client.PostalCode,
	)
	if err != nil {
		return Client{}, fmt.Errorf("failed to insert client: %v", err)
	}

	return client, nil
}

func (r *ClientRepository) CreateAccount(account Account) (Account, error) {
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
		return Account{}, fmt.Errorf("failed to insert account: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Account{}, fmt.Errorf("failed to retrieve inserted account ID: %v", err)
	}

	account.AccountID = int(id)
	return account, nil
}

// GetClientByID retrieves a client by their ID
func (r *ClientRepository) GetClientByID(clientID string) (Client, error) {
    query := `SELECT * FROM client WHERE client_id = ?`
    
    var client Client
    err := database.DB.QueryRow(query, clientID).Scan(
        &client.ClientID, &client.FirstName, &client.LastName, 
        &client.DOB, &client.Gender, &client.Email, 
        &client.Phone, &client.Address, &client.City, 
        &client.State, &client.Country, &client.PostalCode,
    )
    
    if err != nil {
        return Client{}, fmt.Errorf("failed to retrieve client: %v", err)
    }
    
    return client, nil
}

// UpdateClient updates an existing client's information
func (r *ClientRepository) UpdateClient(client Client) (Client, error) {
    query := `
    UPDATE client 
    SET first_name = ?, last_name = ?, dob = ?, gender = ?, 
        email = ?, phone = ?, address = ?, city = ?, 
        state = ?, country = ?, postal_code = ?
    WHERE client_id = ?`
    
    _, err := database.DB.Exec(query,
        client.FirstName, client.LastName, client.DOB, client.Gender,
        client.Email, client.Phone, client.Address, client.City, 
        client.State, client.Country, client.PostalCode, client.ClientID,
    )
    
    if err != nil {
        return Client{}, fmt.Errorf("failed to update client: %v", err)
    }
    
    return client, nil
}

// DeleteClient removes a client from the database
func (r *ClientRepository) DeleteClient(clientID string) error {
    query := `DELETE FROM client WHERE client_id = ?`
    
    result, err := database.DB.Exec(query, clientID)
    if err != nil {
        return fmt.Errorf("failed to delete client: %v", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("error checking rows affected: %v", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("client with ID %s not found", clientID)
    }
    
    return nil
}

// VerifyClient updates a client's verification status
func (r *ClientRepository) VerifyClient(clientID string) error {
    // In a real application, you would have a verification_status column
    // For this implementation, we'll just check if the client exists
    
    query := `SELECT client_id FROM client WHERE client_id = ?`
    
    var id string
    err := database.DB.QueryRow(query, clientID).Scan(&id)
    if err != nil {
        return fmt.Errorf("client with ID %s not found: %v", clientID, err)
    }
    
    // In a real application, you would update the verification status here
    // UPDATE client SET verification_status = 'verified' WHERE client_id = ?
    
    return nil
}
