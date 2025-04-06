package client

import (
	"backend/database"
	"backend/services/observer"
	"database/sql"
	"fmt"
	"log"

	"backend/models"
)

// // Client struct represents a client in the system
// type Client struct {
// 	ClientID           string `json:"client_id"`
// 	FirstName          string `json:"first_name"`
// 	LastName           string `json:"last_name"`
// 	DOB                string `json:"dob"`
// 	Gender             string `json:"gender"`
// 	Email              string `json:"email"`
// 	Phone              string `json:"phone"`
// 	Address            string `json:"address"`
// 	City               string `json:"city"`
// 	State              string `json:"state"`
// 	Country            string `json:"country"`
// 	PostalCode         string `json:"postal_code"`
// 	VerificationStatus string `json:"verification_status"`
// }

// ClientRepository struct for interacting with the database
type ClientRepository struct {
	ObserverManager *observer.ObserverManager
}

// NewClientRepository initializes a new ClientRepository
func NewClientRepository(observerManager *observer.ObserverManager) *ClientRepository {
	repo := &ClientRepository{ObserverManager: observerManager}
	repo.InitClientTables() // Ensure tables exist when the repository is created
	return repo
}


// InitClientTables creates the client table if it doesn't exist
func (r *ClientRepository) InitClientTables() {
	// Create client table
	query := `
	CREATE TABLE IF NOT EXISTS client (
		client_id VARCHAR(50) PRIMARY KEY,
		first_name CHAR(50) NOT NULL,
		last_name CHAR(50) NOT NULL,
		dob DATE NOT NULL,
		gender VARCHAR(20) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		phone VARCHAR(15) UNIQUE NOT NULL,
		address VARCHAR(100) NOT NULL,
		city VARCHAR(50) NOT NULL,
		state VARCHAR(50) NOT NULL,
		country VARCHAR(50) NOT NULL,
		postal_code VARCHAR(10) NOT NULL,
		verification_status VARCHAR(20) DEFAULT 'unverified'
	);`
	_, err := database.DB.Exec(query)
	if err != nil {
		log.Fatal("❌ Error creating client table:", err)
	}

	counterTable := `
	CREATE TABLE IF NOT EXISTS counter (
		id INT PRIMARY KEY AUTO_INCREMENT,
		name VARCHAR(50) UNIQUE NOT NULL,
		value INT NOT NULL
	);`
	_, err = database.DB.Exec(counterTable)
	if err != nil {
		log.Fatal("❌ Error creating counter table:", err)
	}

	// Initialize client counter if not exists
	_, err = database.DB.Exec("INSERT IGNORE INTO counter (name, value) VALUES ('client', 0);")
	if err != nil {
		log.Fatal("❌ Error initializing client counter:", err)
	}

	fmt.Println("✅ Client table checked/created!")
}

// EmailExists checks if an email already exists in the database
func (r *ClientRepository) EmailExists(email string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM client WHERE email = ?`
	err := database.DB.QueryRow(query, email).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check email uniqueness: %v", err)
	}
	return count > 0, nil
}

// PhoneExists checks if a phone number already exists in the database
func (r *ClientRepository) PhoneExists(phone string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM client WHERE phone = ?`
	err := database.DB.QueryRow(query, phone).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check phone uniqueness: %v", err)
	}
	return count > 0, nil
}

