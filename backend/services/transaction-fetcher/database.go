package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

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

	// Set connection pool parameters
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

// insertTransactionLogs inserts transaction logs into the database
func insertTransactionLogs(db *sql.DB, id int, clientID, transactionType string, amount float64, date time.Time, status string) error {
	query := `
		INSERT INTO transaction_logs (id, clientid, transaction_type, amount, transaction_date, status) 
		VALUES (?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE 
			clientid = VALUES(clientid),
			transaction_type = VALUES(transaction_type),
			amount = VALUES(amount),
			transaction_date = VALUES(transaction_date),
			status = VALUES(status)
	`
	_, err := db.Exec(query, id, clientID, transactionType, amount, date, status)
	return err
}