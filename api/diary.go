package api

import (
	"net/http"
	"slices"
	"stone-api/internal/db"
	"stone-api/internal/model"
	"stone-api/internal/response"
	"stone-api/internal/utils"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type DiaryHandler struct {
	userStore  *db.UserStore
	diaryStore *db.DiaryStore
}

func (api *API) initDiaryAPI(router *mux.Router) {
	api.diary = &DiaryHandler{
		userStore:  api.serv.Store().UserStore(),
		diaryStore: api.serv.Store().DiaryStore(),
	}

	router.Handle("", api.AuthHandler(api.diary.listWithRange)).Methods(http.MethodGet).Name("List Diary")
	router.Handle("", api.AuthHandler(api.diary.create)).Methods(http.MethodPost).Name("Create Diary")
	router.Handle("/{id}", api.AuthHandler(api.diary.update)).Methods(http.MethodPatch).Name("Update Diary")
	router.Handle("/{id}", api.AuthHandler(api.diary.delete)).Methods(http.MethodDelete).Name("Delete Diary")
}

type ListDiaryWithRangeRequest struct {
	Start time.Time
	End   time.Time
}

func (r *ListDiaryWithRangeRequest) Parse(req *http.Request) error {
	startDateRange := strings.Trim(req.URL.Query().Get("start"), " ")
	endDateRange := strings.Trim(req.URL.Query().Get("end"), " ")
	if startDateRange == "" && endDateRange == "" {
		startDateRange = time.Now().Format("2006-01-02")
		endDateRange = time.Now().Format("2006-01-02")
	} else if startDateRange == "" || endDateRange == "" {
		log.Error().Msg("invalid date range")
		return model.ErrBadRequest
	}

	if !utils.IsDate.Match([]byte(startDateRange)) || !utils.IsDate.Match([]byte(endDateRange)) {
		log.Error().Msg("invalid date range")
		return model.ErrBadRequest
	}

	startDate, err := time.Parse("2006-01-02 15:04:05", utils.AppendString(startDateRange, " 00:00:00"))
	if err != nil {
		log.Error().Err(err).Msg("failed to parse start date")
		return model.ErrBadRequest
	}
	endDate, err := time.Parse("2006-01-02 15:04:05", utils.AppendString(endDateRange, " 23:59:59"))
	if err != nil {
		log.Error().Err(err).Msg("failed to parse end date")
		return model.ErrBadRequest
	}

	r.Start = startDate
	r.End = endDate
	return nil
}

func (h *DiaryHandler) listWithRange(r *http.Request) (any, error) {
	var sessionUser model.User
	if err := getUser(r, &sessionUser); err != nil {
		return nil, err
	}

	payload := ListDiaryWithRangeRequest{}
	if err := payload.Parse(r); err != nil {
		return nil, err
	}

	diaryEntities, err := h.diaryStore.FindWithRange(db.BUID(sessionUser.ID), payload.Start, payload.End)
	if err != nil {
		return nil, err
	}

	var diaries = make([]model.Diary, len(diaryEntities))
	for i, diary := range diaryEntities {
		diaries[i] = diary.ConvertToModel()
	}

	return diaries, nil
}

type CreateDiaryRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Mood    string `json:"mood"`
}

func (r CreateDiaryRequest) Validate() error {
	if r.Title == "" {
		log.Error().Msg("title is required")
		return model.ErrBadRequest
	}
	titleLength := utf8.RuneCountInString(r.Title)
	if titleLength > 255 {
		log.Error().Msg("title is too long")
		return model.ErrBadRequest
	}
	if r.Content == "" {
		log.Error().Msg("content is required")
		return model.ErrBadRequest
	}
	contentLength := utf8.RuneCountInString(r.Content)
	if contentLength > 1024 {
		log.Error().Msg("content is too long")
		return model.ErrBadRequest
	}
	if r.Mood == "" {
		log.Error().Msg("mood is required")
		return model.ErrBadRequest
	}
	if idx := slices.Index(db.DiaryMoodAll, r.Mood); idx == -1 {
		log.Error().Msg("invalid mood")
		return model.ErrBadRequest
	}

	return nil
}

