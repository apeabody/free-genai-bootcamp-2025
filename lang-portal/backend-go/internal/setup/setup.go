package setup

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func Setup() {
	fmt.Println("Setting up database...")
	if err := initialize(); err != nil {
		fmt.Printf("Error initializing database: %v\n", err)
		os.Exit(1)
	}
	if err := seed(); err != nil {
		fmt.Printf("Error seeding database: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Setup complete!")
}

func initialize() error {
	fmt.Println("Initializing database...")

	// Create database file
	db, err := sql.Open("sqlite3", "words.db")
	if err != nil {
		return fmt.Errorf("failed to create database: %v", err)
	}
	defer db.Close()

	// Read and execute migration
	migrationPath := filepath.Join("internal", "db", "migrations", "001_create_tables.up.sql")
	migration, err := os.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %v", err)
	}

	// Split migration into individual statements
	statements := strings.Split(string(migration), ";")

	// Execute each statement
	for _, stmt := range statements {
		if strings.TrimSpace(stmt) != "" {
			if _, err := db.Exec(stmt); err != nil {
				return fmt.Errorf("failed to execute migration: %v", err)
			}
		}
	}

	fmt.Println("Database initialized successfully")
	return nil
}

func seed() error {
	fmt.Println("Seeding database...")

	db, err := sql.Open("sqlite3", "words.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	// Read and execute seed file
	seedPath := filepath.Join("internal", "db", "seeds", "initial_data.sql")
	seed, err := os.ReadFile(seedPath)
	if err != nil {
		return fmt.Errorf("failed to read seed file: %v", err)
	}

	// Split seed into individual statements
	statements := strings.Split(string(seed), ";")

	// Execute each statement
	for _, stmt := range statements {
		if strings.TrimSpace(stmt) != "" {
			if _, err := db.Exec(stmt); err != nil {
				return fmt.Errorf("failed to execute seed: %v", err)
			}
		}
	}

	fmt.Println("Database seeded successfully")
	return nil
}
