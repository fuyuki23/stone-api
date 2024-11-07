package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      *string   `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

// var ErrUserNotFound = NewError(404, "api.user.not_found", "user not found")

var ErrInvalidCredentials = NewError(400, "api.user.invalid_credential", "invalid credentials")
var ErrUserAlreadyExists = NewError(409, "api.user.already_exists", "user already exists")
