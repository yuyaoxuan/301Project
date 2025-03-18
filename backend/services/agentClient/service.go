package agentClient

import (
	"fmt"
)

// UserService struct to interact with the repository layer
type AgentClientService struct {
	repo *AgentClientRepository
}

// NewUserService initializes the user service
func NewAgentClientService(repo *AgentClientRepository) *AgentClientService {
	return &AgentClientService{repo: repo}
}

// CreateUser processes user creation request
func (s *AgentClientService) UpdateAgentToClient(clientID string, newID int) (Client, error) {
	// ✅ Check if client exists
	exists, err := s.repo.ClientExists(clientID)
	if err != nil {
		return Client{}, fmt.Errorf("failed to check client existence: %v", err)
	}

	if !exists {
		return Client{}, fmt.Errorf("client_id not found")
	}

	// ✅ Update agent ID
	err = s.repo.UpdateAgentToClient(clientID, newID)
	if err != nil {
		return Client{}, fmt.Errorf("failed to update agent ID: %v", err)
	}

	// ✅ Return the updated client
	return Client{
		ClientID: clientID,
		AgentID:  newID,
	}, nil
}

// func (s *AgentClientService) AssignAgentsToUnassignedClients() error {
// 	// ✅ Get all unassigned clients
// 	unassignedClients, err := s.repo.GetUnassignedClients()
// 	if err != nil {
// 		return fmt.Errorf("failed to get unassigned clients: %v", err)
// 	}

// 	if len(unassignedClients) == 0 {
// 		return nil // No unassigned clients, nothing to do
// 	}

// 	// ✅ Get all agents
// 	agents, err := s.repo.GetAllAgents()
// 	if err != nil {
// 		return fmt.Errorf("failed to get agents: %v", err)
// 	}

// 	if len(agents) == 0 {
// 		return fmt.Errorf("no agents available to assign clients")
// 	}

// 	// ✅ Assign clients evenly among agents
// 	clientCount := len(unassignedClients)
// 	agentCount := len(agents)
// 	clientsPerAgent := clientCount / agentCount
// 	extraClients := clientCount % agentCount

// 	index := 0
// 	for i, agent := range agents {
// 		// First agent gets the extra client if division is uneven
// 		numToAssign := clientsPerAgent
// 		if i < extraClients {
// 			numToAssign++
// 		}

// 		for j := 0; j < numToAssign && index < clientCount; j++ {
// 			err = s.repo.UpdateAgentID(unassignedClients[index].ClientID, agent.ID)
// 			if err != nil {
// 				return fmt.Errorf("failed to assign client %s to agent %d: %v", unassignedClients[index].ClientID, agent.ID, err)
// 			}
// 			index++
// 		}
// 	}

// 	return nil
// }



