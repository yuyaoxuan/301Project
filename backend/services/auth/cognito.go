package auth

import (
	"context"
	"fmt"
	"os"
	_ "backend/services/envloader"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/aws/credentials"
	cognito "github.com/aws/aws-sdk-go/service/cognitoidentityprovider"
	"github.com/coreos/go-oidc/v3/oidc"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"

)

var (
	ClientID      = os.Getenv("COGNITO_CLIENT_ID")
	ClientSecret  = os.Getenv("COGNITO_CLIENT_SECRET")
	CognitoDomain = os.Getenv("COGNITO_DOMAIN")
	AWSRegion      = os.Getenv("AWS_REGION")
	userPoolID    = extractUserPoolIDFromDomain(CognitoDomain)
	AuthURL       = CognitoDomain + "/oauth2/authorize"
	TokenURL      = CognitoDomain + "/oauth2/token"

	provider    *oidc.Provider
	verifier    *oidc.IDTokenVerifier
	CognitoSvc  *cognito.CognitoIdentityProvider
)
func calculateSecretHash(username, clientID, clientSecret string) string {
	mac := hmac.New(sha256.New, []byte(clientSecret))
	mac.Write([]byte(username + clientID))
	return base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func init() {

	var err error

	// Init OIDC Verifier
	provider, err = oidc.NewProvider(context.Background(), "https://cognito-idp.ap-southeast-1.amazonaws.com/"+userPoolID)
	if err != nil {
		panic("Failed to initialize OIDC provider: " + err.Error())
	}
	verifier = provider.Verifier(&oidc.Config{ClientID: ClientID})

	// Init AWS Cognito client
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ACCESS_KEY_ID"),
			os.Getenv("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	}))
	CognitoSvc = cognito.New(sess)
}


func LoginWithCognito(email, password string) (map[string]string, error) {
	secretHash := calculateSecretHash(email, ClientID, ClientSecret)

	input := &cognito.InitiateAuthInput{
		AuthFlow: aws.String("USER_PASSWORD_AUTH"),
		ClientId: aws.String(ClientID),
		AuthParameters: map[string]*string{
			"USERNAME":     aws.String(email),
			"PASSWORD":     aws.String(password),
			"SECRET_HASH":  aws.String(secretHash),
		},
	}

	result, err := CognitoSvc.InitiateAuth(input)
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"access_token":  *result.AuthenticationResult.AccessToken,
		"id_token":      *result.AuthenticationResult.IdToken,
		"refresh_token": *result.AuthenticationResult.RefreshToken,
		"expires_in":    fmt.Sprint(*result.AuthenticationResult.ExpiresIn),
	}, nil
}


// Verifier returns the OIDC verifier
func Verifier() *oidc.IDTokenVerifier {
	return verifier
}

// RegisterUserInCognito registers a user with email/password and assigns group
func RegisterUserInCognito(email, password, group string) error {
	_, err := CognitoSvc.AdminCreateUser(&cognito.AdminCreateUserInput{
		UserPoolId: aws.String(userPoolID),
		Username:   aws.String(email),
		UserAttributes: []*cognito.AttributeType{
			{
				Name:  aws.String("email"),
				Value: aws.String(email),
			},
			{
				Name:  aws.String("email_verified"),
				Value: aws.String("true"),
			},
		},
		MessageAction: aws.String("SUPPRESS"),
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	_, err = CognitoSvc.AdminSetUserPassword(&cognito.AdminSetUserPasswordInput{
		UserPoolId: aws.String(userPoolID),
		Username:   aws.String(email),
		Password:   aws.String(password),
		Permanent:  aws.Bool(true),
	})
	if err != nil {
		return fmt.Errorf("failed to set password: %v", err)
	}

	_, err = CognitoSvc.AdminAddUserToGroup(&cognito.AdminAddUserToGroupInput{
		GroupName:  aws.String(group),
		UserPoolId: aws.String(userPoolID),
		Username:   aws.String(email),
	})
	if err != nil {
		return fmt.Errorf("failed to assign group: %v", err)
	}

	return nil
}

// Hardcoded for now — improve later if needed
func extractUserPoolIDFromDomain(domain string) string {
	// Example domain: https://ap-southeast-1ztmaj2omi.auth.ap-southeast-1.amazoncognito.com
	return "ap-southeast-1_ZTmaj2omi"
}
