package response

import (
	"stone-api/internal/model"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

func Fail(err error) APIResponse {
	var appError model.AppError
	if errors.As(err, &appError) {
		log.Debug().Err(err).Int("code", appError.Status()).Msg("app error")
		return APIResponse{data: appError, status: appError.Status()}
	}
	log.Error().Err(err).Msg("unknown error")
	return APIResponse{data: model.ErrUnknown, status: 500}
}
