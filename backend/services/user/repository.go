package user

import (
	"backend/database"
	"fmt"
	"log"
)

// User struct represents a user in the system (no password).
type User struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	Status    string `json:"status"`
}

// UserRepository handles DB logic for users.
type UserRepository struct{}

// NewUserRepository initializes a new UserRepository and ensures table exists.
func NewUserRepository() *UserRepository {
	repo := &UserRepository{}
	repo.InitUserTable()
	return repo
}

// InitUserTable creates the users table if not exists (no password)
func (r *UserRepository) InitUserTable() {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		first_name VARCHAR(100),
		last_name VARCHAR(100),
		email VARCHAR(100) UNIQUE NOT NULL,
		role ENUM('Admin', 'Agent') NOT NULL DEFAULT 'Agent',
		status ENUM('active', 'inactive') NOT NULL DEFAULT 'active'
	);`

	_, err := database.DB.Exec(query)
	if err != nil {
		log.Fatal("❌ Error creating users table:", err)
	}

	fmt.Println("✅ Users table ready (no password stored)")

	// 🔽 Insert seed users if not exists
	seed := `
	INSERT IGNORE INTO users (id, first_name, last_name, email, role, status)
	VALUES 
		(1, 'Admin', 'Root', 'Admin@root.com', 'Admin', 'active'),
	(2, 'Agent1', 'Agent', 'Agent1@agent.com', 'Agent', 'active'),
	(3, 'Admin2', 'Admin', 'Admin2@notRoot.com', 'Admin', 'active'),
	(4, 'Admin3', 'Admin', 'Admin3@notRoot.com', 'Admin', 'active'),
	(5, 'Admin4', 'Admin', 'Admin4@notRoot.com', 'Admin', 'active'),
	(6, 'Agent2', 'Agent', 'Agent2@agent.com', 'Agent', 'active'),
	(7, 'Agent3', 'Agent', 'Agent3@agent.com', 'Agent', 'active'),
	(8, 'Agent4', 'Agent', 'Agent4@agent.com', 'Agent', 'active');`

	_, err = database.DB.Exec(seed)
	if err != nil {
		log.Fatal("❌ Failed to seed users:", err)
	}

	fmt.Println("✅ Seed users inserted (if not already present)")
}


// CreateUser inserts metadata for a new user (already registered in Cognito).
func (r *UserRepository) CreateUser(firstName, lastName, email, role string) (User, error) {
	query := `
	INSERT INTO users (first_name, last_name, email, role, status)
	VALUES (?, ?, ?, ?, 'active')`

	result, err := database.DB.Exec(query, firstName, lastName, email, role)
	if err != nil {
		return User{}, fmt.Errorf("failed to insert user: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return User{}, fmt.Errorf("failed to retrieve user ID: %v", err)
	}

	return User{
		ID:        int(id),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Role:      role,
		Status:    "active",
	}, nil
}

// GetUserByEmail returns user by email.
func (r *UserRepository) GetUserByEmail(email string) (User, error) {
	var user User
	query := `SELECT id, first_name, last_name, email, role, status FROM users WHERE email = ?`
	err := database.DB.QueryRow(query, email).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Role, &user.Status,
	)
	if err != nil {
		return User{}, fmt.Errorf("failed to fetch user: %v", err)
	}
	return user, nil
}

// GetUserByID returns user by ID.
func (r *UserRepository) GetUserByID(userID string) (User, error) {
	var user User
	query := `SELECT id, first_name, last_name, email, role, status FROM users WHERE id = ?`
	err := database.DB.QueryRow(query, userID).Scan(
		&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Role, &user.Status,
	)
	if err != nil {
		return User{}, fmt.Errorf("failed to fetch user by ID: %v", err)
	}
	return user, nil
}

// DisableUser sets status = 'inactive'
func (r *UserRepository) DisableUser(userID string) error {
	_, err := database.DB.Exec(`UPDATE users SET status = 'inactive' WHERE id = ?`, userID)
	return err
}

// UpdateUser allows admins to modify user fields.
func (r *UserRepository) UpdateUser(userID string, user User) error {
	_, err := database.DB.Exec(`
		UPDATE users SET first_name = ?, last_name = ?, email = ?, role = ? WHERE id = ?
	`, user.FirstName, user.LastName, user.Email, user.Role, userID)
	return err
}

// SyncOrInsertUserByEmailAndRole ensures user exists, else inserts (used by JWT middleware)
func (r *UserRepository) SyncOrInsertUserByEmailAndRole(email, role string) (int, error) {
	user, err := r.GetUserByEmail(email)
	if err == nil {
		return user.ID, nil // Already exists
	}

	// Insert placeholder with minimal info
	result, err := database.DB.Exec(`
		INSERT INTO users (email, role, status)
		VALUES (?, ?, 'active')
	`, email, role)
	if err != nil {
		return 0, fmt.Errorf("failed to sync user: %v", err)
	}

	newID, _ := result.LastInsertId()
	return int(newID), nil
}
// InsertUserFromCognito inserts a user synced from Cognito (no password stored)
func (r *UserRepository) InsertUserFromCognito(email, role string) (User, error) {
	query := `
		INSERT INTO users (first_name, last_name, email, role, status)
		VALUES ('', '', ?, ?, 'active')
	`
	result, err := database.DB.Exec(query, email, role)
	if err != nil {
		return User{}, fmt.Errorf("failed to insert Cognito user: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return User{}, fmt.Errorf("failed to retrieve inserted ID: %v", err)
	}

	return User{
		ID:     int(id),
		Email:  email,
		Role:   role,
		Status: "active",
	}, nil
}
