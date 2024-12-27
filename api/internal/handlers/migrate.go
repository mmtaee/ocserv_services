package handlers

import (
	"api/internal/models"
	"api/pkg/database"
	"log"
)

func Migrate() {
	engine := database.Connection()
	err := engine.AutoMigrate(
		&models.Admin{},
		&models.Permission{},
		&models.PanelConfig{},
		&models.User{},
		&models.TrafficStatistics{},
	)
	if err != nil {
		log.Fatalf("error sync tables: %v", err)
	}
	log.Println("migrating tables successfully")
}
