package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *slog.Logger
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "port to listen for request")
	flag.StringVar(&cfg.env, "env", "development", "development | staging | production")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil)) //change for json later

	app := &application{
		config: cfg,
		logger: logger,
	}

	// func (app *application) healthcheckHandler(w )

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", cfg.port),
		Handler:      app.routes(),
		ReadTimeout:  time.Second * 5,
		WriteTimeout: time.Second * 10,
		IdleTimeout:  time.Second * 60,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
	logger.Info("starting server", "port", cfg.port)

	err := srv.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

}
