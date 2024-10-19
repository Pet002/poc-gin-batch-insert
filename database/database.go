package database

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

func New() (*sql.DB, error) {
	client, err := sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/local_database")
	if err != nil {
		return nil, err
	}
	// try to ping to database for check connection
	if err := client.Ping(); err != nil {
		client.Close()
		return nil, err
	}
	return client, nil
}
