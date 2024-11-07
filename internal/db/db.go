package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"stone-api/internal/config"
	"strings"
	"time"
)

func Init() (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", config.Get().Database.URI)
	if err != nil {
		return nil, err
	}

	db.Mapper = reflectx.NewMapperFunc("db", strings.ToLower)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
