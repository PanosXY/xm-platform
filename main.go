package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/PanosXY/xm-platform/config"
	"github.com/PanosXY/xm-platform/utils/logger"
)

func main() {
	configuration, err := config.LoadConfiguration()
	if err != nil {
		fmt.Println("failed to load configuration:", err.Error())
		os.Exit(1)
	}

	component := "main"

	log := logger.NewLogger(configuration.App.Env != "production", false)
	defer log.Info(component, "exit")

	log.SetDefaultField("application_name", configuration.App.Name)
	log.SetDefaultField("version", configuration.App.Version)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	log.Info(component, "initializing server...")
	server := NewServer(configuration, log)

	log.Info(component, "starting server...")
	server.ListenAndServe()

	<-shutdown
}
