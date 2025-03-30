package main

import (
	"fmt"
	"log"
	"net/http"

	"backend/database"
	"backend/routes"           // Import routes from the routes package
	"backend/services/account" // Import acciunt service to initialize table

	"backend/services/agentClient"
	"backend/services/client" // Import client service to initialize table
	"backend/services/user"   // Import user service to initialize table

	"backend/services/agentclient_logs" // Import agent-client logs to initialize its table
	"backend/services/observer"         // import observer
)

func main() {

	database.ConnectDB()

	// Ensure user table is created
	userRepo := user.NewUserRepository() // Initializes user repo (which ensures table exists)
	_ = userRepo                         // Avoid unused variable warning

	// Initialize the ObserverManager
	observerManager := &observer.ObserverManager{}

	// Initialize repositories
	agentClientLogRepo := agentclient_logs.NewAgentClientLogRepository() // Agent-client logs repo
	agentClientRepo := agentClient.NewAgentClientRepository()
	agentClientService := agentClient.NewAgentClientService(agentClientRepo)
	_ = agentClientRepo

	// Create the LogService which will use the repository to log actions
	logService := agentclient_logs.NewAgentClientLogService(agentClientLogRepo)

	// Create client and account observers, passing the LogService to them
	clientObserver := &observer.ClientObserver{LogService: logService}
	accountObserver := &observer.AccountObserver{LogService: logService}

	// Register observers
	observerManager.AddClientObserver(clientObserver)
	observerManager.AddAccountObserver(accountObserver)

	// Initialize repo and services
	clientRepo := client.NewClientRepository(observerManager)            
	clientService := client.NewClientService(clientRepo, observerManager)
	
	accountRepo := account.NewAccountRepository(observerManager, clientRepo, agentClientRepo) 
	accountService := account.NewAccountService(observerManager, accountRepo, agentClientService, clientService)
	
	// Set up routes
	router := routes.SetupRoutes(clientService, accountService)
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
