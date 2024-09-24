package main

import (
	"github.com/ex0rcist/gophermart/internal/app"
	"github.com/ex0rcist/gophermart/internal/config"
	"github.com/ex0rcist/gophermart/internal/logging"
)

func main() {
	logging.Setup()
	logging.LogInfo("starting server...")

	config, err := config.Parse()
	if err != nil {
		logging.LogFatal(err)
	}

	apl, err := NewApp(config)
	if err != nil {
		logging.LogFatal(err)
	}

	err = apl.Run()
	if err != nil {
		logging.LogFatal(err)
	}
}

func NewApp(config *config.Config) (*app.App, error) {
	return app.New(config, nil, nil, nil)
}
