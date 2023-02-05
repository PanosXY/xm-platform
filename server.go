package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/PanosXY/xm-platform/config"
	"github.com/PanosXY/xm-platform/router"
	"github.com/PanosXY/xm-platform/utils/logger"
	"github.com/PanosXY/xm-platform/utils/postgres"
)

type Server interface {
	ListenAndServe()
}

type server struct {
	databaseClient *postgres.PostgresClient
	router         *chi.Mux
	configuration  *config.Configuration
	log            *logger.Logger
}

func NewServer(configuration *config.Configuration, log *logger.Logger) Server {
	db := loadDatabase(configuration, log)

	return &server{
		databaseClient: db,
		router:         router.NewRouter(configuration, log, db),
		configuration:  configuration,
		log:            log,
	}
}

func (srv *server) ListenAndServe() {
	component := "server"

	address := fmt.Sprintf("%s:%d", srv.configuration.HttpServer.Host, srv.configuration.HttpServer.Port)

	go func() {
		if err := http.ListenAndServe(address, srv.router); err != nil {
			srv.log.Fatal(component, "failed to start server", err)
		}
	}()

	srv.log.Info(component, fmt.Sprintf("listening %s", address))
}

func loadDatabase(configuration *config.Configuration, log *logger.Logger) *postgres.PostgresClient {
	component := "server/database"
	ctx := context.Background()

	client, err := postgres.NewPostgresClient(ctx, getPostgresConfiguration(configuration))
	if err != nil {
		log.Fatal(component, "failed to create db client", err)
	}

	if err := client.Connect(ctx); err != nil {
		log.Fatal(component, "failed to connect to db", err)
	}

	return client
}

func getPostgresConfiguration(configuration *config.Configuration) *postgres.Config {
	pool := &postgres.PoolConfig{
		MaxOpenConnections:        configuration.Database.Settings.MaxOpenConns,
		MaxIdleConnections:        configuration.Database.Settings.MaxIdleConns,
		ConnectionMaxLifetimeSecs: configuration.Database.Settings.ConnMaxLifetimeSecs,
		ConnectionMaxIdleTimeSecs: configuration.Database.Settings.ConnMaxIdleTimeSecs,
	}

	return &postgres.Config{
		Pool:               pool,
		Host:               configuration.Database.Host,
		Port:               fmt.Sprintf("%d", configuration.Database.Port),
		Password:           configuration.Database.Password,
		Database:           configuration.Database.Name,
		SSL:                configuration.Database.SSL,
		AppName:            fmt.Sprintf("%s-%s", configuration.App.Name, configuration.App.Version),
		Username:           configuration.Database.Username,
		ConnectTimeoutSecs: configuration.Database.Settings.ConnectTimeout,
	}
}
