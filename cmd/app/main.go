package main

import (
	"awesomeProject/pkg/config"
	"awesomeProject/pkg/di"
	"log"
)

func main() {
	cfg, configErr := config.LoadConfig()
	if configErr != nil {
		log.Fatal("cannot load config: ", configErr)
	}

	server, diErr := di.InitializeAPI(cfg)
	if diErr != nil {
		log.Fatal("cannot start server: ", diErr)
	} else {
		server.Start()
	}
}
