package database

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"time"
)

const (
	dsnDefault         = "file::memory:" // DSN for connection (file or memory, default is memory)
	defaultTablePrefix = "bsv_"
)

// Config is a struct that configures the database connection
type Config struct {
	// Engine is the database engine (PostgreSQL, SQLite)
	Engine defs.DBType

	// SQLiteConfig is configuration struct for SQLite database
	SQLiteConfig SQLiteConfig

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

// SQLiteConfig is configuration struct for SQLite database
type SQLiteConfig struct {
	// ConnectionString is the path to SQLite DB
	ConnectionString string
}
