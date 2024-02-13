package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/segment3d-app/segment3d-storage/util"
	"github.com/segment3d-app/segment3d-storage/api"
	_ "github.com/swaggo/files"
	_ "github.com/swaggo/gin-swagger"
)

// @title Segment3d App API Documentation
// @version 1.0
// @description This is a documentation for Segment3d App API

// @host localhost:8081
// @BasePath /

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("can't load config")
	}

	server, err := api.NewServer(&config)
	if err != nil {
		log.Fatal("can't create server", err)
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("can't start server", err)
	}
}
