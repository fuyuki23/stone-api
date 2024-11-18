package api

import (
	"net/http"
	"stone-api/internal/web"

	"github.com/gorilla/mux"
)

type API struct {
	serv *web.Server
	root *mux.Router

	health *HealthHandler // health handler.
	user   *UserHandler   // user handler.
	diary  *DiaryHandler  // diary handler.
}

func NewAPI(serv *web.Server) *API {
	api := &API{
		serv: serv,
	}

	// root mux router.
	api.root = api.serv.BaseRouter

	api.initUserAPI(api.root.PathPrefix("/users").Subrouter())
	api.initDiaryAPI(api.root.PathPrefix("/diaries").Subrouter())
	api.initHealthAPI(api.root.PathPrefix("/health").Subrouter())

	api.root.NotFoundHandler = http.HandlerFunc(NotFound)
	api.root.MethodNotAllowedHandler = http.HandlerFunc(NotFound)

	return api
}
