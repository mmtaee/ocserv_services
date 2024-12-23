package postgres

import (
	"api/pkg/config"
	_ "github.com/lib/pq"
	"log"
	"time"
	"xorm.io/xorm"
)

var engine *xorm.Engine

func Connect() {
	dsn := config.GetDSN()
	var err error
	engine, err = xorm.NewEngine("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}
	sqlDB := engine.DB()
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

func GetEngine() *xorm.Engine {
	cfg := config.GetApp()
	if cfg.Debug {
		engine.ShowSQL(true)
	}
	return engine
}

func Close() {
	sqlDB := engine.DB()
	err := sqlDB.Close()
	if err != nil {
		log.Fatal(err)
	}
}
