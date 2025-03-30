package models

// Account struct
type Account struct {
	AccountID     int     `json:"account_id"`
	ClientID      string  `json:"client_id"`
	AccountType   string  `json:"account_type"`
	AccountStatus string  `json:"account_status"`
	OpeningDate   string    `json:"opening_date"`
	InitialDeposit float64 `json:"initial_deposit"`
	Currency      string  `json:"currency"`
	BranchID      string  `json:"branch_id"`
}
