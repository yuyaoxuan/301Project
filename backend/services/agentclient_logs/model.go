package agentclient_logs

// AgentClientLog represents a log entry for agent-client operations
type AgentClientLog struct {
	ID             int                    `json:"id"`
	AgentID        int                    `json:"agent_id"`  // agent_id as an integer
	ClientID       string                 `json:"client_id"` // client_id as a string
	Action         string                 `json:"action"`
	ModifiedFields map[string]interface{} `json:"modified_fields"`
	Timestamp      string                 `json:"timestamp"`
}
