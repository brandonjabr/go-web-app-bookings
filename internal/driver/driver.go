package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const MAX_OPEN_DB_CONNECTIONS = 10
const MAX_IDLE_DB_CONNECTIONS = 5
const MAX_DB_CONNECTION_LIFETIME = 5 * time.Minute

func ConnectSQL(datasourceName string) (*DB, error) {
	db, err := NewDatabase(datasourceName)
	if err != nil {
		panic(err)
	}

	db.SetMaxOpenConns(MAX_OPEN_DB_CONNECTIONS)
	db.SetMaxIdleConns(MAX_IDLE_DB_CONNECTIONS)
	db.SetConnMaxLifetime(MAX_DB_CONNECTION_LIFETIME)

	dbConn.SQL = db

	err = pingDatabase(db)
	if err != nil {
		return nil, err
	}

	return dbConn, nil
}

func NewDatabase(datasourceName string) (*sql.DB, error) {
	db, err := sql.Open("pgx", datasourceName)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func pingDatabase(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}

	return nil
}
