package web

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"stone-api/internal/cache"
	"stone-api/internal/config"
	"stone-api/internal/db"
	"stone-api/internal/utils"
	"syscall"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
)

type Server struct {
	serv       *http.Server
	BaseRouter *mux.Router

	store *db.Store
	cache *cache.Manager
}

func NewServer(store *db.Store, cacheManager *cache.Manager) *Server {
	server := &Server{
		BaseRouter: mux.NewRouter().StrictSlash(true),
		store:      store,
		cache:      cacheManager,
	}

	return server
}

func (server *Server) DB() *sqlx.DB {
	return server.store.DB()
}

func (server *Server) Store() *db.Store {
	return server.store
}

func (server *Server) Cache() *cache.Manager {
	return server.cache
}

func (server *Server) Start() {
	addr := fmt.Sprintf(":%d", config.Get().Server.Port)

	ctx, cancel := context.WithCancel(context.Background())

	server.BaseRouter.Use(handlers.RecoveryHandler(handlers.PrintRecoveryStack(true)))
	server.BaseRouter.Use(RequestID)

	server.serv = &http.Server{
		Addr: addr,
		Handler: handlers.LoggingHandler(
			log.Logger,
			handlers.CompressHandler(server.BaseRouter)),
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	server.PrintRoutes()

	log.Info().Str("addr", addr).Msg("server started")
	go func() {
		if err := server.serv.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatal().Err(err).Send()
		}
	}()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	// wait for signal.
	<-signalChan
	log.Info().Msg("shutting down server")
	// set exit code
	// Exit with status 0, 2nd defer.
	defer os.Exit(0)
	// Close database connection after server shutdown, 1st defer.
	defer func(db *sqlx.DB) {
		log.Info().Msg("closing database connection")
		err := db.Close()
		if err != nil {
			log.Error().Err(err).Msg("failed to close database")
		}
	}(server.DB())

	// force shutdown server if it receives another signal within 10 seconds.
	go func() {
		<-signalChan
		log.Fatal().Msg("force shutdown server")
	}()

	// set timeout for graceful shutdown to 10 seconds.
	timeoutCtx, cancelTimeoutCtx := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelTimeoutCtx()

	// graceful shutdown, wait for all connections to close for 10 seconds.
	if err := server.serv.Shutdown(timeoutCtx); err != nil {
		log.Error().Err(err).Msg("failed to shutdown server")
		cancel()
		defer os.Exit(1)
		return
	}
	log.Info().Msg("server stopped")

	cancel()
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
