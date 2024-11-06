package web

import (
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"net/http"
	"stone-api/internal/config"
)

type Server struct {
	serv       *http.Server
	BaseRouter *mux.Router

	db *sqlx.DB
}

func NewServer(db *sqlx.DB) *Server {
	server := &Server{
		BaseRouter: mux.NewRouter().StrictSlash(true),
		db:         db,
	}

	return server
}

func (server *Server) DB() *sqlx.DB {
	return server.db
}

func (server *Server) Start() {
	addr := fmt.Sprintf(":%d", config.Get().Server.Port)

	server.serv = &http.Server{
		Addr:    addr,
		Handler: handlers.LoggingHandler(log.Logger, handlers.CompressHandler(server.BaseRouter)),
	}

	server.PrintRoutes()

	log.Info().Str("addr", addr).Msg("server started")
	log.Fatal().Err(server.serv.ListenAndServe()).Send()
}

func (server *Server) PrintRoutes() {
	err := server.BaseRouter.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}
		pathRegexp, err := route.GetPathRegexp()
		if err != nil {
			return err
		}
		methods, err := route.GetMethods()
		if err != nil {
			return nil
			//return err
		}
		for _, method := range methods {
			log.Info().Str("Path", pathRegexp).Msg(fmt.Sprintf("[%s] %s", method, pathTemplate))
		}
		return nil
	})

	if err != nil {
		log.Error().Err(err).Msg("failed to walk through routes")
	}
}
