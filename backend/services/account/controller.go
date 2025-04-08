package account

import (
	"backend/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// CreateAccountHandler restricts to Admin or Agent
func CreateAccountHandler(service *AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userCtx := r.Context().Value("user").(map[string]interface{})
		role := userCtx["role"].(string)
		if role != "Admin" && role != "Agent" {
			http.Error(w, "Unauthorized: only Admin or Agent can create accounts", http.StatusForbidden)
			return
		}

		var account models.Account
		if err := json.NewDecoder(r.Body).Decode(&account); err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		createdAccount, err := service.CreateAccount(account)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdAccount)
	}
}

// DeleteAccountHandler restricts to Admin or Agent
func DeleteAccountHandler(service *AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userCtx := r.Context().Value("user").(map[string]interface{})
		role := userCtx["role"].(string)
		if role != "Admin" && role != "Agent" {
			http.Error(w, "Unauthorized: only Admin or Agent can delete accounts", http.StatusForbidden)
			return
		}

		vars := mux.Vars(r)
		accountID, err := strconv.Atoi(vars["account_id"])
		if err != nil {
			http.Error(w, "Invalid account ID", http.StatusBadRequest)
			return
		}

		if err := service.DeleteAccount(accountID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"message": "Account deleted successfully"})
	}
}

// GetAllAccountsHandler restricts to Admin or Agent
func GetAllAccountsHandler(service *AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userCtx := r.Context().Value("user").(map[string]interface{})
		role := userCtx["role"].(string)
		if role != "Admin" && role != "Agent" {
			http.Error(w, "Unauthorized: only Admin or Agent can view accounts", http.StatusForbidden)
			return
		}

		accounts, err := service.GetAllAccounts()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(accounts)
	}
}

// GetAccountsByClientHandler restricts to Admin or Agent
func GetAccountsByClientHandler(service *AccountService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userCtx := r.Context().Value("user").(map[string]interface{})
		role := userCtx["role"].(string)
		if role != "Admin" && role != "Agent" {
			http.Error(w, "Unauthorized: only Admin or Agent can view client accounts", http.StatusForbidden)
			return
		}

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
