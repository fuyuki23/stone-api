package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

type HealthHandler struct{}

func (api *API) initHealthAPI(router *mux.Router) {
	api.health = &HealthHandler{}

	router.Handle("/status", api.BaseHandler(api.health.status)).Methods(http.MethodGet).Name("Liveness")
}

func (h *HealthHandler) status(_ *http.Request) (any, error) {
	return "ok", nil
}
