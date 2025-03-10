package main

import (
	"fmt"
	"log"
	"net/http"

	"backend/database"
	"backend/routes"          // Import routes from the routes package
	"backend/services/client" // Import client service to initialize table
	"backend/services/user"   // Import user service to initialize table
)

func main() {

	database.ConnectDB()

	// Ensure user table is created
	userRepo := user.NewUserRepository() // Initializes user repo (which ensures table exists)
	_ = userRepo // Avoid unused variable warning

	// Ensure user table is created
	clientRepo := client.NewClientRepository() // Initializes user repo (which ensures table exists)
	_ = clientRepo // Avoid unused variable warning
	
	router := routes.SetupRoutes()
	fmt.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}