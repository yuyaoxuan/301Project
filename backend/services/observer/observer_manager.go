package observer

import (
	"backend/services/account" // Import the account package
	"backend/services/client"  // Import the client package
)

// ObserverManager manages the observers for client and account actions
type ObserverManager struct {
	clientObservers  []LogObserver
	accountObservers []LogObserver
}

// AddClientObserver adds a client observer to the manager
func (om *ObserverManager) AddClientObserver(observer LogObserver) {
	om.clientObservers = append(om.clientObservers, observer)
}

// AddAccountObserver adds an account observer to the manager
func (om *ObserverManager) AddAccountObserver(observer LogObserver) {
	om.accountObservers = append(om.accountObservers, observer)
}

// NotifyClientCreate notifies all client observers to create a log
func (om *ObserverManager) NotifyClientCreate(agentID int, clientID string, client *client.Client) {
	for _, observer := range om.clientObservers {
		observer.NotifyCreate(agentID, clientID, client)
	}
}

// NotifyClientUpdate notifies all client observers to update a log
func (om *ObserverManager) NotifyClientUpdate(agentID int, clientID string, before, after *client.Client) {
	for _, observer := range om.clientObservers {
		observer.NotifyUpdate(agentID, clientID, before, after)
	}
}

// NotifyClientDelete notifies all client observers to delete a log
func (om *ObserverManager) NotifyClientDelete(agentID int, clientID string, client *client.Client) {
	for _, observer := range om.clientObservers {
		observer.NotifyDelete(agentID, clientID, client)
	}
}

// NotifyAccountCreate notifies all account observers to create a log
func (om *ObserverManager) NotifyAccountCreate(agentID int, clientID string, account *account.Account) {
	for _, observer := range om.accountObservers {
		observer.NotifyCreate(agentID, clientID, account)
	}
}

// NotifyAccountUpdate notifies all account observers to update a log
func (om *ObserverManager) NotifyAccountUpdate(agentID int, clientID string, before, after *account.Account) {
	for _, observer := range om.accountObservers {
		observer.NotifyUpdate(agentID, clientID, before, after)
	}
}

// NotifyAccountDelete notifies all account observers to delete a log
func (om *ObserverManager) NotifyAccountDelete(agentID int, clientID string, account *account.Account) {
	for _, observer := range om.accountObservers {
		observer.NotifyDelete(agentID, clientID, account)
	}
}
