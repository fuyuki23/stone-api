package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"net/mail"
	"stone-api/internal/db"
	"stone-api/internal/model"
	"stone-api/internal/response"
	"stone-api/internal/token"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	userStore *db.UserStore
}

func (api *API) initUserAPI(router *mux.Router) {
	api.user = &UserHandler{
		userStore: api.serv.Store().UserStore(),
	}

	router.Handle("/login", api.BaseHandler(api.user.login)).Methods(http.MethodPost).Name("Login")
	router.Handle("/register", api.BaseHandler(api.user.register)).Methods(http.MethodPost).Name("Register")
  router.Handle("/refresh", api.BaseHandler(api.user.refresh)).Methods(http.MethodPost).Name("Refresh Tokens")
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

func (h *UserHandler) login(r *http.Request) (any, error) {
	payload := LoginRequest{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("failed to decode request body")
		return nil, model.ErrBadRequest
	}

	userEntity, err := h.userStore.FindByEmail(payload.Email)
	if err != nil {
		log.Error().Err(err).Msg("failed to query user")
		return nil, err
	}
	if userEntity == nil {
		return nil, model.ErrInvalidCredentials
	}

	if err = bcrypt.CompareHashAndPassword([]byte(userEntity.Password), []byte(payload.Password)); err != nil {
		return nil, model.ErrInvalidCredentials
	}

	user := userEntity.ConvertToModel()
	tokens, err := token.CreateTokens(user)
	if err != nil {
		log.Error().Err(err).Msg("failed to create tokens")
		return nil, model.ErrUnknown
	}

	return LoginResponse{
		User:   user,
		Tokens: *tokens,
	}, nil
}

type RegisterRequest struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Name     *string `json:"name"`
}

func (r RegisterRequest) Validate() error {
	if r.Email == "" {
		log.Error().Msg("email is required")
		return model.ErrBadRequest
	}
	emailLength := utf8.RuneCountInString(r.Email)
	if emailLength > 100 {
		log.Error().Msg("email is too long")
		return model.ErrBadRequest
	}
	if email, err := mail.ParseAddress(r.Email); err != nil || email.Address != r.Email {
		log.Error().Msg("email is invalid")
		return model.ErrBadRequest
	}
	if r.Password == "" {
		log.Error().Msg("password is required")
		return model.ErrBadRequest
	}
	if passwordLength := utf8.RuneCountInString(r.Password); passwordLength < 6 || passwordLength > 32 {
		log.Error().Msg("password length must be between 6 and 32")
		return model.ErrBadRequest
	}
	if r.Name != nil {
		nameLength := utf8.RuneCountInString(*r.Name)
		if nameLength > 50 {
			log.Error().Msg("name is too long")
			return model.ErrBadRequest
		}
	}

	return nil
}

func (r *RegisterRequest) Sanitize() error {
	r.Email = strings.Trim(r.Email, " ")
	if r.Name != nil {
		*r.Name = strings.Trim(*r.Name, " ")
		if len(*r.Name) == 0 {
			r.Name = nil
		}
	}

	return nil
}

func (h *UserHandler) register(r *http.Request) (any, error) {
	payload := RegisterRequest{}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		log.Error().Err(err).Msg("failed to decode request body")
		return nil, model.ErrBadRequest
	}
	if err := payload.Validate(); err != nil {
		return nil, err
	}
	log.Debug().Interface("payload", payload).Msg("before")
	if err := payload.Sanitize(); err != nil {
		return nil, err
	}
	log.Debug().Interface("payload", payload).Msg("after")

	existUser, err := h.userStore.FindByEmail(payload.Email)
	if err != nil {
		log.Error().Err(err).Msg("failed to query user")
		return nil, err
	}
	if existUser != nil {
		return nil, model.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error().Err(err).Msg("failed to hash password")
		return nil, model.ErrUnknown
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
		return nil, model.ErrUnknown
	}

	newUser := db.UserEntity{
		ID:       db.BUID(newUserID),
		Email:    payload.Email,
		Password: string(hashedPassword),
		Name:     maybeName,
	}
	if err = h.userStore.Create(&newUser); err != nil {
		log.Error().Err(err).Msg("failed to create user")
		return nil, err
	}

	return response.Ok("ok").Status(http.StatusCreated), nil
}

type RefreshRequest struct {
  AccessToken string
  RefreshToken string
}

type RefreshResponse struct {
  User  model.User   `json:"user"`
  Tokens model.Tokens `json:"tokens"`
}

func (h *UserHandler) refresh(r *http.Request) (any, error) {
  var payload RefreshRequest
  if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
    return nil, model.ErrBadRequest
  }

  isValid, err := token.ValidateToken("access", payload.AccessToken, true)
  if err != nil {
    return nil, model.ErrUnauthorized
  }
  if !isValid {
    return nil, model.ErrUnauthorized
  }

  isValid, err = token.ValidateToken("refresh", payload.RefreshToken, false)
  if err != nil {
    return nil, model.ErrUnauthorized
  }
  if !isValid {
    return nil, model.ErrUnauthorized
  }

  accessEmail, err := token.GetEmailFromToken(payload.AccessToken, true)
  if err != nil {
    return nil, model.ErrUnauthorized
  }
  refreshEmail, err := token.GetEmailFromToken(payload.RefreshToken, false)
  if err != nil {
    return nil, model.ErrUnauthorized
  }
  if accessEmail != refreshEmail {
    return nil, model.ErrUnauthorized
  }

  userEntity, err := h.userStore.FindByEmail(accessEmail)
  if err != nil {
    return nil, model.ErrUnauthorized
  }
  if userEntity == nil {
    return nil, model.ErrUnauthorized
  }
  user := userEntity.ConvertToModel()

  tokens, err := token.CreateTokens(user)
  if err != nil {
    return nil, model.ErrUnknown
  }

  // TODO: after refreshToken is used, it should be invalidated in the redis

  return RefreshResponse {
    User: user,
    Tokens: *tokens,
  }, nil
}

func (h *UserHandler) me(r *http.Request) (any, error) {
	var sessionUser model.User
	if err := getUser(r, &sessionUser); err != nil {
		return nil, err
	}

	log.Debug().Interface("sessionUser", sessionUser).Msg("authorized")

	return sessionUser, nil
}
