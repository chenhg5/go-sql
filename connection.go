package connection

import (
	"database/sql"
)

const (
	// DriverMysql is a const value of mysql driver.
	DriverMysql = "mysql"
	// DriverSqlite is a const value of sqlite driver.
	DriverSqlite = "sqlite"
	// DriverPostgresql is a const value of postgresql driver.
	DriverPostgresql = "postgresql"
	// DriverMssql is a const value of mssql driver.
	DriverMssql = "mssql"
)

// Connection is a connection handler of database.
type Connection interface {
	// InitDB initialize the database connections.
	InitDB(cfg map[string]Database) Connection

	// Name get the connection`s name.
	Name() string

	// GetDelimiter get the default delimiter.
	GetDelimiter() string

	// Query is the query method of sql.
	Query(query string, args ...interface{}) ([]map[string]interface{}, error)

	// Exec is the exec method of sql.
	Exec(query string, args ...interface{}) (sql.Result, error)

	// QueryWithConnection is the query method with given connection of sql.
	QueryWithConnection(conn, query string, args ...interface{}) ([]map[string]interface{}, error)

	// ExecWithConnection is the exec method with given connection of sql.
	ExecWithConnection(conn, query string, args ...interface{}) (sql.Result, error)

	// Transaction API
	// ===================================

	QueryWithTx(tx *sql.Tx, query string, args ...interface{}) ([]map[string]interface{}, error)

	ExecWithTx(tx *sql.Tx, query string, args ...interface{}) (sql.Result, error)

	BeginTxWithReadUncommitted() *sql.Tx
	BeginTxWithReadCommitted() *sql.Tx
	BeginTxWithRepeatableRead() *sql.Tx
	BeginTx() *sql.Tx
	BeginTxWithLevel(level sql.IsolationLevel) *sql.Tx

	BeginTxWithReadUncommittedAndConnection(conn string) *sql.Tx
	BeginTxWithReadCommittedAndConnection(conn string) *sql.Tx
	BeginTxWithRepeatableReadAndConnection(conn string) *sql.Tx
	BeginTxAndConnection(conn string) *sql.Tx
	BeginTxWithLevelAndConnection(conn string, level sql.IsolationLevel) *sql.Tx
}

// GetConnectionByDriver return the Connection by given driver name.
func GetConnectionByDriver(driver string) Connection {
	switch driver {
	case DriverMysql:
		return GetMysqlDB()
	case DriverMssql:
		return GetMssqlDB()
	case DriverSqlite:
		return GetSqliteDB()
	case DriverPostgresql:
		return GetPostgresqlDB()
	default:
		panic("driver not found!")
	}
}

func MysqlConnection() Connection {
	return GetMysqlDB()
}

func MssqlConnection() Connection {
	return GetMssqlDB()
}

func SqliteConnection() Connection {
	return GetSqliteDB()
}

func PostgresqlConnection() Connection {
	return GetPostgresqlDB()
}

type Database struct {
	Dsn        string
	User       string
	Pwd        string
	Host       string
	Port       string
	Name       string
	File       string
	Params     Params
	MaxIdleCon int
	MaxOpenCon int
}

type Params map[string]string

type Databases map[string]Database

func (d Databases) Add(key string, db Database) Databases {
	d[key] = db
	return d
}

func DefaultDatabases(database Database) Databases {
	return Databases{"default": database}
}