func (h *DiaryHandler) create(r *http.Request) (any, error) {
	var sessionUser model.User
	if err := getUser(r, &sessionUser); err != nil {
		return nil, err
	}

	var payload CreateDiaryRequest
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return nil, model.ErrBadRequest
	}
	if err := payload.Validate(); err != nil {
		return nil, err
	}

	isExist, err := h.diaryStore.FindByDate(db.BUID(sessionUser.ID), time.Now())
	if err != nil {
		return nil, err
	} else if isExist != nil {
		log.Error().Msg("diary already exists")
		return nil, model.ErrDiaryAlreadyExists
	}

	newDiaryID, err := uuid.NewV7()
	if err != nil {
		return nil, err
	}
	newDiary := db.DiaryEntity{
		ID:      db.BUID(newDiaryID),
		UserID:  db.BUID(sessionUser.ID),
		Title:   payload.Title,
		Content: payload.Content,
		Mood:    payload.Mood,
	}
	if err = h.diaryStore.Create(&newDiary); err != nil {
		return nil, err
	}

	return response.Ok(newDiary.ConvertToModel()).Status(http.StatusCreated), nil
}

type UpdateDiaryRequest struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
	Mood    *string `json:"mood"`
}

func (r UpdateDiaryRequest) Validate() error {
	if r.Title == nil && r.Content == nil && r.Mood == nil {
		log.Error().Msg("no field to update")
		return model.ErrBadRequest
	}

	if r.Title != nil {
		if *r.Title == "" {
			log.Error().Msg("title is required")
			return model.ErrBadRequest
		}
		titleLength := utf8.RuneCountInString(*r.Title)
		if titleLength > 255 {
			log.Error().Msg("title is too long")
			return model.ErrBadRequest
		}
	}
	if r.Content != nil {
		if *r.Content == "" {
			log.Error().Msg("content is required")
			return model.ErrBadRequest
		}
		contentLength := utf8.RuneCountInString(*r.Content)
		if contentLength > 1024 {
			log.Error().Msg("content is too long")
			return model.ErrBadRequest
		}
	}
	if r.Mood != nil {
		if *r.Mood == "" {
			log.Error().Msg("mood is required")
			return model.ErrBadRequest
		}
		if idx := slices.Index(db.DiaryMoodAll, *r.Mood); idx == -1 {
			log.Error().Msg("invalid mood")
			return model.ErrBadRequest
		}
	}

	return nil
}

func (h *DiaryHandler) update(r *http.Request) (any, error) {
	diaryPathID, ok := mux.Vars(r)["id"]
	if !ok {
		return nil, model.ErrNotFound
	}
	diaryID, err := uuid.Parse(diaryPathID)
	if err != nil {
		return nil, model.ErrNotFound
	}

	var sessionUser model.User
	if err = getUser(r, &sessionUser); err != nil {
		return nil, err
	}

	var payload UpdateDiaryRequest
	if err = json.NewDecoder(r.Body).Decode(&payload); err != nil {
		return nil, model.ErrBadRequest
	}
	if err = payload.Validate(); err != nil {
		return nil, err
	}

	diary, err := h.diaryStore.FindByID(db.BUID(diaryID), db.BUID(sessionUser.ID))
	if err != nil {
		return nil, err
	}
	if diary == nil {
		return nil, model.ErrDiaryNotFound
	}

	if !utils.IsSameDate(diary.CreatedAt, time.Now()) {
		return nil, model.ErrDiaryNotToday
	}

	if payload.Title != nil {
		diary.Title = *payload.Title
	}
	if payload.Content != nil {
		diary.Content = *payload.Content
	}
	if payload.Mood != nil {
		diary.Mood = *payload.Mood
	}

	if err = h.diaryStore.Update(diary); err != nil {
		return nil, err
	}

	return diary.ConvertToModel(), nil
}

func (h *DiaryHandler) delete(r *http.Request) (any, error) {
	diaryPathID, ok := mux.Vars(r)["id"]
	if !ok {
		return nil, model.ErrNotFound
	}
	diaryID, err := uuid.Parse(diaryPathID)
	if err != nil {
		return nil, model.ErrNotFound
	}

	var sessionUser model.User
	if err = getUser(r, &sessionUser); err != nil {
		return nil, err
	}

	diary, err := h.diaryStore.FindByID(db.BUID(diaryID), db.BUID(sessionUser.ID))
	if err != nil {
		return nil, err
	}
	if diary == nil {
		return nil, model.ErrDiaryNotFound
	}

	err = h.diaryStore.DeleteByID(db.BUID(diaryID), db.BUID(sessionUser.ID))
	if err != nil {
		return nil, err
	}

	return response.Ok("").Status(http.StatusNoContent), nil
}
