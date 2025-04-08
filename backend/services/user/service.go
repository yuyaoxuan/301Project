package user

import (
	"errors"
	"fmt"
)

// UserService handles business logic for users.
type UserService struct {
	repo *UserRepository
}

// NewUserService initializes a new UserService.
func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser stores user metadata (after Cognito registration)
func (s *UserService) CreateUser(firstName, lastName, email, role string) (User, error) {
	if firstName == "" || lastName == "" || email == "" || role == "" {
		return User{}, errors.New("missing required fields")
	}

	user, err := s.repo.CreateUser(firstName, lastName, email, role)
	if err != nil {
		return User{}, fmt.Errorf("failed to create user: %v", err)
	}

	return user, nil
}

// DisableUser disables a user, checking business rules
func (s *UserService) DisableUser(targetUserID string, requesterID int, requesterRole string) error {
	targetUser, err := s.repo.GetUserByID(targetUserID)
	if err != nil {
		return fmt.Errorf("target user not found: %v", err)
	}

	// Prevent root admin deletion
	if targetUser.Role == "Admin" && targetUser.ID == 1 {
		return errors.New("cannot disable root admin")
	}

	// Only root admin (ID 1) can disable other admins
	if targetUser.Role == "Admin" && requesterID != 1 {
		return errors.New("only root admin can disable other admins")
	}

	return s.repo.DisableUser(targetUserID)
}

// UpdateUser updates an existing user's details.
func (s *UserService) UpdateUser(userID string, user User) error {
	err := s.repo.UpdateUser(userID, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}

// GetUserByEmail retrieves a user's details by their email.
func (s *UserService) GetUserByEmail(email string) (User, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return User{}, fmt.Errorf("failed to fetch user by email: %v", err)
	}
	return user, nil
}

// SyncOrInsertUser checks if the user exists by email; if not, inserts it with given role.
func SyncOrInsertUser(email, role string) (int, error) {
	repo := NewUserRepository()
	existingUser, err := repo.GetUserByEmail(email)
	if err == nil {
		// User exists
		return existingUser.ID, nil
	}

	// Insert new user with empty name fields
	newUser, err := repo.InsertUserFromCognito(email, role)
	if err != nil {
		return 0, fmt.Errorf("failed to insert user: %v", err)
	}

	return newUser.ID, nil
}
