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
	Password  string `json:"password"`
	Role      string `json:"role"`
	Status    string `json:"status"`
}

type UserRepository struct{}

// UserRepository struct for interacting with database
// GetUserByID fetches a user from the database using their ID
func (r *UserRepository) GetUserByID(userID string) (User, error) {
	var user User
	query := "SELECT id, first_name, last_name, email, role, password, status FROM users WHERE id = ?"
	err := database.DB.QueryRow(query, userID).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Role, &user.Password, &user.Status)
	if err != nil {
		return User{}, fmt.Errorf("failed to fetch user by ID: %v", err)
	}
	return user, nil
}


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
		password VARCHAR(255) NOT NULL,
		role ENUM('Admin', 'Agent') NOT NULL DEFAULT 'Agent',
		status ENUM('active', 'inactive') NOT NULL DEFAULT 'active'
	);`

	_, err := database.DB.Exec(query)
	if err != nil {
		log.Fatal("❌ Error creating users table:", err)
	}

	fmt.Println("✅ Users table checked/created!")
}

// CreateUser inserts a new user into the database
func (r *UserRepository) CreateUser(firstName, lastName, email, password, role string) (User, error) {
	query := "INSERT INTO users (first_name, last_name, email, password, role, status) VALUES (?, ?, ?, ?, ?, 'active')"
	result, err := database.DB.Exec(query, firstName, lastName, email, password, role)
	if err != nil {
		return User{}, fmt.Errorf("failed to insert user: %v", err)
	}

	// Get last inserted ID
	id, err := result.LastInsertId()
	if err != nil {
		return User{}, fmt.Errorf("failed to retrieve inserted user ID: %v", err)
	}

	// Return created user
	return User{
		ID:        int(id),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password, 
		Role:      role,
		Status:    "active",
	}, nil
}

// DisableUser and update user status
func (r *UserRepository) DisableUser(userID string) error {
	query := "UPDATE users SET status = 'inactive' WHERE id = ?"
	_, err := database.DB.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to disable user: %v", err)
	}
	return nil
}

// UpdateUser 
func (r *UserRepository) UpdateUser(userID string, user User) error {
	query := "UPDATE users SET first_name = ?, last_name = ?, email = ?, role = ? WHERE id = ?"
	_, err := database.DB.Exec(query, user.FirstName, user.LastName, user.Email, user.Role, userID)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}

//Validate user credentials and get User
func (r *UserRepository) GetUserByEmail(email string) (User, error) {
	var user User
	query := "SELECT id, first_name, last_name, email, role, password FROM users WHERE email = ?"
	err := database.DB.QueryRow(query, email).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Role, &user.Password)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

//Update New Password
func (r *UserRepository) UpdatePassword(email, hashedPassword string) error {
	query := "UPDATE users SET password = ? WHERE email = ?"
	_, err := database.DB.Exec(query, hashedPassword, email)
	if err != nil {
		return fmt.Errorf("failed to reset password: %v", err)
	}
	return nil
}
