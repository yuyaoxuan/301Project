package observer

import (
	"backend/models"
	"fmt"
)

// ObserverManager manages the observers for client and account actions
type ObserverManager struct {
	clientObservers        []LogObserver
	accountObservers       []LogObserver
	communicationObservers []LogObserver
}

// AddClientObserver adds a client observer to the manager
func (om *ObserverManager) AddClientObserver(observer LogObserver) {
	om.clientObservers = append(om.clientObservers, observer)
	fmt.Println("Number of client observers after registration:", len(om.clientObservers))
}

// AddAccountObserver adds an account observer to the manager
func (om *ObserverManager) AddAccountObserver(observer LogObserver) {
	om.accountObservers = append(om.accountObservers, observer)
	fmt.Println("Number of account observers after registration:", len(om.accountObservers))
}

// AddCommunicationObserver adds a communication observer to the manager
func (om *ObserverManager) AddCommunicationObserver(observer LogObserver) {
	om.communicationObservers = append(om.communicationObservers, observer)
	fmt.Println("Number of communication observers after registration:", len(om.communicationObservers))
}

// NotifyClientCreate notifies all client observers to create a log
func (om *ObserverManager) NotifyClientCreate(agentID int, clientID string, client *models.Client) {
	fmt.Println("ObserverManager: Notifying client creation for client ID:", client.ClientID)
	for _, observer := range om.clientObservers {
		observer.NotifyCreate(agentID, clientID, client)
	}
}

// NotifyClientUpdate notifies all client observers to update a log
func (om *ObserverManager) NotifyClientUpdate(agentID int, clientID string, before, after *models.Client) {
	fmt.Println("ObserverManager: Notifying client update for client ID:", clientID)
	for _, observer := range om.clientObservers {
		observer.NotifyUpdate(agentID, clientID, before, after)
	}
}

// NotifyClientDelete notifies all client observers to delete a log
func (om *ObserverManager) NotifyClientDelete(agentID int, clientID string, client *models.Client) {
	fmt.Println("ObserverManager: Notifying client delete for client ID:", clientID)
	for _, observer := range om.clientObservers {
		observer.NotifyDelete(agentID, clientID, client)
	}
}

// NotifyAccountCreate notifies all account observers to create a log
func (om *ObserverManager) NotifyAccountCreate(agentID int, clientID string, account *models.Account) {
	fmt.Println("ObserverManager: Notifying account creation for client ID:", clientID)
	for _, observer := range om.accountObservers {
		observer.NotifyCreate(agentID, clientID, account)
	}
}

// NotifyAccountUpdate notifies all account observers to update a log
func (om *ObserverManager) NotifyAccountUpdate(agentID int, clientID string, before, after *models.Account) {
	fmt.Println("ObserverManager: Notifying account update for client ID:", clientID)
	for _, observer := range om.accountObservers {
		observer.NotifyUpdate(agentID, clientID, before, after)
	}
}

// NotifyAccountDelete notifies all account observers to delete a log
func (om *ObserverManager) NotifyAccountDelete(agentID int, clientID string, account *models.Account) {
	fmt.Println("ObserverManager: Notifying account delete for client ID:", clientID)
	for _, observer := range om.accountObservers {
		observer.NotifyDelete(agentID, clientID, account)
	}
}

// NotifyCommunication logs the communication via the observer
func (om *ObserverManager) NotifyCommunication(agentID int, clientID string, log models.AgentClientLog) {
	// Notify all communication observers
	fmt.Println("ObserverManager: Notifying communication for client ID:", clientID)
	for _, observer := range om.communicationObservers {
		observer.NotifyCreate(agentID, clientID, log)
	}
}
