package user

import (
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// UserService handles business logic for users.
type UserService struct {
	repo *UserRepository
}

// NewUserService initializes a new UserService.
func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser creates a new user.
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

// ResetPassword updates a user's password
func (s *UserService) ResetPassword(email, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %v", err)
	}

	return s.repo.UpdatePassword(email, string(hashedPassword))
}

// DisableUser disables a user by setting their status to 'inactive'.
func (s *UserService) DisableUser(userID string) error {
	err := s.repo.DisableUser(userID)
	if err != nil {
		return fmt.Errorf("failed to disable user: %v", err)
	}
	return nil
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
