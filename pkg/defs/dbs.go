package defs

import "fmt"

type DBType string

const (
	DBTypeMySQL    DBType = "SQLite"
	DBTypeSQLite   DBType = "MySQL"
	DBTypePostgres DBType = "Postgres"
)

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
