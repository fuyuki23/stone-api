package db

import "github.com/jmoiron/sqlx"

type Store struct {
	db   *sqlx.DB
	user *UserStore
}

func (s *Store) UserStore() *UserStore {
	return s.user
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db:   db,
		user: NewUserStore(db),
	}
}

func (s *Store) DB() *sqlx.DB {
	return s.db
}
