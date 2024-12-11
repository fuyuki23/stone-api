package app

import (
	"stone-api/api"
	"stone-api/internal/cache"
	"stone-api/internal/db"
	"stone-api/internal/web"

	"github.com/pkg/errors"
)

type App struct {
	serv *web.Server
	api  *api.API
}

func New() (*App, error) {
	dbConn, err := db.Init()
	if err != nil {
		return nil, errors.Wrap(err, "failed to initialize db")
	}

	store := db.NewStore(dbConn)

	cacheManager := cache.New(cache.Init())

	serv := web.NewServer(store, cacheManager)
	localAPI := api.NewAPI(serv)

	app := &App{
		serv: serv,
		api:  localAPI,
	}

	return app, nil
}

func (app *App) Serve() {
	app.serv.Start()
}
