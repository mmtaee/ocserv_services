package command

import (
	"api/internal/model"
	"api/pkg/postgres"
	"log"
)

func Migrate() {
	engine := postgres.GetEngine()
	err := engine.Sync(
		&model.Admin{},
		&model.Permission{},
		&model.User{},
		&model.TrafficStatistics{},
	)
	if err != nil {
		log.Fatalf("error sync tables: %v", err)
	}
	log.Println("migrating tables successfully")
}
