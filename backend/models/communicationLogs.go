package models

// CommunicationLog represents an email notification event
type CommunicationLog struct {
	ID           int    `json:"id"`
	LogID        int    `json:"log_id"`
	ClientID     string `json:"client_id"`
	AgentID      int    `json:"agent_id"`
	EmailSubject string `json:"email_subject"`
	EmailStatus  string `json:"email_status"`
	Timestamp    string `json:"timestamp"`
}
