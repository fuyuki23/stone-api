package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"stone-api/internal/db"
	"stone-api/internal/model"
	"stone-api/internal/response"
	"stone-api/internal/token"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userStore *db.UserStore
}

func (api *Api) initUserApi(router *mux.Router) {
	api.user = &UserHandler{
		userStore: api.serv.Store().UserStore(),
	}

	router.Handle("/login", api.BaseHandler(api.user.login)).Methods(http.MethodPost).Name("Login")
	router.Handle("/register", api.BaseHandler(api.user.register)).Methods(http.MethodPost).Name("Register")
	router.Handle("/me", api.AuthHandler(api.user.me)).Methods(http.MethodGet).Name("Me")
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	User   model.User   `json:"user"`
	Tokens model.Tokens `json:"tokens"`
}

func (h *UserHandler) login(w http.ResponseWriter, r *http.Request) error {
	payload := LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("failed to decode request body")
		return model.ErrBadRequest
	}

	userEntity, err := h.userStore.FindByEmail(payload.Email)
	if err != nil {
		log.Error().Err(err).Msg("failed to query user")
		return err
	}
	if userEntity == nil {
		return model.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(userEntity.Password), []byte(payload.Password)); err != nil {
		return model.ErrInvalidCredentials
	}

	user := userEntity.ConvertToUser()
	tokens, err := token.CreateTokens(user)
	if err != nil {
		log.Error().Err(err).Msg("failed to create tokens")
		return model.ErrUnknown
	}

	res := LoginResponse{
		User:   userEntity.ConvertToUser(),
		Tokens: *tokens,
	}
	return response.Ok(res).Send(w)
}

type RegisterRequest struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Name     *string `json:"name"`
}

//type RegisterResponse = string

func (h *UserHandler) register(w http.ResponseWriter, r *http.Request) error {
	payload := RegisterRequest{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("failed to decode request body")
		return model.ErrBadRequest
	}

	existUser, err := h.userStore.FindByEmail(payload.Email)
	if err != nil {
		log.Error().Err(err).Msg("failed to query user")
		return err
	}
	if existUser != nil {
		return model.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash password")
		return model.ErrUnknown
	}

	var maybeName sql.NullString
	if payload.Name != nil {
		maybeName = sql.NullString{
			String: *payload.Name,
			Valid:  true,
		}
	} else {
		maybeName = sql.NullString{
			Valid: false,
		}
	}

	newUserID, err := uuid.NewV7()
	if err != nil {
		log.Error().Err(err).Msg("failed to generate user id")
		return model.ErrUnknown
	}

	newUser := db.UserEntity{
		ID:        model.BUID(newUserID),
		Email:     payload.Email,
		Password:  string(hashedPassword),
		Name:      maybeName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	if err := h.userStore.Create(newUser); err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return err
	}

	return response.Ok("ok").Send(w)
}

func (h *UserHandler) me(w http.ResponseWriter, r *http.Request) error {
	var sessionUser model.User
	if err := getUser(r, &sessionUser); err != nil {
		return err
	}

	log.Debug().Interface("sessionUser", sessionUser).Msg("authorized")

	return response.Ok(sessionUser).Send(w)
}
