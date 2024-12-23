package bootstrap

import (
	"api/pkg/config"
	"api/pkg/postgres"
	"api/pkg/rabbitmq"
	"api/pkg/routing"
)

func Serve(debug bool) {
	config.Set(debug)
	postgres.Connect()
	rabbitmq.Connect()
	routing.Serve()
}
