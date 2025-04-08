package agentClient

import (
	"backend/models"
	"fmt"
	"sort"
)

// UserService struct to interact with the repository layer
type AgentClientService struct {
	repo *AgentClientRepository
}

// NewUserService initializes the user service
func NewAgentClientService(repo *AgentClientRepository) *AgentClientService {
	return &AgentClientService{repo: repo}
}

// Add GetAccountByClientId method to implement the interface
func (s *AgentClientService) GetUnassignedClients() ([]models.AgentClient, error) {
	return s.repo.GetUnassignedClients()
}

// // CreateUser processes user creation request
// func (s *AgentClientService) UpdateAgentToClient(clientID string, newID int) (AgentClient, error) {
// 	// ✅ Check if client exists
// 	exists, err := s.repo.ClientExists(clientID)
// 	if err != nil {
// 		return AgentClient{}, fmt.Errorf("failed to check client existence: %v", err)
// 	}

// 	if !exists {
// 		return AgentClient{}, fmt.Errorf("client_id not found")
// 	}

// 	// ✅ Update agent ID
// 	err = s.repo.UpdateAgentToClient(clientID, newID)
// 	if err != nil {
// 		return AgentClient{}, fmt.Errorf("failed to update agent ID: %v", err)
// 	}

// 	// ✅ Return the updated client
// 	return AgentClient{
// 		ClientID: clientID,
// 		AgentID:  newID,
// 	}, nil
// }

func (s *AgentClientService) AssignAgentsToUnassignedClients() error {
	// Get all unassigned clients
	unassignedClients, err := s.repo.GetUnassignedClients()
	if err != nil {
		return fmt.Errorf("failed to get unassigned clients: %v", err)
	}

	if len(unassignedClients) == 0 {
		return nil // No unassigned clients, nothing to do
	}

	// Get all agents from users table
	agents, err := s.repo.GetAllAgents()
	if err != nil {
		return fmt.Errorf("failed to get agents: %v", err)
	}

	if len(agents) == 0 {
		return fmt.Errorf("no agents available to assign clients")
	}

	// Get current agent client counts
	agentClientCounts, err := s.repo.GetAgentClientCount()
	if err != nil {
		return fmt.Errorf("failed to get agent client counts: %v", err)
	}

	// Create agent workload data structure
	type AgentWorkload struct {
		Agent     models.Agent
		ClientNum int
	}

	// Initialize workloads for all agents (including those without clients)
	var agentWorkloads []AgentWorkload
	for _, agent := range agents {
		// Get count from map, which will be 0 if agent doesn't exist in the map
		clientCount := agentClientCounts[agent.ID]
		agentWorkloads = append(agentWorkloads, AgentWorkload{
			Agent:     agent,
			ClientNum: clientCount,
		})
	}

	// Process each unassigned client
	for _, client := range unassignedClients {
		// Sort agents by client count (ascending) and then by ID (ascending)
		sort.SliceStable(agentWorkloads, func(i, j int) bool {
			if agentWorkloads[i].ClientNum == agentWorkloads[j].ClientNum {
				return agentWorkloads[i].Agent.ID < agentWorkloads[j].Agent.ID
			}
			return agentWorkloads[i].ClientNum < agentWorkloads[j].ClientNum
		})

		// Assign to agent with lowest workload
		selectedAgent := agentWorkloads[0].Agent
		err = s.repo.UpdateAgentToClient(client.ClientID, selectedAgent.ID)
		if err != nil {
			return fmt.Errorf("failed to assign client %s to agent %d: %v", 
				client.ClientID, selectedAgent.ID, err)
		}

		// Update workload count
		agentWorkloads[0].ClientNum++
	}

	return nil
}

func (s *AgentClientService) GetAgentIDByClientID(clientID string) (int, error) {
    if s.repo == nil {
        return 0, fmt.Errorf("repository is not initialized")
    }
    return s.repo.GetAgentIDByClientID(clientID)
}