// CreateClient inserts a new client into the database
func (r *ClientRepository) CreateClient(client models.Client, AgentID int) (models.Client, error) {
	var currentValue int

	// Begin a transaction to ensure atomicity
	tx, err := database.DB.Begin()
	if err != nil {
		return models.Client{}, fmt.Errorf("failed to begin transaction: %v", err)
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Retrieve the current value of the client counter
	query := `SELECT value FROM counter WHERE name = 'client' FOR UPDATE`
	err = tx.QueryRow(query).Scan(&currentValue)
	if err != nil {
		return models.Client{}, fmt.Errorf("failed to get client counter: %v", err)
	}

	// Increment the counter and generate a new client ID
	newValue := currentValue + 1
	client.ClientID = fmt.Sprintf("client%d", newValue)

	// Update the counter value in the database
	updateQuery := `UPDATE counter SET value = ? WHERE name = 'client'`
	_, err = tx.Exec(updateQuery, newValue)
	if err != nil {
		return models.Client{}, fmt.Errorf("failed to update client counter: %v", err)
	}

	// Set default verification status
	client.VerificationStatus = "unverified"

	// Insert the new client into the database
	insertQuery := `
        INSERT INTO client 
        (client_id, first_name, last_name, dob, gender, email, phone, address, city, state, country, postal_code, verification_status)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err = tx.Exec(insertQuery,
		client.ClientID, client.FirstName, client.LastName, client.DOB,
		client.Gender, client.Email, client.Phone, client.Address,
		client.City, client.State, client.Country, client.PostalCode,
		client.VerificationStatus,
	)
	if err != nil {
		return models.Client{}, fmt.Errorf("failed to insert client: %v", err)
	}

	// Commit the transaction
	if err = tx.Commit(); err != nil {
		return models.Client{}, fmt.Errorf("failed to commit transaction: %v", err)
	}

	// ✅ Insert into agent_client with agent_id
	agentClientQuery := `
	INSERT INTO agent_client 
	(client_id, id) 
	VALUES (?, ?)`

	_, err = database.DB.Exec(agentClientQuery, client.ClientID, AgentID)
	if err != nil {
		return models.Client{}, fmt.Errorf("failed to insert into agent_client: %v", err)
	}

	return client, nil
}

func (r *ClientRepository) AgentExists(AgentID int) (bool, error) {
	query := `SELECT 1 FROM users WHERE id = ? AND role = 'agent'`
	// check with agent exisit
	var exists int
	err := database.DB.QueryRow(query, AgentID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetClientByID retrieves a client by their ID
func (r *ClientRepository) GetClientByID(clientID string) (models.Client, error) {
	query := `SELECT * FROM client WHERE client_id = ?`

	var client models.Client
	err := database.DB.QueryRow(query, clientID).Scan(
		&client.ClientID, &client.FirstName, &client.LastName,
		&client.DOB, &client.Gender, &client.Email,
		&client.Phone, &client.Address, &client.City,
		&client.State, &client.Country, &client.PostalCode,
		&client.VerificationStatus,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return models.Client{}, fmt.Errorf("client with ID %s not found", clientID)
		}
		return models.Client{}, fmt.Errorf("failed to retrieve client: %v", err)
	}

	return client, nil
}

// UpdateClient updates an existing client's information
func (r *ClientRepository) UpdateClient(client models.Client) (models.Client, error) {
	// Check if client exists
	_, err := r.GetClientByID(client.ClientID)
	if err != nil {
		return models.Client{}, err
	}

	// Check email uniqueness if changed
	var currentEmail string
	err = database.DB.QueryRow("SELECT email FROM client WHERE client_id = ?", client.ClientID).Scan(&currentEmail)
	if err != nil {
		return models.Client{}, fmt.Errorf("failed to retrieve current email: %v", err)
	}

	if currentEmail != client.Email {
		exists, err := r.EmailExists(client.Email)
		if err != nil {
			return models.Client{}, err
		}
		if exists {
			return models.Client{}, fmt.Errorf("email address already exists")
		}
	}

	// Check phone uniqueness if changed
	var currentPhone string
	err = database.DB.QueryRow("SELECT phone FROM client WHERE client_id = ?", client.ClientID).Scan(&currentPhone)
	if err != nil {
		return models.Client{}, fmt.Errorf("failed to retrieve current phone: %v", err)
	}

	if currentPhone != client.Phone {
		exists, err := r.PhoneExists(client.Phone)
		if err != nil {
			return models.Client{}, err
		}
		if exists {
			return models.Client{}, fmt.Errorf("phone number already exists")
		}
	}

	query := `
    UPDATE client 
    SET first_name = ?, last_name = ?, dob = ?, gender = ?,
        email = ?, phone = ?, address = ?, city = ?,
        state = ?, country = ?, postal_code = ?
    WHERE client_id = ?`

	_, err = database.DB.Exec(query,
		client.FirstName, client.LastName, client.DOB, client.Gender,
		client.Email, client.Phone, client.Address, client.City,
		client.State, client.Country, client.PostalCode, client.ClientID,
	)

	if err != nil {
		return models.Client{}, fmt.Errorf("failed to update client: %v", err)
	}

	// Retrieve the updated client to return
	return r.GetClientByID(client.ClientID)
}

// DeleteClient removes a client's profile from the database
func (r *ClientRepository) DeleteClient(clientID string) error {
	query := `DELETE FROM client WHERE client_id = ?`

	result, err := database.DB.Exec(query, clientID)
	if err != nil {
		return fmt.Errorf("failed to delete client: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected; client with ID %s not found", clientID)
	}

	return nil
}

// VerifyClient updates a client's verification status
func (r *ClientRepository) VerifyClient(clientID string) error {
	// Check if client exists
	_, err := r.GetClientByID(clientID)
	if err != nil {
		return err
	}

	query := `
	SELECT verification_status FROM client WHERE client_id = ?`

	var currentStatus string
	err = database.DB.QueryRow(query, clientID).Scan(&currentStatus)
	if err != nil {
		return fmt.Errorf("failed to retrieve verification status: %v", err)
	}

	if currentStatus == "verified" {
		return fmt.Errorf("client %s is already verified", clientID)
	}

	updateQuery := `
    UPDATE client 
    SET verification_status = 'verified' 
    WHERE client_id = ?`

	_, err = database.DB.Exec(updateQuery, clientID)
	if err != nil {
		return fmt.Errorf("failed to update verification status: %v", err)
	}

	return nil
}
