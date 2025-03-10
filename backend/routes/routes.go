package routes

import (
	"net/http"

	"backend/services/client" // import user service
	"backend/services/user"   // import user service

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
	r.HandleFunc("/clients", client.CreateClientHandler).Methods("POST")
	r.HandleFunc("/accounts", client.CreateAccountHandler).Methods("POST")
	r.HandleFunc("/accounts/{client_id}", client.DeleteAccountHandler).Methods("DELETE")
	
	return r
}