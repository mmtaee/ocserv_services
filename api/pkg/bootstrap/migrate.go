package bootstrap

import (
	"api/internal/command"
	"api/pkg/config"
	"api/pkg/postgres"
)

func Migrate(debug bool) {
	config.Set(debug)
	postgres.Connect()
	command.Migrate()
}
