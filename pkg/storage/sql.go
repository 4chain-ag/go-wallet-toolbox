package storage

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/sqlite3extended"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const (
	dsnDefault                = "file::memory:" // DSN for connection (file or memory, default is memory)
	defaultPreparedStatements = false           // Flag for prepared statements for SQL
)

// openSQLiteDatabase will open a SQLite database connection
func openSQLiteDatabase(cfg *Config, logger glogger.Interface) (*gorm.DB, error) {
	dialector := sqlite.New(sqlite.Config{
		Conn:       cfg.ExistingConnection,
		DriverName: sqlite3extended.NAME,
		DSN:        getDSN(cfg.SQLiteConfig.DatabasePath, cfg.SQLiteConfig.Shared),
	})

	// create new connection
	db, err := gorm.Open(dialector, getGormConfig(
		cfg.TablePrefix,
		defaultPreparedStatements,
		cfg.Debug,
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

// getDSN will return the DNS string
func getDSN(databasePath string, shared bool) string {
	var dsn string
	// Use a file based path?
	if len(databasePath) > 0 {
		dsn = databasePath
	} else { // Default is in-memory
		dsn = dsnDefault
	}

	if shared {
		dsn += "?cache=shared"
	}

	return dsn
}

// getGormConfig returns valid gorm.Config for database connections
func getGormConfig(tablePrefix string, preparedStatement, debug bool, logger glogger.Interface) *gorm.Config {
	// Set the prefix
	if len(tablePrefix) > 0 {
		tablePrefix = tablePrefix + "_"
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
		PrepareStmt:            preparedStatement, // default is: false
		QueryFields:            false,
		SkipDefaultTransaction: false,
		TranslateError:         true,
	}

	// Optional logger vs basic
	if logger == nil {
		logLevel := glogger.Silent
		if debug {
			logLevel = glogger.Info
		}

		config.Logger = glogger.New(
			log.New(os.Stdout, "\r\n ", log.LstdFlags), // io writer
			glogger.Config{
				SlowThreshold:             5 * time.Second, // Slow SQL threshold
				LogLevel:                  logLevel,        // Log level
				IgnoreRecordNotFoundError: true,            // Ignore ErrRecordNotFound error for logger
				Colorful:                  false,           // Disable color
			},
		)
	}

	return config
}
