package model

import (
	"github.com/google/uuid"
	"time"
)

type Diary struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Mood      string    `json:"mood"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

var ErrDiaryAlreadyExists = NewError(409, "api.diary.already_exists", "diary already exists")
var ErrDiaryNotFound = NewError(404, "api.diary.not_found", "diary not found")
var ErrDiaryNotToday = NewError(400, "api.diary.not_today", "diary is not today")
