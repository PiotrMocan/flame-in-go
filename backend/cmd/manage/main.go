package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mokan/flame-crm-backend/internal/db"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using OS env vars or defaults")
	}

	action := flag.String("action", "", "Action to perform: createdb, seed")
	flag.Parse()

	cmd := *action
	if cmd == "" {
		if len(os.Args) > 1 {
			cmd = os.Args[1]
		}
	}

	if cmd == "" {
		fmt.Println("Usage: go run cmd/manage/main.go [createdb|seed]")
		return
	}

	switch cmd {
	case "createdb":
		createDB()
	case "seed":
		db.ConnectDatabase()
		db.Seed(db.DB)
	default:
		fmt.Printf("Unknown action: %s\n", cmd)
		fmt.Println("Available actions: createdb, seed")
	}
}

func createDB() {
	targetDBName := os.Getenv("DB_NAME")
	if targetDBName == "" {
		targetDBName = "flame_crm"
	}

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	sslmode := os.Getenv("DB_SSLMODE")
	timezone := os.Getenv("DB_TIMEZONE")

	if host == "" {
		host = "localhost"
	}
	if user == "" {
		user = "postgres"
	}
	if password == "" {
		password = "postgres"
	}
	if port == "" {
		port = "5432"
	}
	if sslmode == "" {
		sslmode = "disable"
	}
	if timezone == "" {
		timezone = "UTC"
	}

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=postgres port=%s sslmode=%s TimeZone=%s",
		host, user, password, port, sslmode, timezone)

	maintenanceDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database instance (postgres):", err)
	}

	var exists bool
	checkSQL := fmt.Sprintf("SELECT EXISTS(SELECT datname FROM pg_catalog.pg_database WHERE datname = '%s')", targetDBName)
	if err := maintenanceDB.Raw(checkSQL).Scan(&exists).Error; err != nil {
		log.Fatal("Failed to check if database exists:", err)
	}

	if exists {
		fmt.Printf("Database '%s' already exists.\n", targetDBName)
	} else {
		fmt.Printf("Creating database '%s'...\n", targetDBName)
		if err := maintenanceDB.Exec(fmt.Sprintf("CREATE DATABASE %s", targetDBName)).Error; err != nil {
			log.Fatal("Failed to create database:", err)
		}
		fmt.Printf("Database '%s' created successfully.\n", targetDBName)
	}
}
