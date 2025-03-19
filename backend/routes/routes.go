package routes

import (
	"net/http"

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

	// TRANSACTION LOGS ROUTES
	r.HandleFunc("/agentclient_logs", agentclientlogs.CreateAgentClientLogHandler).Methods("POST")
	r.HandleFunc("/agentclient_logs/{logID}", agentclientlogs.DeleteAgentClientLogHandler).Methods("DELETE")
	r.HandleFunc("/agentclient_logs/client/{clientID}", agentclientlogs.GetAgentClientLogsByClientHandler).Methods("GET")
	r.HandleFunc("/agentclient_logs/agent/{agentID}", agentclientlogs.GetAgentClientLogsByAgentHandler).Methods("GET")
	r.HandleFunc("/agentclient_logs", agentclientlogs.GetAllAgentClientLogsHandler).Methods("GET")

	return r
}
