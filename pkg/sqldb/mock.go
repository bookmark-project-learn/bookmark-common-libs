package sqldb

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewMiniPostgres(t *testing.T) (*gorm.DB, error) {
	dsn := ":memory:?cache=shared:"
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Failed to start DB")
	}

	return db, err
}
