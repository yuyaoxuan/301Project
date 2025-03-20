package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

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

	// Return results as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(transactions)
}