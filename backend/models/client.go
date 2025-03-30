package models

// Client struct represents a client in the system
type Client struct {
	ClientID           string `json:"client_id"`
	FirstName          string `json:"first_name"`
	LastName           string `json:"last_name"`
	DOB                string `json:"dob"`
	Gender             string `json:"gender"`
	Email              string `json:"email"`
	Phone              string `json:"phone"`
	Address            string `json:"address"`
	City               string `json:"city"`
	State              string `json:"state"`
	Country            string `json:"country"`
	PostalCode         string `json:"postal_code"`
	VerificationStatus string `json:"verification_status"`
}