package api

import (
	"context"
	"net/http"
	"stone-api/internal/db"
	"stone-api/internal/model"
	"stone-api/internal/response"
	"stone-api/internal/token"
	"strings"

	"github.com/rs/zerolog/log"
)

type StoneHandler = func(w http.ResponseWriter, r *http.Request) error

func (api *Api) BaseHandler(handle StoneHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := handle(w, r); err != nil {
			_ = response.Fail(err).Send(w)
		}
	})
}

func (api *Api) AuthHandler(handle StoneHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorization := strings.Trim(r.Header.Get("Authorization"), " ")
		if authorization == "" || !strings.HasPrefix(authorization, "Bearer ") {
			log.Debug().Msg("no or invalid authorization header")
			_ = response.Fail(model.ErrUnauthorized).Status(http.StatusUnauthorized).Send(w)
			return
		}
		accessToken := authorization[7:]
		isOk, err := token.ValidateToken("access", accessToken)
		if err != nil || !isOk {
			log.Debug().Err(err).Msg("failed to validate access token")
			_ = response.Fail(model.ErrUnauthorized).Status(http.StatusUnauthorized).Send(w)
			return
		}
		email, err := token.GetEmailFromToken(accessToken)
		if err != nil {
			log.Debug().Err(err).Msg("failed to get email from token")
			_ = response.Fail(model.ErrUnauthorized).Status(http.StatusUnauthorized).Send(w)
			return
		}

		user, err := api.user.userStore.FindByEmail(email)
		if err != nil {
			log.Debug().Err(err).Msg("failed to query user")
			_ = response.Fail(model.ErrUnknown).Status(http.StatusInternalServerError).Send(w)
			return
		}
		if user == nil {
			log.Debug().Msg("user not found")
			_ = response.Fail(model.ErrUnauthorized).Status(http.StatusUnauthorized).Send(w)
			return
		}

		r = r.WithContext(context.WithValue(r.Context(), "session", user))

		if err := handle(w, r); err != nil {
			_ = response.Fail(err).Send(w)
		}
	})
}

func getUser(r *http.Request, user *model.User) error {
	maybeUser, ok := r.Context().Value("session").(*db.UserEntity)
	if !ok {
		log.Debug().Msg("failed to get user from context")
		return model.ErrUnauthorized
	}

	*user = maybeUser.ConvertToModel()
	return nil
}

var NotFound = func(w http.ResponseWriter, r *http.Request) {
	_ = response.Fail(model.ErrNotFound).Status(http.StatusNotFound).Send(w)
}
