package database

import (
	"errors"
	"fmt"
	"log/slog"

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
func NewDatabase(cfg *Config, logger *slog.Logger) (*Database, error) {
	var database *gorm.DB

	if logger == nil {
		logger = slog.Default()
	}

	switch cfg.Engine {
	case SQLite:
		db, err := openSQLiteDatabase(cfg, &SlogGormLogger{
			logger: logger,
		})
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

// openSQLiteDatabase will open a SQLite database connection
func openSQLiteDatabase(cfg *Config, logger glogger.Interface) (*gorm.DB, error) {
	dsn := cfg.SQLiteConfig.ConnectionString
	if dsn == "" {
		dsn = dsnDefault
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

// createGormConfig returns valid gorm.Config for database connections
func createGormConfig(logger glogger.Interface) *gorm.Config {
	// Set the prefix
	tablePrefix := defaultTablePrefix

	if logger == nil {
		panic("Could not create gorm config. When creating database configuration you need to specify the logger to use")
	}

	// Create the configuration
	config := &gorm.Config{
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

	return config
}
