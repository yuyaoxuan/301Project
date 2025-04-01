package communicationlogs

import (
	"backend/models"
	"fmt"
)

// CommunicationLogService handles log operations
type CommunicationLogService struct {
	repo *CommunicationLogRepository
}

// NewCommunicationLogService initializes the service
func NewCommunicationLogService(repo *CommunicationLogRepository) *CommunicationLogService {
	return &CommunicationLogService{repo: repo}
}

func (s *CommunicationLogService) LogCommunication(agentClientLog models.AgentClientLog) error {
	// Extract values from the AgentClientLog
	clientID := agentClientLog.ClientID
	agentID := agentClientLog.AgentID
	action := agentClientLog.Action
	modifiedFields := agentClientLog.ModifiedFields
	logType := modifiedFields["log_type"].(string)

	// Generate the email subject based on the action
	var emailSubject string
	if action == "Create" {
		emailSubject = fmt.Sprintf("%s created!", logType)
	} else if action == "Update" {
		emailSubject = fmt.Sprintf("%s updated!", logType)
	} else if action == "Delete" {
		emailSubject = fmt.Sprintf("%s deleted!", logType)
	} else {
		emailSubject = fmt.Sprintf("%s action performed!", action)
	}

	// Commented out the body generation for local testing
	/*
		body := fmt.Sprintf("Action: %s\nClient ID: %s\n", action, clientID)
		if details, ok := modifiedFields["details"].(*models.Client); ok {
			// Directly access fields of the *models.Client struct
			body += fmt.Sprintf("Name: %s\nEmail: %s\nPhone: %s\n", details.Name, details.Email, details.Phone)
		} else {
			body += fmt.Sprintf("Details format is not recognized.\n")
		}
	*/

	// For testing, just log the details
	fmt.Printf("Test log: %s - %s\n", emailSubject, clientID)

	// Create CommunicationLog object with logID
	logID := agentClientLog.ID
	emailStatus := "Pending" // Default status (you can change based on email sending result)

	// Skip actual email sending for local testing
	// emailSender, err := NewEmailSender()
	// if err != nil {
	// 	return err
	// }

	// Replace this with the actual recipient's email (from DB or payload)
	// to := "recipient@example.com"
	// if err := emailSender.SendEmail(to, emailSubject, body); err != nil {
	// 	emailStatus = "Failed"
	// } else {
	// 	emailStatus = "Sent"
	// }

	// Just simulate success here for testing
	emailStatus = "Sent"

	// Insert into CommunicationLog repository
	return s.repo.InsertCommunicationLog(logID, clientID, agentID, emailSubject, emailStatus)
}

// GetClientCommunicationLogsByLogID to get communication by logID
func (s *CommunicationLogService) GetCommunicationLogByLogID(logID int) (models.CommunicationLog, error) {
	return s.repo.GetCommunicationLogByLogID(logID)
}
