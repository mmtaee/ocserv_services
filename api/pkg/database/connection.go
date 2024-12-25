package database

import (
	"api/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

var db *gorm.DB

func Connect() {
	dsn := config.GetDSN()
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		err = sqlDB.Ping()
		if err == nil {
			break
		}
		log.Printf("Database is not ready: %v\n", err)
		time.Sleep(1 * time.Second)
	}
	if err != nil {
		log.Fatalf("postgres connection failed: %v", err)
	}
	log.Println("Database connection established")
}

func Connection() *gorm.DB {
	cfg := config.GetApp()
	if cfg.Debug {
		db.Debug()
		db.Logger = logger.Default.LogMode(logger.Info)
	}
	return db
}

func Close() {
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	err = sqlDB.Close()
	if err != nil {
		log.Fatal(err)
	}
}
