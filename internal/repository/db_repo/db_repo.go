package db_repo

import (
	"database/sql"

	"github.com/brandonjabr/go-web-app-bookings/internal/config"
	"github.com/brandonjabr/go-web-app-bookings/internal/repository"
)

type postgresDBRepo struct {
	AppConfig *config.AppConfig
	DB        *sql.DB
}

func NewPostgresRepo(conn *sql.DB, appConfig *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		AppConfig: appConfig,
		DB:        conn,
	}
}
