package main

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lucho2027/workout-tracker/apps/api/internal/config"
)

func main() {
	cfg := config.Load()
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	if _, err := pool.Exec(ctx, `
CREATE TABLE IF NOT EXISTS schema_migrations (
  version TEXT PRIMARY KEY,
  applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
)`); err != nil {
		log.Fatalf("create schema_migrations: %v", err)
	}

	files, err := collectSQL("migrations")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		var exists bool
		if err := pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version=$1)`, file).Scan(&exists); err != nil {
			log.Fatalf("check migration %s: %v", file, err)
		}
		if exists {
			fmt.Printf("skip %s\n", file)
			continue
		}

		content, err := os.ReadFile(filepath.Join("migrations", file))
		if err != nil {
			log.Fatalf("read migration %s: %v", file, err)
		}

		tx, err := pool.Begin(ctx)
		if err != nil {
			log.Fatalf("begin tx for %s: %v", file, err)
		}
		if _, err := tx.Exec(ctx, string(content)); err != nil {
			_ = tx.Rollback(ctx)
			log.Fatalf("apply %s: %v", file, err)
		}
		if _, err := tx.Exec(ctx, `INSERT INTO schema_migrations(version) VALUES ($1)`, file); err != nil {
			_ = tx.Rollback(ctx)
			log.Fatalf("record %s: %v", file, err)
		}
		if err := tx.Commit(ctx); err != nil {
			log.Fatalf("commit %s: %v", file, err)
		}

		fmt.Printf("applied %s\n", file)
	}

	fmt.Println("migrations complete")
}

func collectSQL(dir string) ([]string, error) {
	out := make([]string, 0)
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".sql" {
			out = append(out, filepath.Base(path))
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(out)
	return out, nil
}
