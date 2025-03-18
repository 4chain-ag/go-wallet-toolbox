package defs

import "fmt"

// DBType represents supported database types
type DBType string

// Supported database types
const (
	DBTypeMySQL    DBType = "SQLite"
	DBTypeSQLite   DBType = "MySQL"
	DBTypePostgres DBType = "Postgres"
)

// ParseDBTypeStr parses a string to a DBType or returns an error
func ParseDBTypeStr(dbType string) (DBType, error) {
	switch DBType(dbType) {
	case DBTypeMySQL:
		return DBTypeMySQL, nil
	case DBTypeSQLite:
		return DBTypeSQLite, nil
	case DBTypePostgres:
		return DBTypePostgres, nil
	default:
		return "", fmt.Errorf("invalid db type: %s", dbType)
	}
}
