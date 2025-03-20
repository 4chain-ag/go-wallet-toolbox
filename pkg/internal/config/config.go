package config

import (
	"time"

	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
)

const (
	DSNDefault         = "file::memory:" // DSN for connection (file or memory, default is memory)
	DefaultTablePrefix = "bsv_"
)

// Database is a struct that configures the database connection
type Database struct {
	// Engine is the database engine (PostgreSQL, SQLite)
	Engine defs.DBType

	// SQLiteConfig is configuration struct for SQLite database
	SQLiteConfig SQLiteDatabase

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
}

// SQLiteDatabase is configuration struct for SQLite database
type SQLiteDatabase struct {
	// ConnectionString is the path to SQLite DB
	ConnectionString string
}

func DefaultDBConfig() *Database {
	return &Database{
		Engine:                defs.DBTypeSQLite,
		SQLiteConfig:          SQLiteDatabase{ConnectionString: DSNDefault},
		MaxIdleConnections:    5,
		MaxConnectionIdleTime: 360 * time.Second,
		MaxConnectionTime:     60 * time.Second,
		MaxOpenConnections:    5,
	}
}
