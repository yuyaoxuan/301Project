package user

import (
	"errors"
	"fmt"
)

// UserService struct to interact with the repository layer
type UserService struct {
	repo *UserRepository
}

// NewUserService initializes the user service
func NewUserService(repo *UserRepository) *UserService {
	return &UserService{repo: repo}
}

// CreateUser processes user creation request
func (s *UserService) CreateUser(firstName, lastName, email, role string) (User, error) {
	// Validate inputs
	if firstName == "" || lastName == "" || email == "" {
		return User{}, errors.New("missing required user fields")
	}

	// Ensure role is valid
	if role != "Admin" && role != "Agent" {
		return User{}, errors.New("invalid role: must be 'Admin' or 'Agent'")
	}

	// Call repository function to insert user
	user, err := s.repo.CreateUser(firstName, lastName, email, role)
	if err != nil {
		return User{}, fmt.Errorf("failed to create user: %v", err)
	}

	return user, nil
}
