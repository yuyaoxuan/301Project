package routes

import (
	"net/http"

	"backend/services/account" // import account service
	"backend/services/agentClient"
	"backend/services/client"
	"backend/services/user" // import user service

	"github.com/gorilla/mux"
)

// SetupRoutes initializes the router and returns it
func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is running!"))
	}).Methods("GET")

	// USER ROUTES
	r.HandleFunc("/users", user.CreateUserHandler).Methods("POST")

	// CLIENT ROUTES
	r.HandleFunc("/clients/{agent_id}", client.CreateClientHandler).Methods("POST")
	
	// ACCOUNT Routes
	r.HandleFunc("/accounts", account.CreateAccountHandler).Methods("POST")
	r.HandleFunc("/accounts/{account_id}", account.DeleteAccountHandler).Methods("DELETE")
	
	// agentClient routes 
	r.HandleFunc("/agentClient/{client_id}", agentClient.UpdateAgentToClientHandler).Methods("PUT") 
	r.HandleFunc("/agentClient", agentClient.AssignAgentsToUnassignedClientsHandler).Methods("PUT") 
	r.HandleFunc("/api/clients", client.CreateClientHandler).Methods("POST")
    r.HandleFunc("/api/clients/{clientId}", client.GetClientHandler).Methods("GET")
    r.HandleFunc("/api/clients/{clientId}", client.UpdateClientHandler).Methods("PUT")
    r.HandleFunc("/api/clients/{clientId}", client.DeleteClientHandler).Methods("DELETE")
    r.HandleFunc("/api/clients/{clientId}/verify", client.VerifyClientHandler).Methods("POST")

	return r
}