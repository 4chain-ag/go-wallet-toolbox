package storage

import (
	"fmt"
	"log/slog"
	"os"
	"time"

	"gorm.io/gorm"
)

// Engine is the different engines that are supported (database)
type Engine string

// Supported engines (databases)
const (
	SQLite Engine = "sqlite"
)

// Storage is a struct that holds logger for database connecton and the connection itself
type Storage struct {
	DB     *gorm.DB
	logger *slog.Logger
}

// StorageConfig is a struct that configures the database connection
type StorageConfig struct {
	// LogLevel is the importance and amount of information printed: debug, info, warn, error
	LogLevel string
	// Datastore engine (PostgreSQL, SQLite)
	Engine Engine
	// TablePrefix is a prefix that will be added to each tablePrefix
	TablePrefix string
	// Debug is a flag that will allow debugging in the database
	Debug bool

	// SQLiteConfig is configuration struct for SQLite database
	SQLiteConfig *SQLiteConfig

	// MaxIdleConnections defines the maximum number of idle connections allowed for the database.
	MaxIdleConnections int

	// MaxConnectionIdleTime sets the maximum duration an idle connection can remain open before being closed.
	// Typically set in seconds.
	MaxConnectionIdleTime time.Duration

	// MaxConnectionTime defines the maximum amount of time a connection may be reused.
	// Typically set in seconds.
	MaxConnectionTime time.Duration

	// MaxOpenConnections specifies the maximum number of open connections to the database.
	MaxOpenConnections int

	// ExistingConnection is a param used for connecting to existing db connection
	ExistingConnection gorm.ConnPool
}

type SQLiteConfig struct {
	// DatabasePath is the path to sqlite DB
	DatabasePath string
	// Shared determines whether the database connection is shared among multiple instances.
	Shared bool
}

func NewStorage(cfg *StorageConfig, logger *slog.Logger) *Storage {
	var store *gorm.DB
	logLevel := parseLogLevel(cfg.LogLevel)

	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
			Level: logLevel,
		}))
	}

	gormLogger := &SlogGormLogger{
		logger: logger,
		level:  logLevel,
	}

	switch cfg.Engine {
	case SQLite:
		db, err := openSQLiteDatabase(cfg, gormLogger)
		if err != nil {
			panic(fmt.Errorf("failed to create gorm instance, caused by: %w", err))
		}

		store = db
	default:
		panic(fmt.Errorf("Engine: %s is not supported", cfg.Engine))
	}

	return &Storage{
		DB:     store,
		logger: logger,
	}
}
