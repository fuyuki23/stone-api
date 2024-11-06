package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"strings"
	"time"
)

func Init() (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", "stone:stone1234@tcp(127.0.0.1:12345)/stone?charset=utf8&parseTime=True&loc=UTC")
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
