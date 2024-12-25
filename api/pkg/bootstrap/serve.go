package bootstrap

import (
	"api/pkg/config"
	"api/pkg/database"
	"api/pkg/rabbitmq"
	"api/pkg/routing"
	"log"
)

func Serve(debug bool) {
	config.Set(debug)
	database.Connect()
	if err := rabbitmq.CheckConnection(); err != nil {
		log.Fatal(err)
	}
	routing.Serve()
}
