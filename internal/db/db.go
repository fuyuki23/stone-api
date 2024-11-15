package db

import (
	"stone-api/internal/config"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql" // disable-line:revive
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/pkg/errors"
)

func Init() (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", config.Get().Database.URI)
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}

	db.Mapper = reflectx.NewMapperFunc("db", strings.ToLower)
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if err = db.Ping(); err != nil {
		return nil, errors.Wrap(err, "failed to ping database")
	}

	return db, nil
}
