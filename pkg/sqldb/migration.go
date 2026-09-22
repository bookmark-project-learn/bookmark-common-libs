package sqldb

import (
	"fmt"

	"github.com/rs/zerolog/log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

func CatchDBError(err error) {
	if err != nil {
		log.Error().Err(err).Str("db error", err.Error()).Msg("DB Error")
		panic(err)
	}
}

type Migrator struct {
	migrator *migrate.Migrate
}

// MigrateLog

type MigrateLog struct{}

func NewMigrationLog() *MigrateLog {
	return &MigrateLog{}
}

// MigrateLog Printf
func (m *MigrateLog) Printf(format string, v ...interface{}) {
	log.Info().Msgf(format, v...)
}

// MigrateLog Verbose
func (m *MigrateLog) Verbose() bool {
	return true
}

// end log

// BuildMigrate return migrate.Migrate
func BuildMigrate(db *gorm.DB, migrationPath string) *Migrator {
	sqlDb, err := db.DB()
	CatchDBError(err)
	pgDriver, err := postgres.WithInstance(sqlDb, &postgres.Config{})
	CatchDBError(err)
	migrator, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%v", migrationPath),
		db.Name(), pgDriver)
	CatchDBError(err)
	return &Migrator{migrator}
}
func (m *Migrator) MigrateVersion() (uint, bool, error) {
	return m.migrator.Version()
}

func (m *Migrator) SetLogging() {
	m.migrator.Log = NewMigrationLog()
}

// MigrateUp: no step is migrate all
func (m *Migrator) MigrateUp(step ...int) error {
	if len(step) == 0 {
		err := m.migrator.Up()
		if err == migrate.ErrNoChange {
			return nil
		}
		return err
	}
	defaultStep := 1
	if len(step) > 0 {
		defaultStep = step[0]
	}
	return m.migrator.Steps(defaultStep)
}

// MigrateDown: no step is migrate single one
func (m *Migrator) MigrateDown(step ...int) error {
	defaultStep := -1
	if len(step) > 0 && step[0] < -1 {
		defaultStep = step[0]
	}
	return m.migrator.Steps(defaultStep)
}
