package communication_observer

import (
	"backend/models"
	communicationlogs "backend/services/communication_logs"
	"fmt"
)

type CommunicationObserver struct {
	LogService *communicationlogs.CommunicationLogService
}

func (co *CommunicationObserver) NotifyCreate(agentID int, clientID string, object interface{}) {
	log, ok := object.(models.AgentClientLog) // ✅ safe type assertion
	if !ok {
		fmt.Println("❌ CommunicationObserver: expected AgentClientLog, got something else")
		return
	}

	err := co.LogService.LogCommunication(log)
	if err != nil {
		fmt.Println("❌ Failed to log communication:", err)
	}
}

func (co *CommunicationObserver) NotifyUpdate(agentID int, clientID string, before, after interface{}) {
	// Leave blank or implement if needed
}

func (co *CommunicationObserver) NotifyDelete(agentID int, clientID string, object interface{}) {
	// Leave blank or implement if needed
}
