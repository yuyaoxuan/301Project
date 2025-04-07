package account

import (
	"backend/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CreateUserHandler handles the HTTP request to create a user
func CreateAccountHandler(service *AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		var account models.Account

		// Decode JSON request body into account struct
		err := json.NewDecoder(r.Body).Decode(&account)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		// Call service layer to create account
		createdAccount, err := service.CreateAccount(account)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")  // Set the response header to indicate it's JSON
		w.WriteHeader(http.StatusCreated)                   // Set the HTTP status code to 201 (Created)
		json.NewEncoder(w).Encode(createdAccount)            // Write the actual data (new client) as a response body

	}
}

// CreateUserHandler handles the HTTP request to create a user
func DeleteAccountHandler(service *AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		AccountID, stringToInt_err := strconv.Atoi(vars["account_id"])

		if stringToInt_err != nil {
			http.Error(w, "Unable to convert string to int", http.StatusInternalServerError)
			return
		 }

		// Call service layer to create account
		err := service.DeleteAccount(AccountID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Client deleted successfully"})
	
	}
}

func GetAllAccountsHandler(service *AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		accounts, err := service.GetAllAccounts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(accounts)
	}
}

func GetAccountsByClientHandler(service *AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		clientID := vars["clientId"]

		accounts, err := service.GetAccountsByClientID(clientID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(accounts)
	}
}