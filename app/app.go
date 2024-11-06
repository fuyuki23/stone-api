package app

import (
	"stone-api/api"
	"stone-api/internal/db"
	"stone-api/internal/web"
)

type App struct {
	serv *web.Server
	api  *api.Api
}

func New() (*App, error) {
	dbConn, err := db.Init()
	if err != nil {
		return nil, err
	}

	store := db.NewStore(dbConn)

	serv := web.NewServer(store)
	localApi := api.NewApi(serv)

	app := &App{
		serv: serv,
		api:  localApi,
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
