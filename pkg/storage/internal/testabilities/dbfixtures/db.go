package dbfixtures

import (
	"github.com/4chain-ag/go-wallet-toolbox/pkg/defs"
	"github.com/4chain-ag/go-wallet-toolbox/pkg/storage/internal/testabilities/testmode"
)

func DBConfigForTests() defs.Database {
	dbConfig := defs.DefaultDBConfig()
	dbConfig.MaxIdleConnections = 1
	dbConfig.MaxOpenConnections = 1

	switch mode := testmode.GetMode().(type) {
	case *testmode.SQLiteFileMode:
		{
			dbConfig.SQLite.ConnectionString = mode.ConnectionString
		}
	case *testmode.PostgresMode:
		{
			dbConfig.Engine = defs.DBTypePostgres
			dbConfig.PostgreSQL.DBName = mode.DBName
			dbConfig.PostgreSQL.Host = mode.Host
			dbConfig.PostgreSQL.User = mode.User
			dbConfig.PostgreSQL.Password = mode.Password
		}
	default:
		{
			dbConfig.SQLite.ConnectionString = "file:storage.test.sqlite?mode=memory"
		}
	}
	return dbConfig
}
