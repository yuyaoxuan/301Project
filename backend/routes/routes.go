
package routes

import (
	"net/http"
	"backend/services/account"
	"backend/services/agentclient_logs"
	"backend/services/client"
	"backend/services/communication_logs"
	"backend/services/user"
	"github.com/gorilla/mux"
	"backend/services/middleware"
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
r.HandleFunc("/api/users/authenticate", user.AuthenticateUserHandler).Methods("GET") // OAuth login
r.HandleFunc("/api/auth/callback", user.AuthCallbackHandler).Methods("GET")
r.HandleFunc("/api/users", user.CreateUserHandler).Methods("POST") // Registers a user into Cognito and inserts into DB
// Removed password reset endpoint here (now protected only)

// Protected Routes (Require JWT)
protected := r.PathPrefix("/api").Subrouter()
protected.Use(middleware.JWTAuthMiddleware) // Use middleware from correct package

// User Routes (protected)
protected.HandleFunc("/users/{userId}", user.DisableUserHandler).Methods("DELETE")
protected.HandleFunc("/users/{userId}", user.UpdateUserHandler).Methods("PUT")
//protected.HandleFunc("/users/reset-password", user.ResetPasswordHandler).Methods("POST") // Reset password (if supported)

	// Client Routes
	r.HandleFunc("/api/clients/{agent_id}", client.CreateClientHandler(clientService)).Methods("POST")
	r.HandleFunc("/api/clients/{clientId}", client.GetClientHandler(clientService)).Methods("GET")
	r.HandleFunc("/api/clients/{agent_id}/{clientId}", client.UpdateClientHandler(clientService)).Methods("PUT")
	r.HandleFunc("/api/clients/{clientId}", client.DeleteClientHandler(clientService)).Methods("DELETE")
	r.HandleFunc("/api/clients/{clientId}/verify", client.VerifyClientHandler(clientService)).Methods("POST")

	// Account Routes
	r.HandleFunc("/api/accounts", account.CreateAccountHandler(accountService)).Methods("POST")
	r.HandleFunc("/api/accounts/{account_id}", account.DeleteAccountHandler(accountService)).Methods("DELETE")

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
