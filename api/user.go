package api

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"net/http"
	"stone-api/internal/db"
	"stone-api/internal/model"
	"stone-api/internal/response"
	"stone-api/internal/web"
)

type UserHandler struct {
	userStore *db.UserStore
}

func (api *Api) initUserApi(router *mux.Router) {
	api.user = &UserHandler{
		userStore: api.serv.Store().UserStore(),
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

	userEntity, err := h.userStore.FindByEmail(payload.Email)
	if err != nil {
		log.Error().Err(err).Msg("failed to query user")
		return err
	}
	if userEntity == nil {
		return model.ErrUserNotFound
	}

	res := LoginResponse{
		User: userEntity.ConvertToUser(),
	}
	return response.Ok(res).Send(w)
}
