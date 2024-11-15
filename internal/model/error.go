package model

type AppError struct {
	status  int
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e AppError) Error() string {
	return e.Message
}

func (e AppError) Status() int {
	return e.status
}

func NewError(status int, code, message string) error {
	return AppError{
		status:  status,
		Code:    code,
		Message: message,
	}
}

var ErrBadRequest = NewError(400, "api.common.bad_request", "bad request")
var ErrUnauthorized = NewError(401, "api.common.unauthorized", "unauthorized")
var ErrNotFound = NewError(404, "api.common.not_found", "not found")
var ErrUnknown = NewError(500, "api.common.unknown", "unknown error")
