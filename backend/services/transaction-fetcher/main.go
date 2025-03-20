package main

import (
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/joho/godotenv"
	"github.com/cs301-itsa/project-2024-25t2-g1-t5/backend/services/transaction-fetcher" // Import models
)

var (
	status          ServiceStatus
	statusMutex     sync.Mutex
	processTrigger  = make(chan bool, 1)
	shutdownTrigger = make(chan bool, 1)
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Initialize the database
	db, err := initializeDatabase()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// SFTP Config
	sftpServer := os.Getenv("SFTP_SERVER")
	username := os.Getenv("SFTP_USERNAME")
	privateKeyPath := os.Getenv("SFTP_PRIVATE_KEY")
	remotePath := os.Getenv("SFTP_REMOTE_LOG_PATH")
	processedPath := os.Getenv("SFTP_PROCESSED_PATH")

	// Debug logging for paths
	log.Printf("Remote path: %s", remotePath)
	log.Printf("Processed path: %s", processedPath)

	// Expand ~ to home directory in privateKeyPath if needed
	if len(privateKeyPath) > 0 && privateKeyPath[0] == '~' {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal("Error getting home directory:", err)
		}
		privateKeyPath = filepath.Join(homeDir, privateKeyPath[1:])
	}

	// Start API server in a separate goroutine
	go startAPIServer()

	// Start the main processing loop
	go processLoop(sftpServer, username, privateKeyPath, remotePath, processedPath)

	// Block indefinitely (until shutdown)
	<-shutdownTrigger
	log.Println("Service shutting down gracefully...")
}