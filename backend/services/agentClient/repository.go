package agentClient

import (
	"backend/database"
	"database/sql"
	"fmt"
	"log"
)

// User struct represents a user in the system
type AgentClient struct {
	ClientID string `json:"client_id"`
	AgentID int `json:"id"`
}

type Agent struct {
    ID        int    `json:"id"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
    Email     string `json:"email"`
    Role      string `json:"role"`
}

// UserRepository struct for interacting with database
type AgentClientRepository struct{}

// NewClientRepository initializes a new ClientRepository
func NewAgentClientRepository() *AgentClientRepository {
	repo := &AgentClientRepository{}
	repo.InitAgentClientTables() // ✅ Ensure tables exist when the repository is created
	return repo
}

// InitClientTables creates the client and account tables if they don't exist
func (r *AgentClientRepository) InitAgentClientTables() {
	// if client is deleted, delete value. if agent is deleted, set id to null 
	query := `
	CREATE TABLE IF NOT EXISTS agent_client (
		client_id VARCHAR(50) NOT NULL,
		id INT,
		FOREIGN KEY (client_id) REFERENCES client(client_id) ON DELETE CASCADE ON UPDATE CASCADE,
		FOREIGN KEY (id) REFERENCES users(id) ON DELETE SET NULL
	);`

	_, err := database.DB.Exec(query)
	if err != nil {	
		log.Fatal("❌ Error creating agent client table:", err)
	}
}	

func (r *AgentClientRepository) ClientExists(clientID string) (bool, error) {
	query := `SELECT 1 FROM agent_client WHERE client_id = ?`
	var exists int
	err := database.DB.QueryRow(query, clientID).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// IsAgentNull checks if the agent ID for a given client is NULL
func (r *AgentClientRepository) IsAgentNull(clientID string) (bool, error) {
	query := `SELECT id FROM agent_client WHERE client_id = ?`
	var agentID *int

	err := database.DB.QueryRow(query, clientID).Scan(&agentID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, fmt.Errorf("client with ID %s not found", clientID)
		}
		return false, err
	}

	return agentID == nil, nil
}

func (r *AgentClientRepository) UpdateAgentToClient(clientID string, newID int) error {
	// Check if the agent ID is NULL
	isNull,user_id_err := r.IsAgentNull(clientID)
	if user_id_err != nil  {
		return user_id_err
	}
	if !isNull {
		return fmt.Errorf("agent is already assigned to client %s", clientID)
	}

	query := `UPDATE agent_client SET id = ? WHERE client_id = ?`
	_, err := database.DB.Exec(query, newID, clientID)
	return err
}

func (r *AgentClientRepository) GetUnassignedClients() ([]AgentClient, error) {
	query := `SELECT client_id FROM agent_client WHERE id IS NULL`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var clients []AgentClient
	for rows.Next() {
		var client AgentClient
		err := rows.Scan(&client.ClientID)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	return clients, nil
}

func (r *AgentClientRepository) GetAllAgents() ([]Agent, error) {
	query := `SELECT id, first_name, last_name, email, role FROM users WHERE role = 'Agent'`
	rows, err := database.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var agents []Agent
	for rows.Next() {
		var agent Agent
		err := rows.Scan(&agent.ID, &agent.FirstName, &agent.LastName, &agent.Email, &agent.Role)
		if err != nil {
			return nil, err
		}
		agents = append(agents, agent)
	}

	return agents, nil
}


func (r *AgentClientRepository) IsAgent(userID int) (bool, error) {
	var role string
	query := `SELECT role FROM users WHERE id = ?`
	err := database.DB.QueryRow(query, userID).Scan(&role)
	if err != nil {
		return false, fmt.Errorf("failed to retrieve user role: %v", err)
	}

	return role == "Agent", nil
}