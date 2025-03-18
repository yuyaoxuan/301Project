package user
//Logic File
import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var jwtSecret = []byte("your-secret-key")


func generateJWT(userID int, role string) (string, error) {
	claims := jwt.MapClaims{
		"userId": userID,
		"role":   role,
		"exp":    time.Now().Add(time.Hour * 2).Unix(), // Token expires in 2 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

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
