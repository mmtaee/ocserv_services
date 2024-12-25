package bootstrap

import (
	"api/internal/command"
	"api/pkg/config"
	"api/pkg/database"
)

func Migrate(debug bool) {
	config.Set(debug)
	database.Connect()
	command.Migrate()
}
