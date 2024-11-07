package api

import (
	"net/http"
	"stone-api/internal/web"

	"github.com/gorilla/mux"
)

type Api struct {
	serv *web.Server
	root *mux.Router

	// handlers
	user *UserHandler
}

func NewApi(serv *web.Server) *Api {
	api := &Api{
		serv: serv,
	}

	// root
	api.root = api.serv.BaseRouter

	api.initUserApi(api.root.PathPrefix("/users").Subrouter())

	api.root.NotFoundHandler = http.HandlerFunc(NotFound)
	api.root.MethodNotAllowedHandler = http.HandlerFunc(NotFound)

	return api
}
