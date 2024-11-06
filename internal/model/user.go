package model

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
}

var ErrUserNotFound = NewError(404, "api.user.not_found", "user not found")
var ErrInvalidCredentials = NewError(400, "api.user.invalid_credential", "invalid credentials")
