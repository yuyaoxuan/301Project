package main

import (
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	_"github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/joho/godotenv"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// TransactionLog represents your transaction log structure for CSV data
type TransactionLog struct {
	ID          int       `json:"id"`
	ClientID    string    `json:"clientid"`
	Transaction string    `json:"transaction"`
	Amount      float64   `json:"amount"`
	Date        time.Time `json:"date"`
	Status      string    `json:"status"`
}

// ServiceStatus represents the current status of the log processor
type ServiceStatus struct {
	Running          bool      `json:"running"`
	LastCheck        time.Time `json:"last_check"`
	CycleCount       int       `json:"cycle_count"`
	FilesProcessed   int       `json:"files_processed"`
	ProcessingErrors int       `json:"processing_errors"`
}

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

func startAPIServer() {
	// Define API endpoints
	http.HandleFunc("/status", getStatusHandler)
	http.HandleFunc("/trigger", triggerProcessingHandler)
	http.HandleFunc("/shutdown", shutdownHandler)
	http.HandleFunc("/transactions/", getTransactionsHandler)

	// Get port from env or use default
	port := os.Getenv("API_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting API server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start API server: %v", err)
	}
}

func getStatusHandler(w http.ResponseWriter, r *http.Request) {
	statusMutex.Lock()
	defer statusMutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

func triggerProcessingHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Send trigger signal (non-blocking)
	select {
	case processTrigger <- true:
		log.Println("Manual processing triggered via API")
	default:
		// Channel is full, which means processing is already triggered
		log.Println("Processing already triggered")
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Processing triggered",
	})
}

func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Service is shutting down",
	})

	// Trigger shutdown after sending response
	go func() {
		time.Sleep(100 * time.Millisecond)
		shutdownTrigger <- true
	}()
}

