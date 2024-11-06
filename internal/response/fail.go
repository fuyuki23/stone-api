package response

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"stone-api/internal/model"
)

func Fail(err error) ApiResponse {
	var appError model.AppError
	if errors.As(err, &appError) {
		log.Debug().Err(err).Int("code", appError.Status()).Msg("app error")
		return ApiResponse{data: appError, status: appError.Status()}
	}
	log.Error().Err(err).Msg("unknown error")
	return ApiResponse{data: model.ErrUnknown, status: 500}
}
