package routes

import (
	"net/http"

	"backend/services/account"
	"backend/services/agentClient"
	"backend/services/agentclient_logs"
	"backend/services/client"
	communicationlogs "backend/services/communication_logs"
	"backend/services/user"

	"github.com/gorilla/mux"
)

func SetupRoutes(
	clientService *client.ClientService,
	accountService *account.AccountService,
	agentClientLogService *agentclient_logs.AgentClientLogService,
) *mux.Router {
	r := mux.NewRouter()

	// Health Check Route
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("API is running!"))
	}).Methods("GET")

	// Public User Routes
	r.HandleFunc("/api/users/authenticate", user.AuthenticateUserHandler).Methods("POST")
	r.HandleFunc("/api/users", user.CreateUserHandler).Methods("POST")
	r.HandleFunc("/api/users/reset-password", user.ResetPasswordHandler).Methods("POST")

	// Protected Routes (Require JWT)
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(user.JWTAuthMiddleware)

	// User Routes (protected)
	protected.HandleFunc("/users/{userId}", user.DisableUserHandler).Methods("DELETE")
	protected.HandleFunc("/users/{userId}", user.UpdateUserHandler).Methods("PUT")
	protected.HandleFunc("/users/reset-password", user.ResetPasswordHandler).Methods("POST")

	// Client Routes
	r.HandleFunc("/api/clients/{agent_id}", client.CreateClientHandler(clientService)).Methods("POST")
	r.HandleFunc("/api/clients/{clientId}", client.GetClientHandler(clientService)).Methods("GET")
	r.HandleFunc("/api/clients/{agent_id}/{clientId}", client.UpdateClientHandler(clientService)).Methods("PUT")
	r.HandleFunc("/api/clients/{clientId}", client.DeleteClientHandler(clientService)).Methods("DELETE")
	r.HandleFunc("/api/clients/{clientId}/verify", client.VerifyClientHandler(clientService)).Methods("POST")

	// Account Routes
	r.HandleFunc("/api/accounts", account.CreateAccountHandler(accountService)).Methods("POST")
	r.HandleFunc("/api/accounts/{account_id}", account.DeleteAccountHandler(accountService)).Methods("DELETE")

	// 
	r.HandleFunc("/api/agentclient", agentClient.AssignAgentsToUnassignedClientsHandler).Methods("PUT")
	
	// Agent Client Log Read Routes
	r.HandleFunc("/agentclient_logs/client/{clientID}", agentclient_logs.GetAgentClientLogsByClientHandler(agentClientLogService)).Methods("GET")
	r.HandleFunc("/agentclient_logs/agent/{agentID}", agentclient_logs.GetAgentClientLogsByAgentHandler(agentClientLogService)).Methods("GET")
	r.HandleFunc("/agentclient_logs", agentclient_logs.GetAllAgentClientLogsHandler(agentClientLogService)).Methods("GET")
	r.HandleFunc("/agentclient_logs/account/client/{clientID}", agentclient_logs.GetAccountLogsByClientHandler(agentClientLogService)).Methods("GET")
	r.HandleFunc("/agentclient_logs/account/agent/{agentID}", agentclient_logs.GetAccountLogsByAgentHandler(agentClientLogService)).Methods("GET")
	r.HandleFunc("/agentclient_logs/account", agentclient_logs.GetAllAccountLogsHandler(agentClientLogService)).Methods("GET")
	r.HandleFunc("/agentclient_logs/all/client/{clientID}", agentclient_logs.GetClientAndAccountLogsByClientIDHandler(agentClientLogService)).Methods("GET")
	r.HandleFunc("/agentclient_logs/all/agent/{agentID}", agentclient_logs.GetClientAndAccountLogsByAgentIDHandler(agentClientLogService)).Methods("GET")
	r.HandleFunc("/agentclient_logs/all", agentclient_logs.GetAllLogsHandler(agentClientLogService)).Methods("GET")
	r.HandleFunc("/agentclient_logs/{logID}", agentclient_logs.DeleteLogHandler(agentClientLogService)).Methods("DELETE")

	// Communication Log Read Routes
	r.HandleFunc("/communication_logs/{logID}", communicationlogs.GetCommunicationLogByLogIDHandler).Methods("GET")

	return r
}
