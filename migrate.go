package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer pool.Close()

	if len(os.Args) > 1 && os.Args[1] == "clean" {
		cleanDatabase(pool)
		return
	}

	applyMigrations(pool)
}

func cleanDatabase(pool *pgxpool.Pool) {
	_, err := pool.Exec(context.Background(), "DROP SCHEMA public CASCADE; CREATE SCHEMA public;")
	if err != nil {
		log.Fatalf("Error cleaning database: %v", err)
	}
	log.Println("Database cleaned successfully!")
}

func applyMigrations(pool *pgxpool.Pool) {
	sqlFile := "./migrations/schema.sql"
	sqlBytes, err := os.ReadFile(sqlFile)
	if err != nil {
		log.Fatalf("Unable to read SQL file: %v", err)
	}

	queries := strings.Split(string(sqlBytes), ";")

	for _, query := range queries {
		query = strings.TrimSpace(query)
		if query == "" {
			continue
		}

		_, err := pool.Exec(context.Background(), query)
		if err != nil {
			log.Fatalf("Error executing query: %v\nQuery: %s", err, query)
		}
	}

	log.Println("Migrations applied successfully!")
}
