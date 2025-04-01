package communicationlogs

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	// "github.com/aws/aws-sdk-go-v2/aws"
	// "github.com/aws/aws-sdk-go-v2/config"
	// "github.com/aws/aws-sdk-go-v2/service/ses"
	// "github.com/aws/aws-sdk-go-v2/service/ses/types"
)

type EmailSender struct {
	client *ses.Client
	from   string
}

func NewEmailSender() (*EmailSender, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS config: %w", err)
	}

	fromEmail := os.Getenv("EMAIL_SENDER") // Make sure to set this in your environment
	if fromEmail == "" {
		return nil, fmt.Errorf("EMAIL_SENDER environment variable not set")
	}

	client := ses.NewFromConfig(cfg)
	return &EmailSender{
		client: client,
		from:   fromEmail,
	}, nil
}

func (s *EmailSender) SendEmail(to, subject, body string) error {
	// input := &ses.SendEmailInput{
	// 	Source: aws.String(s.from),
	// 	Destination: &types.Destination{
	// 		ToAddresses: []string{to},
	// 	},
	// 	Message: &types.Message{
	// 		Subject: &types.Content{
	// 			Data: aws.String(subject),
	// 		},
	// 		Body: &types.Body{
	// 			Text: &types.Content{
	// 				Data: aws.String(body),
	// 			},
	// 		},
	// 	},
	// }

	// _, err := s.client.SendEmail(context.TODO(), input)
	// if err != nil {
	// 	return fmt.Errorf("failed to send email: %w", err)
	// }
	// return nil
	fmt.Println("🚀 [Mock Email Sent]")
	fmt.Println("To:", to)
	fmt.Println("Subject:", subject)
	fmt.Println("Body:", body)
	return nil // Simulate success
}
