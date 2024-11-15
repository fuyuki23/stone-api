package web

import (
	"fmt"
	"net/http"
	"stone-api/internal/config"
	"stone-api/internal/db"
	"stone-api/internal/utils"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type Server struct {
	serv       *http.Server
	BaseRouter *mux.Router

	store *db.Store
}

func NewServer(store *db.Store) *Server {
	server := &Server{
		BaseRouter: mux.NewRouter().StrictSlash(true),
		store:      store,
	}

	return server
}

func (server *Server) DB() *sqlx.DB {
	return server.store.DB()
}

func (server *Server) Store() *db.Store {
	return server.store
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
	err := server.BaseRouter.Walk(func(route *mux.Route, _ *mux.Router, _ []*mux.Route) error {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}
		methods, err := route.GetMethods()
		if err != nil {
			return nil
		}
		name := route.GetName()
		for _, method := range methods {
			log.Info().Msgf("%8s %s [%s]", utils.AppendString("[", method, "]"), pathTemplate, name)
		}
		return nil
	})

	if err != nil {
		log.Error().Err(err).Msg("failed to walk through routes")
	}
}
