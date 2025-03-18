package agentClient

import (
	"backend/database"
	"database/sql"
	"log"
)

// User struct represents a user in the system
type Client struct {
	ClientID string `json:"client_id"`
	AgentID int `json:"id"`
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
	// Create client table
	query := `
	CREATE TABLE IF NOT EXISTS agent_client (
		client_id VARCHAR(50) NOT NULL,
		id INT
	);`

	// todo: agent need to be assigned based on the agent that created the client 

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

func (r *AgentClientRepository) UpdateAgentToClient(clientID string, newID int) error {
	query := `UPDATE agent_client SET id = ? WHERE client_id = ?`
	_, err := database.DB.Exec(query, newID, clientID)
	return err
}

// func (r *AgentClientRepository) GetUnassignedClients() ([]unassignClients, error) {
// 	query := `SELECT client_id FROM agent_client WHERE id IS NULL`
// 	rows, err := database.DB.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var clients []unassignClients
// 	for rows.Next() {
// 		var client unassignClients
// 		err := rows.Scan(&client.ClientID)
// 		if err != nil {
// 			return nil, err
// 		}
// 		clients = append(clients, client)
// 	}

// 	return clients, nil
// }

// func (r *AgentClientRepository) GetAllAgents() ([]Agent, error) {
// 	query := `SELECT id, first_name, last_name, email, role FROM agent`
// 	rows, err := database.DB.Query(query)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer rows.Close()

// 	var agents []Agent
// 	for rows.Next() {
// 		var agent Agent
// 		err := rows.Scan(&agent.ID, &agent.FirstName, &agent.LastName, &agent.Email, &agent.Role)
// 		if err != nil {
// 			return nil, err
// 		}
// 		agents = append(agents, agent)
// 	}

// 	return agents, nil
// }

// func (r *AgentClientRepository) UpdateAgentID(clientID string, newID int) error {
// 	query := `UPDATE agent_client SET id = ? WHERE client_id = ?`
// 	_, err := database.DB.Exec(query, newID, clientID)
// 	return err
// }