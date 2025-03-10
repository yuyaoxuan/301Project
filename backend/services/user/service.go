package user
//Logic File
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

// DisableUser service function
func (s *UserService) DisableUser(userID string) error {
	return s.repo.DisableUser(userID)
}

// Update user Details
func (s *UserService) UpdateUser(userID string, user User) error {
	if user.FirstName == "" || user.LastName == "" || user.Email == "" {
		return errors.New("missing required fields")
	}

	return s.repo.UpdateUser(userID, user)
}

//Authenticate User for Login
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

	// Generate JWT Token (OAuth2 Implementation)
	token, err := generateJWT(user.ID, user.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

//Reset Password
func (s *UserService) ResetPassword(email, newPassword string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.repo.UpdatePassword(email, string(hashedPassword))
}
