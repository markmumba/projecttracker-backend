package database

import (
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

var DB *gorm.DB

func ConnectDB() {
	var err error

	// Check if we're in a production environment
	if os.Getenv("GO_ENV") != "production" {
		// Load .env file only in non-production environments
		err = godotenv.Load()
		if err != nil {
			log.Println("Warning: Error loading .env file:", err)
			// Don't fatal here, as the environment variables might be set another way
		}
	}

	// Get the DATABASE_URL from environment variable
	databaseURL := os.Getenv("DATABASE_URL")

	if databaseURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// Open a connection to the database
	DB, err = gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	log.Println("Database connection established successfully")
}