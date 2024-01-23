package main

import (
	"effective_mobile/api/http"
	"effective_mobile/config"
	"effective_mobile/internal"
	"effective_mobile/pkg/logger"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	log, err := logger.Create()
	if err != nil {
		panic(err)
	}
	cfg, err := config.LoadConfig("./config")
	if err != nil {
		log.Fatal(fmt.Sprintf("load config file: %v", err))
	}

	app := internal.NewApp(&cfg, log)
	err = app.Init()
	if err != nil {
		log.Fatal(fmt.Sprintf("init app: %v", err))
	}

	server := http.NewServer(&cfg, log)
	server.Init()
	server.MapHandlers(app)
	go func() {
		server.Serve()
	}()

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt, syscall.SIGTERM)
	sig := <-terminate
	log.Info("received terminate signal", zap.String("signal", sig.String()))

	server.Shutdown()
	app.DB.Close()
}
