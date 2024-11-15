package main

import (
	"os"
	"stone-api/app"
	"stone-api/internal/config"
	"time"
	_ "time/tzdata"

	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		panic(errors.Wrap(err, "failed to load location"))
	}
	time.Local = loc

	initLogger()
	initConfig()

	localApp, err := app.New()
	if err != nil {
		panic(errors.Wrap(err, "failed to initialize app"))
	}

	localApp.Serve()
}

func initLogger() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339})
}

func initConfig() {
	if err := config.Load(); err != nil {
		panic(errors.Wrap(err, "failed to load config"))
	}
}
