package client

import (
	"backend/database"
	"database/sql"
	"fmt"
	"log"
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
	
	// ✅ Insert into agent_client with id set to NULL 
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