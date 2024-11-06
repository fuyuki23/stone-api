package web

import (
	"net/http"
	"stone-api/internal/model"
	"stone-api/internal/response"
)

type StoneHandler = func(w http.ResponseWriter, r *http.Request) error

func BaseHandler(handle StoneHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := handle(w, r); err != nil {
			_ = response.Fail(err).Send(w)
		}
	})
}

var NotFound = func(w http.ResponseWriter, r *http.Request) {
	_ = response.Fail(model.ErrNotFound).Status(http.StatusNotFound).Send(w)
}
