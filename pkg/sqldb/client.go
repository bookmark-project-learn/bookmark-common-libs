package sqldb

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// coreClient creates a new GORM client
func buildClient(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		panic(fmt.Sprintf("Failed to start DB: %v", err))
	}
	return db, err
}

// getDSNConfig creates a new config struct based on the environment variables
func getClientFromConfig(cfg *config) (*gorm.DB, error) {
	db, err := buildClient(cfg.getDSN())
	return db, err
}

func NewSqlDB() (*gorm.DB, error) {
	cfg, err := getConfig()
	if err != nil {
		return nil, err
	}
	return getClientFromConfig(cfg)
}
