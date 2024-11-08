package db

import (
	"database/sql"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"stone-api/internal/model"
	"stone-api/internal/utils"
	"time"
)

type DiaryMood = string

const (
	DiaryMoodSurprise DiaryMood = "surprise"
	DiaryMoodAngry    DiaryMood = "angry"
	DiaryMoodSad      DiaryMood = "sad"
	DiaryMoodNeutral  DiaryMood = "neutral"
	DiaryMoodMad      DiaryMood = "Mad"
	DiaryMoodCry      DiaryMood = "cry"
	DiaryMoodHappy    DiaryMood = "happy"
	DiaryMoodExhaust  DiaryMood = "exhaust"
)

var DiaryMoodAll = []DiaryMood{
	DiaryMoodSurprise,
	DiaryMoodAngry,
	DiaryMoodSad,
	DiaryMoodNeutral,
	DiaryMoodMad,
	DiaryMoodCry,
	DiaryMoodHappy,
	DiaryMoodExhaust,
}

type DiaryEntity struct {
	ID        BUID      `db:"id"`
	UserID    BUID      `db:"user_id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	Mood      DiaryMood `db:"mood"`
	CreatedAt time.Time `db:"create_at"`
	UpdatedAt time.Time `db:"update_at"`
}

func (d DiaryEntity) ConvertToModel() model.Diary {
	return model.Diary{
		ID:        uuid.UUID(d.ID),
		Title:     d.Title,
		Content:   d.Content,
		Mood:      d.Mood,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

type DiaryStore struct {
	db *sqlx.DB
}

func NewDiaryStore(db *sqlx.DB) *DiaryStore {
	return &DiaryStore{db: db}
}

func (s *DiaryStore) FindWithRange(userID BUID, start time.Time, end time.Time) ([]DiaryEntity, error) {
	var diaries = make([]DiaryEntity, 0)
	rows, err := s.db.Queryx(`
		SELECT * FROM diary
		WHERE user_id = ?
		AND create_at >= ? AND create_at <= ?
		ORDER BY create_at
	`, userID, start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"))
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var diary DiaryEntity
		if err = rows.StructScan(&diary); err != nil {
			return nil, err
		}
		diaries = append(diaries, diary)
	}

	return diaries, nil
}

func (s *DiaryStore) FindByDate(userID BUID, date time.Time) (*DiaryEntity, error) {
	var diary DiaryEntity
	localDate := date.Format("2006-01-02")
	err := s.db.QueryRowx(`
		SELECT * FROM diary
		WHERE user_id = ?
		AND create_at >= ? AND create_at <= ?
		
	`, userID, utils.AppendString(localDate, " 00:00:00"), utils.AppendString(localDate, " 23:59:59")).StructScan(&diary)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &diary, nil
}

func (s *DiaryStore) FindByID(diaryID BUID, userID BUID) (*DiaryEntity, error) {
	var diary DiaryEntity
	err := s.db.QueryRowx(`
		SELECT * FROM diary
		WHERE user_id = ? AND id = ?
	`, userID, diaryID).StructScan(&diary)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &diary, nil
}

func (s *DiaryStore) Create(diary *DiaryEntity) error {
	_, err := s.db.Exec(`
		INSERT INTO diary (id, user_id, title, content, mood)
		VALUES (?, ?, ?, ?, ?)
	`, diary.ID, diary.UserID, diary.Title, diary.Content, diary.Mood)
	if err != nil {
		return err
	}

	if err = s.db.QueryRowx("select * from diary where id = ?", diary.ID).StructScan(diary); err != nil {
		return err
	}

	return nil
}

func (s *DiaryStore) Update(diary *DiaryEntity) error {
	_, err := s.db.Exec(`
		UPDATE diary
		SET title = ?, content = ?, mood = ?, update_at = NOW()
		WHERE id = ? AND user_id = ?
	`, diary.Title, diary.Content, diary.Mood, diary.ID, diary.UserID)
	if err != nil {
		return err
	}

	if err = s.db.QueryRowx("select * from diary where id = ?", diary.ID).StructScan(diary); err != nil {
		return err
	}

	return nil
}

func (s *DiaryStore) DeleteByID(diaryID BUID, userID BUID) error {
	_, err := s.db.Exec(`
		DELETE FROM diary
		WHERE id = ? AND user_id = ?
	`, diaryID, userID)
	if err != nil {
		return err
	}

	return nil
}
