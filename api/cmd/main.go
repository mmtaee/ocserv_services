package main

import (
	_ "api/docs"
	"api/internal/handlers"
	"api/pkg/config"
	"api/pkg/database"
	"api/pkg/routing"
	"flag"
	"log"
)

// @title Ocserv User management Example Api
// @version 1.0
// @description This is a sample Ocserv User management Api server.
// @BasePath /api
func main() {
	var (
		debug   bool
		migrate bool
	)
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.BoolVar(&migrate, "migrate", false, "migrate models to database")
	flag.Parse()
	if debug {
		log.SetFlags(0)
	}
	config.Set(debug)
	database.Connect()
	if migrate {
		handlers.Migrate()
	} else {
		routing.Serve()
	}
}
