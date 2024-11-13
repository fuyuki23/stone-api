package db

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"stone-api/internal/model"
	"time"
)

type UserEntity struct {
	ID        BUID           `db:"id"`
	Email     string         `db:"email"`
	Password  string         `db:"password"`
	Name      sql.NullString `db:"name"`
	CreatedAt time.Time      `db:"create_at"`
	UpdatedAt time.Time      `db:"update_at"`
}

func (u UserEntity) ConvertToModel() model.User {
	var name *string
	if u.Name.Valid {
		name = &u.Name.String
	}
	return model.User{
		ID:        uuid.UUID(u.ID),
		Email:     u.Email,
		Name:      name,
		CreatedAt: u.CreatedAt,
	}
}

type UserStore struct {
	db *sqlx.DB
}

func (s *UserStore) FindByEmail(email string) (*UserEntity, error) {
	var user UserEntity
	if err := s.db.QueryRowx("select * from user where email = ?", email).StructScan(&user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, err
	}

	return &user, nil
}

func (s *UserStore) Create(user *UserEntity) error {
	// create user with sqlx.DB and after insert get the last inserted id
	_, err := s.db.Exec("insert into user (id, email, password, name) values (?, ?, ?, ?)", user.ID, user.Email, user.Password, user.Name)
	if err != nil {
		return err
	}

	if err = s.db.QueryRowx("select * from user where id = user.ID").StructScan(user); err != nil {
		return err
	}

	return nil
}

func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{db: db}
}
