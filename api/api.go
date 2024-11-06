package api

import (
	"github.com/gorilla/mux"
	"net/http"
	"stone-api/internal/web"
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

	api.root.NotFoundHandler = http.HandlerFunc(web.NotFound)
	api.root.MethodNotAllowedHandler = http.HandlerFunc(web.NotFound)

	return api
}
