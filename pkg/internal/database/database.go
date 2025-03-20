package database

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/internal/database/sqlite3extended"
	"gorm.io/driver/sqlite"
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
func NewDatabase(cfg *defs.Database, logger *slog.Logger) (*Database, error) {
	defaultConfig := defs.DefaultDBConfig()
	if cfg == nil {
		return newDatabaseInternal(defaultConfig, logger)
	}

	mergedCfg := mergeConfig(defaultConfig, cfg)
	return newDatabaseInternal(mergedCfg, logger)
}

// newDatabaseInternal configures database with merged default and provided config
func newDatabaseInternal(cfg *defs.Database, logger *slog.Logger) (*Database, error) {
	var database *gorm.DB

	if logger == nil {
		logger = slog.Default()
	}

	gormLogger := &SlogGormLogger{
		logger: logger,
	}

	switch cfg.Engine {
	case defs.DBTypeSQLite:
		db, err := openSQLiteDatabase(cfg, gormLogger)
		if err != nil {
			return nil, fmt.Errorf("failed to create gorm instance, caused by: %w", err)
		}

		database = db
	case defs.DBTypePostgres:
		db, err := openPostgresDatabase(cfg, gormLogger)
		if err != nil {
			return nil, fmt.Errorf("failed to create gorm instance, caused by: %w", err)
		}

		database = db
	case defs.DBTypeMySQL:
		db, err := openMySQLDatabase(cfg, gormLogger)
		if err != nil {
			return nil, fmt.Errorf("failed to create gorm instance, caused by: %w", err)
		}

		database = db
	default:
		panic(fmt.Sprintf("Engine: %s is not supported", cfg.Engine))
	}

	return &Database{
		DB:     database,
		logger: logger,
	}, nil
}

// mergeConfig merges the default configuration with the provided one
func mergeConfig(defaultCfg *defs.Database, providedCfg *defs.Database) *defs.Database {
	if providedCfg == nil {
		return defaultCfg
	}

	mergedCfg := *defaultCfg

	if providedCfg.Engine != "" {
		mergedCfg.Engine = providedCfg.Engine
	}

	if providedCfg.SQLiteConfig.ConnectionString != "" {
		mergedCfg.SQLiteConfig.ConnectionString = providedCfg.SQLiteConfig.ConnectionString
	}

	if providedCfg.MaxIdleConnections > 0 {
		mergedCfg.MaxIdleConnections = providedCfg.MaxIdleConnections
	}

	if providedCfg.MaxConnectionIdleTime > 0 {
		mergedCfg.MaxConnectionIdleTime = providedCfg.MaxConnectionIdleTime
	}

	if providedCfg.MaxConnectionTime > 0 {
		mergedCfg.MaxConnectionTime = providedCfg.MaxConnectionTime
	}

	if providedCfg.MaxOpenConnections > 0 {
		mergedCfg.MaxOpenConnections = providedCfg.MaxOpenConnections
	}

	return &mergedCfg
}

// openSQLiteDatabase will open a SQLite database connection
func openSQLiteDatabase(cfg *defs.Database, logger glogger.Interface) (*gorm.DB, error) {
	dsn := cfg.SQLiteConfig.ConnectionString
	if dsn == "" {
		dsn = defs.DSNDefault
	}

	dialector := sqlite.New(sqlite.Config{
		Conn:       nil,
		DriverName: sqlite3extended.NAME,
		DSN:        dsn,
	})

	// create new connection
	db, err := gorm.Open(dialector, createGormConfig(
		logger,
	))
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to create new database connection with gorm"))
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, errors.Join(err, errors.New("failed to connect to the database with gorm"))
	}
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConnections)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConnections)
	sqlDB.SetConnMaxLifetime(cfg.MaxConnectionTime)
	sqlDB.SetConnMaxIdleTime(cfg.MaxConnectionIdleTime)

	// Return the connection
	return db, nil
}

// openPostgresDatabase will open a PostrgreSQL database connection
func openPostgresDatabase(_ *defs.Database, _ glogger.Interface) (*gorm.DB, error) {

	return nil, nil
}

// openMySQLDatabase will open a MySQL database connection
func openMySQLDatabase(_ *defs.Database, _ glogger.Interface) (*gorm.DB, error) {

	return nil, nil
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
