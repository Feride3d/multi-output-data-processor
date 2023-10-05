package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"multi-output-data-processor/internal/config"
	"multi-output-data-processor/internal/handler"
	"multi-output-data-processor/internal/server"
	"multi-output-data-processor/internal/service"
)

var configFile string

// Init parses config file. Init run before the main function.
func init() {
	flag.StringVar(&configFile, "config", "./config/config.yml", "path to configuration file")
}

func main() {

	flag.Parse()
	cfg, err := config.NewConfig(configFile)
	if err != nil {
		log.Println("failed to start new config")
	}

	service := service.NewPipelineService(cfg)
	handler := handler.NewHandler(service, cfg)
	server := new(server.Server)

	go func() {
		err := server.Run(cfg.Http, handler.InitRoutes())
		if err != nil {
			log.Printf("failed to run the http server: %v", err.Error())
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	err = server.Shutdown(context.Background())
	if err != nil {
		log.Printf("failed to shut down the http server: %v", err.Error())
	}
}
