package main

import (
	"fmt"
	"log"
	"net/http"

	"backend/database"
	"backend/routes"                    // Import routes from the routes package
	"backend/services/agentclient_logs" // Import agent-client logs to initialize its table
	"backend/services/observer"         // import observer
	"backend/services/user"             // Import user service to initialize table
)

func main() {

	database.ConnectDB()

	// Ensure user table is created
	userRepo := user.NewUserRepository() // Initializes user repo (which ensures table exists)
	_ = userRepo                         // Avoid unused variable warning

	// For the agent-client logs table (if you have other services like this)
	agentClientLogRepo := agentclient_logs.NewAgentClientLogRepository() // Initializes agent-client logs repo (which ensures table exists)
	_ = agentClientLogRepo

	// Create the LogService which will use the repository to log actions
	logService := agentclient_logs.NewAgentClientLogService(agentClientLogRepo)

	// Create the observer manager
	observerManager := &observer.ObserverManager{}

	// Create client and account observers, passing the LogService to them
	clientObserver := &observer.ClientObserver{LogService: logService}
	accountObserver := &observer.AccountObserver{LogService: logService}

	// Register observers
	observerManager.AddClientObserver(clientObserver)
	observerManager.AddAccountObserver(accountObserver)

	router := routes.SetupRoutes()
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
