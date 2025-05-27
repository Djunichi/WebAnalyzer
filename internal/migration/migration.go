package migration

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// Direction - points to migration direction.
type Direction string

const (
	// Up - migrate to the latest version
	Up Direction = "up"
	// StepBack - rollback to previous version
	StepBack Direction = "step_back"

	stepBack = -1
)

type Config struct {
	Direction Direction `mapstructure:"direction"`
	Version   int       `mapstructure:"version"`
}

// PostgresMigrate runs postgres migrations. Returns version was migrated to
func PostgresMigrate(connectionStr string, cnf Config, migrationFiles embed.FS) (int, error) {

	db, err := sql.Open("postgres", connectionStr)
	if err != nil {
		return 0, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return 0, fmt.Errorf("failed to init postgres driver: %w", err)
	}

	// use preloaded migration files
	d, err := iofs.New(migrationFiles, "postgres")
	if err != nil {
		return 0, fmt.Errorf("failed to create postgres migrator source: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", d, "postgres", driver)
	if err != nil {
		return 0, fmt.Errorf("failed to create postgres migrator instance: %w", err)
	}

	switch cnf.Direction {
	case Up:
		err = m.Up()
		if err != nil {
			if !errors.Is(err, migrate.ErrNoChange) {
				return 0, fmt.Errorf("failed to run postgres migrations: %w", err)
			}
		}
	case StepBack:
		err = m.Steps(stepBack)
		if err != nil {
			return 0, fmt.Errorf("failed to run postgres migrations: %w", err)
		}
	default:
		return 0, fmt.Errorf("invalid migration config diretion: %s", cnf.Direction)
	}

	v, _, err := m.Version()
	if err != nil {
		return 0, fmt.Errorf(`migration went well, but version fetch failed.
								If it was down to 0 migration then this error is fine, err: %w`, err)
	}

	return int(v), nil
}
