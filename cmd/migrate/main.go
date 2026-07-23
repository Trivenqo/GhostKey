package main

import (
	"errors"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	dbURL := os.Getenv("GHOSTKEY_DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/ghostkey?sslmode=disable"
	}

	log.Printf("Connecting to database: %s", dbURL)

	// "file://migrations" points to the migrations directory
	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		log.Fatalf("Failed to initialize migrator: %v", err)
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("Database is already up to date!")
			return
		}
		log.Fatalf("Migration failed: %v", err)
	}

	log.Println("Successfully applied database migrations!")
}