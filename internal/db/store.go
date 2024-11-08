package db

import "github.com/jmoiron/sqlx"

type Store struct {
	db    *sqlx.DB
	user  *UserStore
	diary *DiaryStore
}

func (s *Store) UserStore() *UserStore {
	return s.user
}

func (s *Store) DiaryStore() *DiaryStore {
	return s.diary
}

func NewStore(db *sqlx.DB) *Store {
	return &Store{
		db:    db,
		user:  NewUserStore(db),
		diary: NewDiaryStore(db),
	}
}

func (s *Store) DB() *sqlx.DB {
	return s.db
}
