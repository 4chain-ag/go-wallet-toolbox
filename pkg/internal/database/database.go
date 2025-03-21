package database

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/logging"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// Database is a struct that holds logger for database connection and the connection itself
type Database struct {
	DB     *gorm.DB
	logger *slog.Logger
}

// NewDatabase will configure and return database based on provided config
func NewDatabase(cfg defs.Database, logger *slog.Logger) (*Database, error) {
	if logger == nil {
		logger = logging.Child(logger, "database")
	}

	gormLogger := &SlogGormLogger{
		logger: logger,
	}

	dialector, ok := dialectors[cfg.Engine]
	if !ok {
		return nil, fmt.Errorf("dialector for engine %s not found", cfg.Engine)
	}

	database, err := createAndConfigureDatabaseConnection(dialector(cfg), cfg, gormLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to create gorm instance, caused by: %w", err)
	}

	return &Database{
		DB:     database,
		logger: logger,
	}, nil
}

func createAndConfigureDatabaseConnection(dialector gorm.Dialector, cfg defs.Database, logger glogger.Interface) (*gorm.DB, error) {
	db, err := gorm.Open(dialector, createGormConfig(
		logger,
	))
	if err != nil {
		return nil, fmt.Errorf("failed to initialize GORM database connection: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve underlying SQL database connection: %w", err)

	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(cfg.MaxConnectionTime)
	sqlDB.SetConnMaxIdleTime(cfg.MaxConnectionIdleTime)

	return db, nil
}

// createGormConfig returns valid gorm.Config for database connections
func createGormConfig(logger glogger.Interface) *gorm.Config {
	// Set the prefix
	tablePrefix := defs.DefaultTablePrefix

	if logger == nil {
		panic("Could not create gorm config. When creating database configuration you need to specify the logger to use")
	}

	// Create the configuration
	gormCfg := &gorm.Config{
		AllowGlobalUpdate:                        false,
		ClauseBuilders:                           nil,
		ConnPool:                                 nil,
		CreateBatchSize:                          0,
		Dialector:                                nil,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: true,
		DisableNestedTransaction:                 false,
		DryRun:                                   false, // toggle for extreme debugging
		FullSaveAssociations:                     false,
		Logger:                                   logger,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   tablePrefix, // table name prefix, table for `User` would be `t_users`
			SingularTable: false,       // use singular table name, table for `User` would be `user` with this option enabled
		},
		NowFunc:                nil,
		Plugins:                nil,
		PrepareStmt:            false,
		QueryFields:            false,
		SkipDefaultTransaction: false,
		TranslateError:         true,
	}

	return gormCfg
}

// normalizeTimeZone changes every "/" in timezone to special char "%2F" for mysql to parse time location correctly
// https://github.com/go-sql-driver/mysql?tab=readme-ov-file#loc
func normalizeTimeZone(tz string) string {
	return strings.Replace(tz, "/", "%2F", -1)
}
