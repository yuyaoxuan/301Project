package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	_ "backend/services/envloader"
	"github.com/joho/godotenv" // ✅ Load .env

	"backend/database"
	"backend/routes"           // Import routes from the routes package
	"backend/services/account" // Import acciunt service to initialize table

	"backend/services/agentClient"
	"backend/services/agentclient_logs"                     // Import agent-client logs to initialize its table
	"backend/services/client"                               // Import client service to initialize table
	communicationlogs "backend/services/communication_logs" // Import communication service to initialize table
	commobserver "backend/services/communication_observer"  // Import communication observer
	"backend/services/observer"                             // import observer
	"backend/services/user"                                 // Import user service to initialize table
)

func main() {
	// ✅ Load environment variables from .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Optional: Print confirmation that AWS creds are present
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		log.Fatal("Missing AWS credentials in .env")
	}

	// Initialize the database connection
	database.ConnectDB()

	// Initialize repositories
	agentClientLogRepo := agentclient_logs.NewAgentClientLogRepository() // Agent-client logs repo
	communicationRepo := communicationlogs.NewCommunicationLogRepository()

	// Ensure that necessary database tables are created
	userRepo := user.NewUserRepository() // Initializes user repo (which ensures table exists)
	_ = userRepo                         // Avoid unused variable warning

	// Initialize the ObserverManager and register observers
	observerManager := &observer.ObserverManager{}

	// Initialize repo and services
	clientRepo := client.NewClientRepository(observerManager)
	agentClientRepo := agentClient.NewAgentClientRepository()
	accountRepo := account.NewAccountRepository(observerManager, agentClientRepo)

	clientService := client.NewClientService(clientRepo, observerManager)
	agentClientService := agentClient.NewAgentClientService(agentClientRepo)
	accountService := account.NewAccountService(observerManager, accountRepo, agentClientService)

	clientService.SetAgentClientService(agentClientService)
	clientService.SetAccountService(accountService)
	accountService.SetClientService(clientService)

	// Create the LogService which will use the repository to log actions
	logService := agentclient_logs.NewAgentClientLogService(agentClientLogRepo, observerManager)
	communicationService := communicationlogs.NewCommunicationLogService(communicationRepo)

	clientObserver := &observer.ClientObserver{LogService: logService}
	accountObserver := &observer.AccountObserver{LogService: logService}
	communicationObserver := &commobserver.CommunicationObserver{LogService: communicationService}

	// Register observers
	observerManager.AddClientObserver(clientObserver)
	observerManager.AddAccountObserver(accountObserver)
	observerManager.AddCommunicationObserver(communicationObserver)

	// Set up routes
	router := routes.SetupRoutes(clientService, accountService, logService)

	// Start the server
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