func getTransactionsHandler(w http.ResponseWriter, r *http.Request) {
	// Only accept GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract clientid from the URL path
	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) < 3 {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}
	clientID := pathParts[2] // The clientid is the third part of the path

	// Connect to MySQL database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("Failed to connect to database:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Query transactions for the given clientid
	query := `
		SELECT id, clientid, transaction_type, amount, transaction_date, status
		FROM transaction_logs
		WHERE clientid = ?
		ORDER BY transaction_date DESC
	`
	rows, err := db.Query(query, clientID)
	if err != nil {
		log.Println("Failed to query transactions:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Collect results
	var transactions []TransactionLog
	for rows.Next() {
		var t TransactionLog
		err := rows.Scan(
			&t.ID,
			&t.ClientID,
			&t.Transaction,
			&t.Amount,
			&t.Date,
			&t.Status,
		)
		if err != nil {
			log.Println("Failed to scan transaction row:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		transactions = append(transactions, t)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		log.Println("Error iterating over transaction rows:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check if no transactions were found
	if len(transactions) == 0 {
		http.Error(w, "No transactions found for "+clientID, http.StatusNotFound)
		return
	}

	// Return results as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}

func ensureDatabaseExists(db *sql.DB, dbName string) error {
	_, err := db.Exec(fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s", dbName))
	if err != nil {
		return fmt.Errorf("failed to create database: %v", err)
	}
	return nil
}

func ensureTableExists(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS transaction_logs (
			id INT PRIMARY KEY,
			clientid VARCHAR(255) NOT NULL,
			transaction_type VARCHAR(255) NOT NULL,
			amount FLOAT(10, 2) NOT NULL,
			transaction_date DATETIME NOT NULL,
			status VARCHAR(255) NOT NULL,
			UNIQUE KEY (id)
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %v", err)
	}
	return nil
}

func initializeDatabase() (*sql.DB, error) {
	// Connect to MySQL without specifying a database
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT")))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %v", err)
	}
	defer db.Close()

	// Ensure the database exists
	err = ensureDatabaseExists(db, os.Getenv("DB_NAME"))
	if err != nil {
		return nil, fmt.Errorf("failed to ensure database exists: %v", err)
	}

	// Reconnect to the specific database
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME")))
	if err != nil {
		return nil, fmt.Errorf("failed to reconnect to database: %v", err)
	}

	// Ensure the table exists
	err = ensureTableExists(db)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure table exists: %v", err)
	}

	return db, nil
}

func processLoop(sftpServer, username, privateKeyPath, remotePath, processedPath string) {
	tickerInterval := 5 * time.Minute
	ticker := time.NewTicker(tickerInterval)
	defer ticker.Stop()

	// Process once immediately on startup
	processOnce(sftpServer, username, privateKeyPath, remotePath, processedPath)

	for {
		select {
		case <-ticker.C:
			// Process on ticker schedule
			log.Println("Scheduled processing triggered")
			processOnce(sftpServer, username, privateKeyPath, remotePath, processedPath)

		case <-processTrigger:
			// Process on manual trigger
			log.Println("Manual processing triggered")
			processOnce(sftpServer, username, privateKeyPath, remotePath, processedPath)

			// Reset the ticker to avoid processing twice in quick succession
			ticker.Reset(tickerInterval)
		}
	}
}

func processOnce(sftpServer, username, privateKeyPath, remotePath, processedPath string) {
	// Update status
	statusMutex.Lock()
	status.Running = true
	status.LastCheck = time.Now()
	status.CycleCount++
	statusMutex.Unlock()

	log.Printf("Starting log processing cycle %d...", status.CycleCount)

	// Load SSH private key
	key, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		log.Println("Unable to read private key:", err)
		updateStatusError()
		return
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Println("Unable to parse private key:", err)
		updateStatusError()
		return
	}

	// SSH Config
	config := &ssh.ClientConfig{
		User:            username,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         30 * time.Second,
	}

	// Connect to SFTP Server
	client, err := ssh.Dial("tcp", sftpServer+":22", config)
	if err != nil {
		log.Println("Failed to connect:", err)
		updateStatusError()
		return
	}
	defer client.Close()

	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		log.Println("Failed to create SFTP client:", err)
		updateStatusError()
		return
	}
	defer sftpClient.Close()

	// Connect to MySQL database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"))

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("Failed to connect to database:", err)
		updateStatusError()
		return
	}
	defer db.Close()

	// Set connection pool parameters
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test database connection
	err = db.Ping()
	if err != nil {
		log.Println("Database connection failed:", err)
		updateStatusError()
		return
	}

	// Create temp directory for downloaded files
	tempDir, err := ioutil.TempDir("", "sftp-logs")
	if err != nil {
		log.Println("Failed to create temp directory:", err)
		updateStatusError()
		return
	}
	defer os.RemoveAll(tempDir)

	// Ensure processed directory exists
	err = ensureRemoteDir(sftpClient, processedPath)
	if err != nil {
		log.Println("Failed to ensure processed directory exists:", err)
		updateStatusError()
		return
	}

	// List client directories
	clientDirs, err := sftpClient.ReadDir(remotePath)
	if err != nil {
		log.Println("Failed to list directory:", err)
		updateStatusError()
		return
	}

	filesProcessed := 0
	errors := 0

	// Process each client directory
	for _, clientDir := range clientDirs {
		// Skip files and hidden directories
		if !clientDir.IsDir() || strings.HasPrefix(clientDir.Name(), ".") {
			continue
		}

		clientID := clientDir.Name()
		// Skip the "processed" directory itself
		if clientID == "processed" {
			continue
		}
		clientPath := filepath.Join(remotePath, clientID)
		log.Printf("Processing client directory: %s", clientID)
		log.Printf("Client path: %s", clientPath)

		// Ensure processed directory exists for this client
		clientProcessedPath := filepath.Join(processedPath, clientID)
		log.Printf("Client processed path: %s", clientProcessedPath)

		err = ensureRemoteDir(sftpClient, clientProcessedPath)
		if err != nil {
			log.Printf("Warning: Could not create processed directory for client %s: %v", clientID, err)
			errors++
			continue
		}

		// List files in the client directory
		files, err := sftpClient.ReadDir(clientPath)
		if err != nil {
			log.Printf("Failed to list files for client %s: %v", clientID, err)
			errors++
			continue
		}

		// Process each log file
		for _, file := range files {
			// Skip directories and hidden files
			if file.IsDir() || strings.HasPrefix(file.Name(), ".") {
				continue
			}

			// Only process CSV files
			if !strings.HasSuffix(strings.ToLower(file.Name()), ".csv") {
				continue
			}

			log.Printf("Processing file: %s/%s", clientID, file.Name())

			// Full path to remote file
			remoteFilePath := filepath.Join(clientPath, file.Name())
			localFilePath := filepath.Join(tempDir, clientID+"_"+file.Name())

			// Download the file
			err = downloadFile(sftpClient, remoteFilePath, localFilePath)
			if err != nil {
				log.Printf("Error downloading file %s/%s: %v", clientID, file.Name(), err)
				errors++
				continue
			}

			// Process the downloaded file
			err = processLogFile(localFilePath, db)
			if err != nil {
				log.Printf("Error processing file %s/%s: %v", clientID, file.Name(), err)
				errors++
				continue
			}

			// Move to processed directory
			processedFilePath := filepath.Join(clientProcessedPath, file.Name())
			err = moveRemoteFile(sftpClient, remoteFilePath, processedFilePath)
			if err != nil {
				log.Printf("Warning: Could not move processed file %s/%s: %v", clientID, file.Name(), err)
				errors++
			}

			log.Printf("Successfully processed %s/%s", clientID, file.Name())
			filesProcessed++
		}
	}

	// Update status
	statusMutex.Lock()
	status.Running = false
	status.FilesProcessed += filesProcessed
	status.ProcessingErrors += errors
	statusMutex.Unlock()

	log.Printf("Processing cycle %d completed. Files processed: %d, Errors: %d", status.CycleCount, filesProcessed, errors)
}

