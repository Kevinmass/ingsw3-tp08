package integration

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testDB *sql.DB
var postgresContainer *postgres.PostgresContainer

// SetupTestDB starts a PostgreSQL container and initializes the database
func SetupTestDB() (*sql.DB, func(), error) {
	ctx := context.Background()

	pgContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("postgres:15.3-alpine"),
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("testuser"),
		postgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(30)),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start postgres container: %w", err)
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get connection string: %w", err)
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create tables
	if err := createTables(db); err != nil {
		db.Close()
		return nil, nil, fmt.Errorf("failed to create tables: %w", err)
	}

	cleanup := func() {
		db.Close()
		if err := testcontainers.TerminateContainer(pgContainer); err != nil {
			log.Printf("failed to terminate container: %v", err)
		}
	}

	return db, cleanup, nil
}

// createTables creates the necessary tables for testing
func createTables(db *sql.DB) error {
	// Create users table
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) UNIQUE NOT NULL,
		password VARCHAR(255) NOT NULL,
		username VARCHAR(255) NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);`

	if _, err := db.Exec(usersTable); err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}

	// Create posts table
	postsTable := `
	CREATE TABLE IF NOT EXISTS posts (
		id SERIAL PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		content TEXT NOT NULL,
		user_id INTEGER NOT NULL REFERENCES users(id),
		created_at TIMESTAMP DEFAULT NOW()
	);`

	if _, err := db.Exec(postsTable); err != nil {
		return fmt.Errorf("failed to create posts table: %w", err)
	}

	// Create comments table
	commentsTable := `
	CREATE TABLE IF NOT EXISTS comments (
		id SERIAL PRIMARY KEY,
		post_id INTEGER NOT NULL REFERENCES posts(id),
		user_id INTEGER NOT NULL REFERENCES users(id),
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT NOW()
	);`

	if _, err := db.Exec(commentsTable); err != nil {
		return fmt.Errorf("failed to create comments table: %w", err)
	}

	return nil
}

// CleanupTestDB truncates tables to clean state
func CleanupTestDB(db *sql.DB) error {
	tables := []string{"comments", "posts", "users"}
	for _, table := range tables {
		query := "TRUNCATE TABLE " + table + " CASCADE"
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}

	return nil
}
