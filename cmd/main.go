package main

import (
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"os"

	"template/src/auth"
	"template/src/config"
	"template/src/handler"
	"template/src/store"

	"github.com/enrichman/httpgrace"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	c := config.NewConfig()

	if c.DropDatabase {
		os.Remove(c.DatabaseURL)
	}

	db, err := sqlx.Open("sqlite3", c.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	migration, err := createMigrationInstance(db.DB, c.RootPath)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	err = runMigration(logger, migration)
	if err != nil {
		log.Fatalf("Failed to run migration: %v", err)
	}

	userStore := store.NewUserStore(db)
	store := store.NewStore(db)
	authHandler := auth.NewAuthHandler(c.PrivateKey, c.PublicKey)
	server := handler.NewServer(userStore, store, authHandler, logger)

	err = httpgrace.ListenAndServe(c.Address+":"+c.Port, server)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func runMigration(logger *slog.Logger, migration *migrate.Migrate) error {
	oldVersion, _, _ := migration.Version()

	err := migration.Up()
	if err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to run migration: %w", err)
	} else if err == migrate.ErrNoChange {
		logger.Info("no migration changes no problem", "version", oldVersion)
		return nil
	}

	newVersion, _, err := migration.Version()
	if err != nil {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	logger.Info("migrated", "from", oldVersion, "to", newVersion)

	return nil
}

func createMigrationInstance(db *sql.DB, rootPath string) (*migrate.Migrate, error) {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to create driver: %w", err)
	}

	migration, err := migrate.NewWithDatabaseInstance(
		"file://"+rootPath+"/src/migration",
		"sqlite3", driver)
	if err != nil {
		return nil, fmt.Errorf("failed to create migration instance: %w", err)
	}

	return migration, nil
}
