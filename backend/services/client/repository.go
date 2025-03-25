package client

import (
	"backend/database"
	"database/sql"
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
	query := `
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
	_, err := database.DB.Exec(query)
	if err != nil {
		log.Fatal("❌ Error creating client table:", err)
	}
}

// CreateAccount inserts a new account into the database
func (r *ClientRepository) CreateClient(client Client, AgentID int) (Client, error) {

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
	
	// ✅ Insert into agent_client with agent_id
	agentClientQuery := `
	INSERT INTO agent_client 
	(client_id, id) 
	VALUES (?, ?)`
	
	_, err = database.DB.Exec(agentClientQuery, client.ClientID, AgentID)
	if err != nil {
		return Client{}, fmt.Errorf("failed to insert into agent_client: %v", err)
	}

	return client, nil
}

func (r *ClientRepository) AgentExists(AgentID int) (bool, error) {
	query := `SELECT 1 FROM users WHERE id = ?`
	// check with agent exisit 

	// to-do: need to check if user is agent or admin 

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
