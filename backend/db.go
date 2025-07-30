package main

import (
	"os"
	"context"
	"fmt"
	"log"
	"time"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	
)

type DatabaseConfig struct {
	Host string
	Port string
	User string
	Password string
	DBName string
	SSLMode string
}

func getEnvOrDefault(k string, d string) string {
	if value := os.Getenv(k); value != "" { 
		return value
		}
	return d
}

func GetConfigFromEnv() *DatabaseConfig {
	return &DatabaseConfig{
		Host: getEnvOrDefault("PGHOST", "db"),
		Port: getEnvOrDefault("PGPORT", "5432"),
		User: getEnvOrDefault("PGUSER", "user"),
		Password: getEnvOrDefault("PGPASSWORD", "password"),
		DBName: getEnvOrDefault("PGDATABASE", "sanbar"),
		SSLMode: getEnvOrDefault("PGSSLMODE", "disable"),
	}
}


func (c *DatabaseConfig) BuildConnectionString(dbName string) string {
	if dbName == "" {
		dbName = c.DBName
	}
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s", c.User, c.Password, c.Host, c.Port, c.DBName, c.SSLMode)
}

type DatabaseManager struct {
	config *DatabaseConfig
	pool *pgxpool.Pool
}
	
func NewDatabaseManager(c *DatabaseConfig) *DatabaseManager {
	return &DatabaseManager{
		config: c,
	}
}

func (dm *DatabaseManager) Initialize(ctx context.Context) error {
	fmt.Println("start init")
	if err := dm.ensureDatabaseExists(ctx); err != nil {
		return fmt.Errorf("failed to ensure database exists: %w", err)
	}
	fmt.Println("createpool")

	if err := dm.createConnectionPool(ctx); err != nil {
		return fmt.Errorf("failed to create connection to connection pool: %w", err)
	}	
	fmt.Println("migrate")

	if err := dm.runMigrations(ctx); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}	

	return nil
}

func (dm *DatabaseManager) ensureDatabaseExists(ctx context.Context) error {
	connString := dm.config.BuildConnectionString("")
	conn, err := pgx.Connect(ctx, connString)
	if err != nil {
		return fmt.Errorf("failed to connect to psql database: %s", connString)
	}
	defer conn.Close(ctx)

	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	err = conn.QueryRow(ctx, query, dm.config.DBName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if database exists: %w", err)
	}

	if !exists {
		// Create the database (Note: database names cannot be parameterized)
		createQuery := fmt.Sprintf("CREATE DATABASE %s", pgx.Identifier{dm.config.DBName}.Sanitize())
		_, err = conn.Exec(ctx, createQuery)
		if err != nil {
			return fmt.Errorf("failed to create database: %w", err)
		}
		log.Printf("Database '%s' created successfully", dm.config.DBName)
	} else {
		log.Printf("Database '%s' already exists", dm.config.DBName)
	}
	fmt.Println("ensuredatabaseexists past")		
	return nil
}

func (dm *DatabaseManager) createConnectionPool(ctx context.Context) error {
	connString := dm.config.BuildConnectionString("")
	
	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Configure connection pool settings
	config.MaxConns = 30
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Test the connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return fmt.Errorf("failed to ping database: %w", err)
	}

	dm.pool = pool
	log.Println("Database connection pool created successfully")
	return nil
}

// runMigrations creates tables and other schema objects idempotently
func (dm *DatabaseManager) runMigrations(ctx context.Context) error {
	migrations := []string{
		// Enable extensions
		`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`,
		
		// Create users table
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			email VARCHAR(255) UNIQUE NOT NULL,
			username VARCHAR(100) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		
		// Create index on email for faster lookups
		`CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`,
		
		// Create posts table
		`CREATE TABLE IF NOT EXISTS posts (
			id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			title VARCHAR(255) NOT NULL,
			content TEXT,
			published BOOLEAN DEFAULT FALSE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		)`,
		
		// Create index on user_id for faster lookups
		`CREATE INDEX IF NOT EXISTS idx_posts_user_id ON posts(user_id)`,
		
		// Create index on published posts
		`CREATE INDEX IF NOT EXISTS idx_posts_published ON posts(published) WHERE published = true`,
	}

	tx, err := dm.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx) // This will be ignored if tx.Commit() succeeds

	for i, migration := range migrations {
		_, err := tx.Exec(ctx, migration)
		if err != nil {
			return fmt.Errorf("failed to execute migration %d: %w", i+1, err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit migrations: %w", err)
	}

	log.Println("All migrations completed successfully")
	return nil
}

// GetPool returns the connection pool for use by the application
func (dm *DatabaseManager) GetPool() *pgxpool.Pool {
	return dm.pool
}

// Close closes the database connection pool
func (dm *DatabaseManager) Close() {
	if dm.pool != nil {
		dm.pool.Close()
		log.Println("Database connection pool closed")
	}
}

