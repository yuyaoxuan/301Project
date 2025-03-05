package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

// Global DB connection
var DB *sql.DB

// LoadEnv loads the .env file
func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("❌ Error loading .env file")
	}
}

// ConnectDB initializes the database and creates tables if needed
func ConnectDB() {
	LoadEnv()

	// Get DB credentials from environment variables
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	// Connection string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/", dbUser, dbPassword, dbHost, dbPort)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ Error connecting to MySQL:", err)
	}

	// Verify connection
	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ Database ping failed:", err)
	}

	fmt.Println("✅ Connected to MySQL!")

	// Initialize database and tables
	initializeDatabase(dbName)
}

// initializeDatabase creates the database and tables if they don't exist
func initializeDatabase(dbName string) {
	// Create database if it doesn't exist
	_, err := DB.Exec("CREATE DATABASE IF NOT EXISTS " + dbName)
	if err != nil {
		log.Fatal("❌ Error creating database:", err)
	}

	fmt.Println("✅ Database checked/created!")

	// Close previous connection and reconnect to the new database
	DB.Close()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), dbName)
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("❌ Error reconnecting to database:", err)
	}

	// Ensure connection is valid
	err = DB.Ping()
	if err != nil {
		log.Fatal("❌ Database ping failed:", err)
	}
}
