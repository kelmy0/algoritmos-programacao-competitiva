package database

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"

	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

var DB *pgxpool.Pool

func ConnectDB(reset bool) bool {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Println("❌ DATABASE_URL not found.")
		return false
	}

	var err error
	DB, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Printf("❌ Error connecting to the database: %v\n", err)
		return false
	}

	// Try running the migrations at startup
	if !runMigrations(connStr, reset) {
		return false
	}

	return true
}

func runMigrations(connStr string, shouldReset bool) bool {
	migrateConnStr := strings.Replace(connStr, "postgres://", "pgx5://", 1)
	migrateConnStr = strings.Replace(migrateConnStr, "postgresql://", "pgx5://", 1)

	m, err := migrate.New("file://db/migrations", migrateConnStr)

	if err != nil {
		log.Printf("❌ Error initializing the migration manager: %v\n", err)
		return false
	}

	m.Force(1)
	if shouldReset {
		log.Println("🗑️ Dropping all tables and types...")
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Printf("⚠️ Warning during database clean: %v\n", err)
		}
		log.Println("✨ Database cleaned successfully!")
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Printf("❌ Error applying migrations: %v\n", err)
		return false
	}

	if err == migrate.ErrNoChange {
		log.Println("✨ The database is already up-to-date.")
	} else {
		log.Println("🚀 Migrations successfully implemented!")
	}

	return true
}

func RunSeeds() {
	files, err := os.ReadDir("db/seeds")
	if err != nil {
		log.Fatalf("❌ Error reading seeds directory: %v\n", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".sql") {
			filePath := filepath.Join("db/seeds", file.Name())

			query, err := os.ReadFile(filePath)
			if err != nil {
				log.Fatalf("❌ Error reading seed file %s: %v\n", file.Name(), err)
			}

			_, err = DB.Exec(context.Background(), string(query))
			if err != nil {
				log.Fatalf("❌ Error executing seed %s: %v\n", file.Name(), err)
			}
		}
	}

	log.Println("✨ Database successfully seeded!")
}
