package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/M-Fisher/companies_api/app/api"
	"github.com/M-Fisher/companies_api/app/config"
	"github.com/M-Fisher/companies_api/app/internal/server"
	"github.com/M-Fisher/companies_api/app/internal/storage/postgres"
)

type App struct {
	Cfg      *config.Config
	Server   *server.Server
	Postgres *postgres.DB
}

func New(cfg *config.Config) App {
	srv := server.NewServer(cfg)
	api.SetAPIRouter(srv, "/api/")

	return App{Server: srv}
}

func (a App) Run() error {
	halt := make(chan os.Signal, 1)

	signal.Notify(halt, syscall.SIGTERM, syscall.SIGINT)

	a.Server.Run()
	<-halt
	a.Server.Stop()
	return nil
}
