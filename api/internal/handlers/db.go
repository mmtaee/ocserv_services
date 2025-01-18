package handlers

import (
	"fmt"
	"github.com/mmtaee/go-oc-utils/database"
	"github.com/mmtaee/go-oc-utils/logger"
	"github.com/mmtaee/go-oc-utils/models"
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
		logger.Log(logger.CRITICAL, fmt.Sprintf("error sync tables: %v", err))
	}
	logger.Log(logger.INFO, "migrating tables successfully")
}

func Drop() {
	engine := database.Connection()
	err := engine.Migrator().DropTable(tables...)
	if err != nil {
		return
	}
	logger.Log(logger.INFO, "database drop table successfully")
}
