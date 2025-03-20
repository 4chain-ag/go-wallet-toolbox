package defs

import (
	"fmt"
	"time"
)

// DBType represents supported database types
type DBType string

// Supported database types
const (
	DBTypeMySQL    DBType = "mysql"
	DBTypeSQLite   DBType = "sqlite"
	DBTypePostgres DBType = "postgres"

	DSNDefault         = "file::memory:" // DSN for connection (file or memory, default is memory)
	DefaultTablePrefix = "bsv_"
)

// Database is a struct that configures the database connection
type Database struct {
	// Engine is the database engine (PostgreSQL, SQLite)
	Engine DBType

	// SQLiteConfig is configuration struct for SQLite database
	SQLiteConfig SQLiteDatabase

	// SQLConfig is configuration for SQL databases such as postgres or mysql
	SQLConfig SQLConfig

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

type SQLConfig struct {
	Host      string
	DBName    string
	Password  string
	Port      string
	Replica   bool
	TimeZone  string
	TxTimeout time.Duration
	User      string
	SslMode   string
}

// ParseDBTypeStr parses a string to a DBType or returns an error
func ParseDBTypeStr(dbType string) (DBType, error) {
	return parseEnumCaseInsensitive(dbType, DBTypeMySQL, DBTypeSQLite, DBTypePostgres)
}

func DefaultDBConfig() *Database {
	return &Database{
		Engine:                DBTypeSQLite,
		SQLiteConfig:          SQLiteDatabase{ConnectionString: DSNDefault},
		MaxIdleConnections:    5,
		MaxConnectionIdleTime: 360 * time.Second,
		MaxConnectionTime:     60 * time.Second,
		MaxOpenConnections:    5,
	}
}

func (db *Database) Validate() (err error) {
	if db.Engine, err = ParseDBTypeStr(string(db.Engine)); err != nil {
		return fmt.Errorf("invalid DB engine: %w", err)
	}

	return nil
}
