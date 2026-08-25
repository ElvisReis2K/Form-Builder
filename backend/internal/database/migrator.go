package database

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strconv"
	"strings"

	migrationfiles "github.com/ElvisReis2K/Form-Builder/backend/migrations"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Migration struct {
	Version int64
	Name    string
	UpPath  string
	DownPath string
	UpSQL   string
	DownSQL string
}

func MigrateUp(ctx context.Context, pool *pgxpool.Pool) error {
	if err := ensureSchemaMigrations(ctx, pool); err != nil {
		return err
	}

	applied, err := appliedVersions(ctx, pool)
	if err != nil {
		return err
	}

	migrations, err := loadMigrations()
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		if applied[migration.Version] {
			continue
		}

		if err := runMigration(ctx, pool, migration, migration.UpSQL, true); err != nil {
			return err
		}
	}

	return nil
}

func MigrateDown(ctx context.Context, pool *pgxpool.Pool) error {
	if err := ensureSchemaMigrations(ctx, pool); err != nil {
		return err
	}

	var version int64
	var name string
	err := pool.QueryRow(ctx, `
		SELECT version, name
		FROM schema_migrations
		ORDER BY version DESC
		LIMIT 1
	`).Scan(&version, &name)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("load latest migration: %w", err)
	}

	migrations, err := loadMigrations()
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		if migration.Version == version {
			return runMigration(ctx, pool, migration, migration.DownSQL, false)
		}
	}

	return fmt.Errorf("down migration not found for version %d (%s)", version, name)
}

func ensureSchemaMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version bigint PRIMARY KEY,
			name text NOT NULL,
			applied_at timestamptz NOT NULL DEFAULT now()
		)
	`)
	if err != nil {
		return fmt.Errorf("ensure schema_migrations: %w", err)
	}

	return nil
}

func appliedVersions(ctx context.Context, pool *pgxpool.Pool) (map[int64]bool, error) {
	rows, err := pool.Query(ctx, "SELECT version FROM schema_migrations")
	if err != nil {
		return nil, fmt.Errorf("load applied migrations: %w", err)
	}
	defer rows.Close()

	applied := map[int64]bool{}
	for rows.Next() {
		var version int64
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("scan applied migration: %w", err)
		}

		applied[version] = true
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate applied migrations: %w", err)
	}

	return applied, nil
}

func loadMigrations() ([]Migration, error) {
	upPaths, err := fs.Glob(migrationfiles.Files, "*.up.sql")
	if err != nil {
		return nil, fmt.Errorf("list migrations: %w", err)
	}
	sort.Strings(upPaths)

	migrations := make([]Migration, 0, len(upPaths))
	for _, upPath := range upPaths {
		version, name, err := parseMigrationName(upPath)
		if err != nil {
			return nil, err
		}

		downPath := strings.TrimSuffix(upPath, ".up.sql") + ".down.sql"
		upSQL, err := fs.ReadFile(migrationfiles.Files, upPath)
		if err != nil {
			return nil, fmt.Errorf("read up migration %s: %w", upPath, err)
		}

		downSQL, err := fs.ReadFile(migrationfiles.Files, downPath)
		if err != nil {
			return nil, fmt.Errorf("read down migration %s: %w", downPath, err)
		}

		migrations = append(migrations, Migration{
			Version:  version,
			Name:     name,
			UpPath:   upPath,
			DownPath: downPath,
			UpSQL:    string(upSQL),
			DownSQL:  string(downSQL),
		})
	}

	return migrations, nil
}

func parseMigrationName(filePath string) (int64, string, error) {
	base := path.Base(filePath)
	versionText, nameWithSuffix, ok := strings.Cut(base, "_")
	if !ok {
		return 0, "", fmt.Errorf("invalid migration name %s", base)
	}

	version, err := strconv.ParseInt(versionText, 10, 64)
	if err != nil {
		return 0, "", fmt.Errorf("invalid migration version %s: %w", base, err)
	}

	name := strings.TrimSuffix(nameWithSuffix, ".up.sql")
	return version, name, nil
}

func runMigration(ctx context.Context, pool *pgxpool.Pool, migration Migration, script string, up bool) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin migration transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback(ctx)
	}()

	if err := execScript(ctx, tx, script); err != nil {
		return fmt.Errorf("run migration %d %s: %w", migration.Version, migration.Name, err)
	}

	if up {
		_, err = tx.Exec(ctx, `
			INSERT INTO schema_migrations (version, name)
			VALUES ($1, $2)
		`, migration.Version, migration.Name)
	} else {
		_, err = tx.Exec(ctx, "DELETE FROM schema_migrations WHERE version = $1", migration.Version)
	}
	if err != nil {
		return fmt.Errorf("record migration %d: %w", migration.Version, err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit migration %d: %w", migration.Version, err)
	}

	return nil
}

func execScript(ctx context.Context, tx pgx.Tx, script string) error {
	for index, statement := range splitStatements(script) {
		if _, err := tx.Exec(ctx, statement); err != nil {
			return fmt.Errorf("statement %d: %w", index+1, err)
		}
	}

	return nil
}

func splitStatements(script string) []string {
	// Migration files in this project use plain PostgreSQL statements.
	// Keep this runner small; move to a dedicated migration library if future
	// migrations need procedural SQL with semicolons inside function bodies.
	rawStatements := strings.Split(script, ";")
	statements := make([]string, 0, len(rawStatements))

	for _, statement := range rawStatements {
		statement = strings.TrimSpace(statement)
		if statement == "" {
			continue
		}

		statements = append(statements, statement)
	}

	return statements
}
