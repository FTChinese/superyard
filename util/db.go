package util

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"
)

// NewDB creates a db connection
func NewDB(host, port, user, pass string) (*sql.DB, error) {
	cfg := &mysql.Config{
		User:                 user,
		Passwd:               pass,
		Net:                  "tcp",
		Addr:                 host + ":" + port,
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
