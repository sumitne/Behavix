package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"behavix-ai/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const schemaMigrationsTable = "schema_migrations"

// RunMigrate runs pending migrations. migrationsDir is the path to the migrations folder (e.g. "migrations").
func RunMigrate(migrationsDir string) {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}
	if cfg.PostgresDSN == "" {
		fmt.Fprintf(os.Stderr, "set POSTGRES_DSN or DB_HOST, DB_USER, DB_PASSWORD, DB_NAME in .env\n")
		os.Exit(1)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.PostgresDSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "postgres connect: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()
	if err := pool.Ping(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "postgres ping: %v\n", err)
		os.Exit(1)
	}

	if err := ensureSchemaMigrationsTable(ctx, pool); err != nil {
		fmt.Fprintf(os.Stderr, "ensure schema_migrations: %v\n", err)
		os.Exit(1)
	}

	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read migrations dir: %v\n", err)
		os.Exit(1)
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		files = append(files, e.Name())
	}
	sort.Strings(files)

	for _, name := range files {
		version := name
		applied, err := isApplied(ctx, pool, version)
		if err != nil {
			fmt.Fprintf(os.Stderr, "check applied %s: %v\n", version, err)
			os.Exit(1)
		}
		if applied {
			continue
		}
		path := filepath.Join(migrationsDir, name)
		sql, err := os.ReadFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "read %s: %v\n", path, err)
			os.Exit(1)
		}
		if _, err := pool.Exec(ctx, string(sql)); err != nil {
			fmt.Fprintf(os.Stderr, "migrate %s: %v\n", version, err)
			os.Exit(1)
		}
		if _, err := pool.Exec(ctx,
			`INSERT INTO `+schemaMigrationsTable+` (version) VALUES ($1)`,
			version,
		); err != nil {
			fmt.Fprintf(os.Stderr, "record migration %s: %v\n", version, err)
			os.Exit(1)
		}
		fmt.Println("migrate-up:", version)
	}
	fmt.Println("migrate-up: done")
}

func ensureSchemaMigrationsTable(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS `+schemaMigrationsTable+` (
			version TEXT PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	return err
}

func isApplied(ctx context.Context, pool *pgxpool.Pool, version string) (bool, error) {
	var n int
	err := pool.QueryRow(ctx,
		`SELECT 1 FROM `+schemaMigrationsTable+` WHERE version = $1`,
		version,
	).Scan(&n)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
