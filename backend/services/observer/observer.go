package observer

import (
	"backend/models"
	"backend/services/account" // Import the account package
	"backend/services/agentclient_logs"
	"fmt"
	// Import the client package
)

// LogObserver is the observer interface for logging actions
type LogObserver interface {
	NotifyCreate(agentID int, clientID string, object interface{})
	NotifyUpdate(agentID int, clientID string, before, after interface{})
	NotifyDelete(agentID int, clientID string, object interface{})
}

// ClientObserver listens for client-related changes and notifies the AgentClientLogService
type ClientObserver struct {
	LogService *agentclient_logs.AgentClientLogService
}

func (co *ClientObserver) NotifyCreate(agentID int, clientID string, object interface{}) {
	fmt.Println("Observer: Notifying client creation for client ID:", clientID)
	client := object.(*models.Client) // Use client type from the client package
	// Call the AgentClientLogService to create the log for client creation
	co.LogService.LogAgentClientAction(agentID, clientID, "Create", map[string]interface{}{"details": client})
}

func (co *ClientObserver) NotifyUpdate(agentID int, clientID string, before, after interface{}) {
	fmt.Println("Observer: Notifying client update for client ID:", clientID)
	beforeClient := before.(*models.Client) // Use client type from the client package
	afterClient := after.(*models.Client)   // Use client type from the client package
	// Prepare the modified fields (before and after comparison)
	changes := Compare(beforeClient, afterClient)
	// Call the AgentClientLogService to create the log for client update
	co.LogService.LogAgentClientAction(agentID, clientID, "Update", map[string]interface{}{"details": changes})
}

func (co *ClientObserver) NotifyDelete(agentID int, clientID string, object interface{}) {
	fmt.Println("Observer: Notifying client delete for client ID:", clientID)
	client := object.(*models.Client) // Use client type from the client package
	// Call the AgentClientLogService to create the log for client deletion
	co.LogService.LogAgentClientAction(agentID, clientID, "Delete", map[string]interface{}{"details": client})
}

// AccountObserver listens for account-related changes and notifies the AgentClientLogService
type AccountObserver struct {
	LogService *agentclient_logs.AgentClientLogService
}

func (ao *AccountObserver) NotifyCreate(agentID int, clientID string, object interface{}) {
	fmt.Println("Observer: Notifying account creation for client ID:", clientID)
	account := object.(*account.Account) // Use account type from the account package
	// Call the AgentClientLogService to create the log for account creation
	ao.LogService.LogAccountChange(agentID, clientID, "Create", map[string]interface{}{"details": account})
}

func (ao *AccountObserver) NotifyUpdate(agentID int, clientID string, before, after interface{}) {
	fmt.Println("Observer: Notifying account update for client ID:", clientID)
	beforeAccount := before.(*account.Account) // Use account type from the account package
	afterAccount := after.(*account.Account)   // Use account type from the account package
	// Prepare the modified fields (before and after comparison)
	changes := Compare(beforeAccount, afterAccount)
	// Call the AgentClientLogService to create the log for account update
	ao.LogService.LogAccountChange(agentID, clientID, "Update", map[string]interface{}{"details": changes})
}

func (ao *AccountObserver) NotifyDelete(agentID int, clientID string, object interface{}) {
	fmt.Println("Observer: Notifying account delete for client ID:", clientID)
	account := object.(*account.Account) // Use account type from the account package
	// Call the AgentClientLogService to create the log for account deletion
	ao.LogService.LogAccountChange(agentID, clientID, "Delete", map[string]interface{}{"details": account})
}
