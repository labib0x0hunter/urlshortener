package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"urlshortener/utils"

	_ "github.com/go-sql-driver/mysql"
)

// NewMysqlDb creates and returns a new MySQL database connection using environment variables for configuration.
// It also clears the 'urls' table and resets its AUTO_INCREMENT value for a clean state.
// Returns the database connection or an error if any step fails.
func NewMysqlDb() (*sql.DB, error) {
	databaseName := os.Getenv("MYSQL_DB")
	host := os.Getenv("MYSQL_HOST")
	port := os.Getenv("MYSQL_PORT")
	user := os.Getenv("MYSQL_USER")
	password := os.Getenv("MYSQL_PASSWORD")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", user, password, host, port, databaseName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		slog.Error(" [db.go] [OPEN DB] ", slog.Any("error", err))
		return nil, utils.ErrDatabaseConnection
	}

	// Test the database connection
	if err := db.Ping(); err != nil {
		slog.Error(" [db.go] [PING DB] ", slog.Any("error", err))
		return nil, utils.ErrDatabaseConnection
	}

	// // Delete all records from the 'urls' table for a clean start
	// if _, err := db.Exec("DELETE FROM urls"); err != nil {
	// 	slog.Error(" [db.go] [DELETE TABLE] ", slog.Any("error", err))
	// 	return nil, utils.ErrDatabaseDelete
	// }

	// // Reset the AUTO_INCREMENT value for the 'urls' table
	// if _, err := db.Exec("ALTER TABLE urls AUTO_INCREMENT = 1"); err != nil {
	// 	slog.Error(" [db.go] [RESET AUTO_INCREMENT] ", slog.Any("error", err))
	// 	return nil, utils.ErrDatabaseUpdate
	// }

	return db, nil
}
