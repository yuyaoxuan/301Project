package routes

import (
	"net/http"

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
	r.HandleFunc("/api/clients", client.CreateClientHandler).Methods("POST")
    r.HandleFunc("/api/clients/{clientId}", client.GetClientHandler).Methods("GET")
    r.HandleFunc("/api/clients/{clientId}", client.UpdateClientHandler).Methods("PUT")
    r.HandleFunc("/api/clients/{clientId}", client.DeleteClientHandler).Methods("DELETE")
    r.HandleFunc("/api/clients/{clientId}/verify", client.VerifyClientHandler).Methods("POST")

	// ACCOUNT ROUTES
	r.HandleFunc("/api/clients/{clientId}/accounts", client.CreateAccountHandler).Methods("POST")

	return r
}
