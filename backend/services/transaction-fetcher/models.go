package main

import (
	"time"
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