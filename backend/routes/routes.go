package routes

import (
	"net/http"

	"backend/services/account" // import account service
	"backend/services/agentClient"
	"backend/services/agentclient_logs" // import transaction logs service
	"backend/services/client"
	"backend/services/user" // import user service

	"github.com/gorilla/mux"
)

// SetupRoutes initializes the router and returns it
func SetupRoutes(clientService *client.ClientService, accountService *account.AccountService) *mux.Router {
	r := mux.NewRouter()

	// Health Check Route
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is running!"))
	}).Methods("GET")

	// Public Routes (No Authentication Needed)
	r.HandleFunc("/api/users/authenticate", user.AuthenticateUserHandler).Methods("POST") // Login
	r.HandleFunc("/api/users", user.CreateUserHandler).Methods("POST")    
	r.HandleFunc("/api/users/reset-password", user.ResetPasswordHandler).Methods("POST")// Register User

	// CLIENT ROUTES
	r.HandleFunc("/api/clients/{agent_id}", client.CreateClientHandler(clientService)).Methods("POST")
	r.HandleFunc("/api/clients/{clientId}", client.GetClientHandler(clientService)).Methods("GET")
	r.HandleFunc("/api/clients/{agent_id}/{clientId}", client.UpdateClientHandler(clientService)).Methods("PUT")
	r.HandleFunc("/api/clients/{clientId}", client.DeleteClientHandler(clientService)).Methods("DELETE")
	r.HandleFunc("/api/clients/{clientId}/verify", client.VerifyClientHandler(clientService)).Methods("POST")

	// ACCOUNT Routes
	r.HandleFunc("/api/accounts", account.CreateAccountHandler(accountService)).Methods("POST")
	r.HandleFunc("/api/accounts/{account_id}", account.DeleteAccountHandler(accountService)).Methods("DELETE")

	// agentClient routes
	r.HandleFunc("/agentClient", agentClient.AssignAgentsToUnassignedClientsHandler).Methods("PUT")

	// Protected Routes (Require JWT)
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(user.JWTAuthMiddleware) // Apply JWT Middleware

	// protected.HandleFunc("/users/{userId}", user.DisableUserHandler).Methods("DELETE")       // Disable User
	protected.HandleFunc("/users/{userId}", user.UpdateUserHandler).Methods("PUT")           // Update User
	protected.HandleFunc("/users/reset-password", user.ResetPasswordHandler).Methods("POST") // Reset Password

	// // AgentClient logs routes
	r.HandleFunc("/agentclient_logs/client/{clientID}", agentclient_logs.GetAgentClientLogsByClientHandler).Methods("GET")
	r.HandleFunc("/agentclient_logs/agent/{agentID}", agentclient_logs.GetAgentClientLogsByAgentHandler).Methods("GET")
	r.HandleFunc("/agentclient_logs", agentclient_logs.GetAllAgentClientLogsHandler).Methods("GET") // Get all agent-client logs

	// // Bank account logs routes
	r.HandleFunc("/agentclient_logs/account/client/{clientID}", agentclient_logs.GetAccountLogsByClientHandler).Methods("GET")
	r.HandleFunc("/agentclient_logs/account/agent/{agentID}", agentclient_logs.GetAccountLogsByAgentHandler).Methods("GET")
	r.HandleFunc("/agentclient_logs/account", agentclient_logs.GetAllAccountLogsHandler).Methods("GET") // Get all bank account logs

	// // Combined logs route (both client and bank account logs)
	r.HandleFunc("/agentclient_logs/all", agentclient_logs.GetAllLogsHandler).Methods("GET") // Get all logs (client + bank account)

	// // Delete log (generalized for all types)
	r.HandleFunc("/agentclient_logs/{logID}", agentclient_logs.DeleteLogHandler).Methods("DELETE")
	

	return r
}
