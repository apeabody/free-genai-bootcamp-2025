//go:build mage
// +build mage

package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// Default target to run when none is specified
var Default = Setup

// Setup runs Initialize and Seed
func Setup() {
	fmt.Println("Setting up database...")
	Initialize()
	Seed()
	fmt.Println("Setup complete!")
}

// Initialize creates a new SQLite database and initializes the schema
func Initialize() {
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

// Seed adds sample data to the database
func Seed() {
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

// Reset drops all tables and reinitializes the database
func Reset() {
	fmt.Println("Resetting database...")

	db, err := sql.Open("sqlite3", "words.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}
	defer db.Close()

	// Read and execute down migration
	downMigrationPath := filepath.Join("internal", "db", "migrations", "001_create_tables.down.sql")
	downMigration, err := os.ReadFile(downMigrationPath)
	if err != nil {
		return fmt.Errorf("failed to read down migration file: %v", err)
	}

	// Split migration into individual statements
	statements := strings.Split(string(downMigration), ";")

	// Execute each statement
	for _, stmt := range statements {
		if strings.TrimSpace(stmt) != "" {
			if _, err := db.Exec(stmt); err != nil {
				return fmt.Errorf("failed to execute down migration: %v", err)
			}
		}
	}

	// Reinitialize the database
	db := DB{}
	if err := db.Initialize(); err != nil {
		return fmt.Errorf("failed to reinitialize database: %v", err)
	}

	fmt.Println("Database reset successfully")
	return nil
}
