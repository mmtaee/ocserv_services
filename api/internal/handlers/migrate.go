package handlers

import (
	"api/internal/models"
	"api/pkg/database"
	"log"
)

func Migrate() {
	engine := database.Connection()
	tables := []interface{}{
		&models.User{},
		&models.UserPermission{},
		&models.UserToken{},
		&models.PanelConfig{},
		&models.OcUser{},
		&models.OcUserActivity{},
		&models.OcUserTrafficStatistics{},
	}

	//errx := engine.Migrator().DropTable(tables...)
	//if errx != nil {
	//	return
	//}

	err := engine.AutoMigrate(tables...)
	if err != nil {
		log.Fatalf("error sync tables: %v", err)
	}
	log.Println("migrating tables successfully")
}
