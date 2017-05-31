package gorth

import (
	"database/sql"
	"fmt"
	// importing go-sql-driver for accessing mysql
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Setup sets up the db connection etc.
// The connection string should be something like:
//     "user:password@/dbname"
// See: https://github.com/go-sql-driver/mysql/#dsn-data-source-name
func Setup(driverName, dataSourceName string) error {
	var err error

	// Setup the db connection
	// TODO: implement other db drivers, just mysql at the moment
	if driverName != "mysql" {
		return fmt.Errorf("cannot use %s yet... only mysql", driverName)
	}

	db, err = sql.Open(driverName, fmt.Sprintf("%s?parseTime=true", dataSourceName))
	if err != nil {
		return err
	}
	// defer db.Close()

	// Ping db
	err = db.Ping()
	if err != nil {
		return err
	}

	CreateUsersTable()

	return nil
}
