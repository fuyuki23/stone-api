package app

import (
	"github.com/jmoiron/sqlx"
	"stone-api/api"
	"stone-api/internal/db"
	"stone-api/internal/web"
)

type App struct {
	serv *web.Server
	api  *api.Api
	db   *sqlx.DB
}

func New() (*App, error) {
	db, err := db.Init()
	if err != nil {
		return nil, err
	}

	serv := web.NewServer(db)
	api := api.NewApi(serv)

	app := &App{
		serv: serv,
		api:  api,
		db:   db,
	}

	return app, nil
}

func (app *App) Serve() {
	app.serv.Start()
	//log.Fatal().Err(a.serv.ListenAndServe()).Send()
}

//func beforeExit(db *sqlx.DB) {
//	// close db connection
//	err := db.Close()
//	if err != nil {
//		log.Error().Err(err).Msg("failed to close db connection")
//	}
//}
