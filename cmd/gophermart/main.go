package main

import (
	"github.com/ex0rcist/gophermart/internal/config"
	"github.com/ex0rcist/gophermart/internal/logging"
	"github.com/ex0rcist/gophermart/internal/server"
)

func main() {
	logging.Setup()
	logging.LogInfo("starting server...")

	config, err := config.Parse()
	if err != nil {
		logging.LogFatal(err)
	}

	srv, err := server.New(config)
	if err != nil {
		logging.LogFatal(err)
	}

	err = srv.Run()
	if err != nil {
		logging.LogFatal(err)
	}
}
