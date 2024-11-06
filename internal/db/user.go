package db

import (
	"github.com/google/uuid"
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
