package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

)

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

	// Connect to SFTP server
	sshClient, sftpClient, err := connectSFTP(sftpServer, username, privateKeyPath)
	if err != nil {
		log.Println("SFTP connection error:", err)
		updateStatusError()
		return
	}
	defer sshClient.Close()
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
	clientDirs, err := listRemoteDir(sftpClient, remotePath)
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
		files, err := listRemoteDir(sftpClient, clientPath)
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