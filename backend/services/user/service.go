package user

// Logic File
import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
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
func (s *UserService) CreateUser(firstName, lastName, email, password, role string) (User, error) {
	// Validate inputs
	if firstName == "" || lastName == "" || email == "" || password == "" {
		return User{}, errors.New("missing required fields")
	}

	// Hash the password before storing
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, fmt.Errorf("failed to hash password: %v", err)
	}

	// Call repository function to insert user
	user, err := s.repo.CreateUser(firstName, lastName, email, string(hashedPassword), role)
	if err != nil {
		return User{}, fmt.Errorf("failed to create user: %v", err)
	}

	return user, nil
}

// // DisableUser service function
// func (s *UserService) DisableUser(targetUserID string, requesterID int, requesterRole string) error {
// 	targetUser, err := s.repo.GetUserByID(targetUserID)
// 	if err != nil {
// 		return err
// 	}

// 	// Rule 1: Prevent root admin deletion
// 	if targetUser.Role == "Admin" && targetUser.ID == 1 {
// 		return errors.New("cannot delete root admin")
// 	}

// 	// Rule 2: Only root admin can delete other admins
// 	if targetUser.Role == "Admin" && requesterID != 1 {
// 		return errors.New("only root admin can delete other admins")
// 	}

// 	return s.repo.DisableUser(targetUserID)
// }

// Update user details
func (s *UserService) UpdateUser(userID string, user User) error {
	if user.FirstName == "" || user.LastName == "" || user.Email == "" {
		return errors.New("missing required fields")
	}

	err := s.repo.UpdateUser(userID, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}
	return nil
}

// GetUserByEmail retrieves a user by email
func (s *UserService) GetUserByEmail(email string) (User, error) {
	return s.repo.GetUserByEmail(email)
}

// AuthenticateUser for Login
func (s *UserService) AuthenticateUser(email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Compare the hashed password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT Token
	token, err := GenerateJWT(user.ID, user.Role) // ✅ Use GenerateJWT from jwt_utils.go
	if err != nil {
		return "", err
	}

	return token, nil
}

// ResetPassword allows users to reset their password
func (s *UserService) ResetPassword(email, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(email, string(hashedPassword))
}
