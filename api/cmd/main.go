package main

import (
	_ "api/docs"
	"api/internal/handlers"
	"api/pkg/config"
	"api/pkg/routing"
	"flag"
	"github.com/mmtaee/go-oc-utils/database"
	"os"
	"strconv"
)

// @title Ocserv User management Example Api
// @version 1.0
// @description This is a sample Ocserv User management Api server.
// @BasePath /services
func main() {
	var (
		debug   bool
		migrate bool
		drop    bool
	)
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.BoolVar(&migrate, "migrate", false, "migrate models to database")
	flag.BoolVar(&drop, "drop", false, "drop models table from database")
	flag.Parse()

	if debugStr := os.Getenv("DEBUG"); debugStr != "" {
		debug, _ = strconv.ParseBool(debugStr)
	}

	config.Set(debug)
	dbCfg := config.GetDB()
	dbConfig := &database.DBConfig{
		Host:     dbCfg.Host,
		Port:     dbCfg.Port,
		User:     dbCfg.User,
		Password: dbCfg.Password,
		Name:     dbCfg.Name,
	}
	database.Connect(dbConfig, debug)
	if migrate {
		handlers.Migrate()
	} else if drop && debug {
		handlers.Drop()
	} else {
		routing.Serve()
	}
}
