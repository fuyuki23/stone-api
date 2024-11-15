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

type contextKey string

var sessionKey = contextKey("session")

type StoneHandler = func(r *http.Request) (any, error)

func (api *API) BaseHandler(handle StoneHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if res, err := handle(r); err != nil {
			err = response.Fail(err).Send(w)
			if err != nil {
				log.Error().Err(err).Send()
			}
		} else {
			if apiRes, ok := res.(response.ApiResponse); !ok {
				err = response.Ok(res).Send(w)
			} else {
				err = apiRes.Send(w)
			}
			if err != nil {
				log.Error().Err(err).Send()
			}
		}
	})
}

func (api *API) AuthHandler(handle StoneHandler) http.Handler {
	return api.BaseHandler(func(r *http.Request) (any, error) {
		authorization := strings.Trim(r.Header.Get("Authorization"), " ")
		if authorization == "" || !strings.HasPrefix(authorization, "Bearer ") {
			log.Debug().Msg("no or invalid authorization header")
			return nil, model.ErrUnauthorized
		}
		accessToken := authorization[7:]
		isOk, err := token.ValidateToken("access", accessToken)
		if err != nil || !isOk {
			log.Debug().Err(err).Msg("failed to validate access token")
			return nil, model.ErrUnauthorized
		}
		email, err := token.GetEmailFromToken(accessToken)
		if err != nil {
			log.Debug().Err(err).Msg("failed to get email from token")
			return nil, model.ErrUnauthorized
		}

		user, err := api.user.userStore.FindByEmail(email)
		if err != nil {
			log.Debug().Err(err).Msg("failed to query user")
			return nil, model.ErrUnknown
		}
		if user == nil {
			log.Debug().Msg("user not found")
			return nil, model.ErrUnauthorized
		}

		r = r.WithContext(context.WithValue(r.Context(), sessionKey, user))

		return handle(r)
	})
}

func getUser(r *http.Request, user *model.User) error {
	maybeUser, ok := r.Context().Value(sessionKey).(*db.UserEntity)
	if !ok {
		log.Debug().Msg("failed to get user from context")
		return model.ErrUnauthorized
	}

	*user = maybeUser.ConvertToModel()
	return nil
}

var NotFound = func(w http.ResponseWriter, _ *http.Request) {
	log.Error().Err(response.Fail(model.ErrNotFound).Status(http.StatusNotFound).Send(w)).Send()
}
