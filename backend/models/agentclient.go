package models

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