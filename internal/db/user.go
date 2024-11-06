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
	ID        uuid.UUID `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"create_at"`
	UpdatedAt time.Time `db:"update_at"`
}

func (u UserEntity) ConvertToUser() model.User {
	return model.User{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
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

func NewUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{db: db}
}
