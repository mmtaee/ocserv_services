package handlers

import (
	"api/internal/models"
	"github.com/mmtaee/go-oc-utils/database"
	"log"
)

var tables = []interface{}{
	&models.User{},
	&models.UserPermission{},
	&models.UserToken{},
	&models.PanelConfig{},
	&models.OcUser{},
	&models.OcUserActivity{},
	&models.OcUserTrafficStatistics{},
}

func Migrate() {
	engine := database.Connection()
	err := engine.AutoMigrate(tables...)
	if err != nil {
		log.Fatalf("error sync tables: %v", err)
	}
	log.Println("migrating tables successfully")
}

func Drop() {
	engine := database.Connection()
	err := engine.Migrator().DropTable(tables...)
	if err != nil {
		return
	}
	log.Println("database drop table successfully")
}