func updateStatusError() {
	statusMutex.Lock()
	status.Running = false
	status.ProcessingErrors++
	statusMutex.Unlock()
}

// downloadFile downloads a file from SFTP server to local path
func downloadFile(client *sftp.Client, remotePath, localPath string) error {
	srcFile, err := client.Open(remotePath)
	if err != nil {
		return fmt.Errorf("failed to open remote file: %v", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("failed to create local file: %v", err)
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return fmt.Errorf("failed to copy file contents: %v", err)
	}

	return nil
}

// moveRemoteFile moves a file on the remote SFTP server
func moveRemoteFile(client *sftp.Client, oldPath, newPath string) error {
	return client.Rename(oldPath, newPath)
}

// ensureRemoteDir ensures a directory exists on the remote server
func ensureRemoteDir(client *sftp.Client, path string) error {
	_, err := client.Stat(path)
	if err == nil {
		// Directory exists
		return nil
	}

	// Create directory if it doesn't exist
	return client.MkdirAll(path)
}

// processLogFile parses the CSV log file and stores entries in the database
func processLogFile(filePath string, db *sql.DB) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	// Create CSV reader
	reader := csv.NewReader(file)

	// Read header row
	header, err := reader.Read()
	if err != nil {
		return fmt.Errorf("failed to read CSV header: %v", err)
	}

	// Verify header format
	expectedHeaders := []string{"ID", "ClientID", "Transaction", "Amount", "Date", "Status"}
	if !validateHeaders(header, expectedHeaders) {
		return fmt.Errorf("invalid CSV header format: %v", header)
	}

	// Begin transaction for batch inserts
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %v", err)
	}

	// Prepare statement for inserting logs - MySQL syntax
	stmt, err := tx.Prepare(`
		INSERT INTO transaction_logs (id, clientid, transaction_type, amount, transaction_date, status) 
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE 
		    clientid = VALUES(clientid),
		    transaction_type = VALUES(transaction_type),
		    amount = VALUES(amount),
		    transaction_date = VALUES(transaction_date),
		    status = VALUES(status)
	`)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to prepare statement: %v", err)
	}
	defer stmt.Close()

	// Process each row
	count := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error reading CSV row: %v", err)
			continue
		}

		// Parse the data
		id, err := strconv.Atoi(record[0])
		if err != nil {
			log.Printf("Invalid ID in row %d: %v", count+1, err)
			continue
		}

		amount, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Printf("Invalid Amount in row %d: %v", count+1, err)
			continue
		}

		date, err := time.Parse(time.RFC3339, record[4])
		if err != nil {
			// Try alternative date formats if RFC3339 fails
			date, err = tryParseDate(record[4])
			if err != nil {
				log.Printf("Invalid Date in row %d: %v", count+1, err)
				continue
			}
		}

		// Insert into database
		_, err = stmt.Exec(
			id,
			record[1], // ClientID
			record[2], // Transaction
			amount,
			date,
			record[5], // Status
		)
		if err != nil {
			log.Printf("Error inserting row %d: %v", count+1, err)
			continue
		}

		count++
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("Successfully processed %d transaction records from %s", count, filepath.Base(filePath))
	return nil
}

// validateHeaders checks if the CSV header matches the expected format
func validateHeaders(headers, expectedHeaders []string) bool {
	if len(headers) < len(expectedHeaders) {
		return false
	}

	for i, expected := range expectedHeaders {
		if i >= len(headers) || headers[i] != expected {
			return false
		}
	}

	return true
}

// tryParseDate attempts to parse a date string using multiple formats
func tryParseDate(dateStr string) (time.Time, error) {
	formats := []string{
		time.RFC3339,
		"2006-01-02T15:04:05",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}
