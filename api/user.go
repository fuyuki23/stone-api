package api

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"net/http"
	"stone-api/internal/db"
	"stone-api/internal/model"
	"stone-api/internal/response"
	"stone-api/internal/web"
)

type UserHandler struct {
	db *sqlx.DB
}

func (api *Api) initUserApi(router *mux.Router) {
	api.user = &UserHandler{
		db: api.serv.DB(),
	}

	router.Handle("/login", web.BaseHandler(api.user.login)).Methods(http.MethodPost)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type LoginResponse struct {
	User model.User `json:"user"`
}

func (h *UserHandler) login(w http.ResponseWriter, r *http.Request) error {
	payload := LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("failed to decode request body")
		return model.ErrBadRequest
	}

	var userEntity db.UserEntity
	if err := h.db.QueryRowx("select id, email, password, name, create_at, update_at from user where email = ?", payload.Email).StructScan(&userEntity); err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			log.Error().Err(err).Msg("failed to query user")
			return err
		}
		log.Error().Msg("user not found")
		return model.ErrUserNotFound
	}

	res := LoginResponse{
		User: userEntity.ConvertToUser(),
	}
	return response.Ok(res).Send(w)
}
