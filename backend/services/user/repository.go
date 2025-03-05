package user

import (
	"backend/database"
	"fmt"
	"log"
)

// User struct represents a user in the system
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
}

// UserRepository struct for interacting with database
type UserRepository struct{}

// NewUserRepository initializes a new UserRepository
func NewUserRepository() *UserRepository {
	repo := &UserRepository{}
	repo.InitUserTable() // ✅ Ensure table exists when the repository is created
	return repo
}

// InitUserTable ensures the users table exists when the service starts
func (r *UserRepository) InitUserTable() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		first_name VARCHAR(100) NOT NULL,
		last_name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		role ENUM('Admin', 'Agent') NOT NULL DEFAULT 'Agent'
	);`

	_, err := database.DB.Exec(query)
	if err != nil {
		log.Fatal("❌ Error creating users table:", err)
	}

	fmt.Println("✅ Users table checked/created!")
}

// CreateUser inserts a new user into the database
func (r *UserRepository) CreateUser(firstName, lastName, email, role string) (User, error) {
	query := "INSERT INTO users (first_name, last_name, email, role) VALUES (?, ?, ?, ?)"
	result, err := database.DB.Exec(query, firstName, lastName, email, role)
	if err != nil {
		return User{}, fmt.Errorf("failed to insert user: %v", err)
	}

	// Get the last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return User{}, fmt.Errorf("failed to retrieve inserted user ID: %v", err)
	}

	// Return the created user
	return User{
		ID:        int(id),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Role:      role,
	}, nil
}
