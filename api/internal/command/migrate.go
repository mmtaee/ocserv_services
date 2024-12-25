package command

import (
	"api/internal/model"
	"api/pkg/database"
	"log"
)

func Migrate() {
	engine := database.Connection()
	err := engine.AutoMigrate(
		&model.Admin{},
		&model.Permission{},
		&model.PanelConfig{},
		&model.User{},
		&model.TrafficStatistics{},
	)
	if err != nil {
		log.Fatalf("error sync tables: %v", err)
	}
	log.Println("migrating tables successfully")
}
