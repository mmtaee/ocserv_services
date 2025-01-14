package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mmtaee/go-oc-utils/logger"
	"os"
	"strings"
	"sync"
)

type Config struct {
	APP      APP
	DB       DB
	RabbitMQ RabbitMQ
}

type APP struct {
	Debug          bool
	SecretKey      string
	Host           string
	Port           string
	AllowOrigins   []string
	InitSecretFile string
	Isolate        bool
}

type DB struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type RabbitMQ struct {
	Host     string
	Port     string
	User     string
	Password string
	Protocol string
	Vhost    string
}

var (
	config  Config
	AppInit bool
	mutex   sync.RWMutex
)

func GetAppInit() bool {
	mutex.RLock()
	defer mutex.RUnlock()
	return AppInit
}

func ActiveAppInit() {
	mutex.Lock()
	defer mutex.Unlock()
	AppInit = true
}
func Set(debug bool) {
	if debug {
		err := godotenv.Load()
		if err != nil {
			logger.Log(logger.CRITICAL, fmt.Sprintf("Error loading .env file: %v", err))
		}
	}
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		secretKey = "SECRET_KEY122456"
		logger.Log(
			logger.WARNING,
			fmt.Sprintf("SECRET_KEY environment variable not set. set default secret key to: %s", secretKey),
		)
	}

	InitSecretFile := os.Getenv("SECRET_KEY_FILE_NAME")
	if InitSecretFile == "" {
		InitSecretFile = "./init_secret"
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	config.APP = APP{
		Debug:          debug,
		Host:           host,
		Port:           port,
		SecretKey:      secretKey,
		InitSecretFile: InitSecretFile,
	}

	isolate := os.Getenv("ISOLATE")
	if isolate == "" {
		config.APP.Isolate = false
	} else {
		config.APP.Isolate = true
	}

	allowOrigins := os.Getenv("ALLOW_ORIGINS")
	if allowOrigins != "" {
		config.APP.AllowOrigins = strings.Split(allowOrigins, ",")
	}

	config.DB = DB{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     os.Getenv("POSTGRES_PORT"),
		Name:     os.Getenv("POSTGRES_NAME"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
	}

	config.RabbitMQ = RabbitMQ{
		Host:     os.Getenv("RABBIT_MQ_HOST"),
		Port:     os.Getenv("RABBIT_MQ_PORT"),
		User:     os.Getenv("RABBIT_MQ_USER"),
		Password: os.Getenv("RABBIT_MQ_PASSWORD"),
	}
	if os.Getenv("RABBIT_MQ_SECURE") == "true" {
		config.RabbitMQ.Protocol = "amqps"
	} else {
		config.RabbitMQ.Protocol = "amqp"
	}
	if vhost := os.Getenv("RABBIT_MQ_VHOST"); vhost != "" {
		config.RabbitMQ.Vhost = vhost
	}
	logger.Log(logger.INFO, "Configuration applied successfully")
}

func GetDB() *DB {
	return &config.DB
}

func GetRSN() string {
	return fmt.Sprintf(
		"%s://%s:%s@%s:%s/%s",
		config.RabbitMQ.Protocol, config.RabbitMQ.User, config.RabbitMQ.Password,
		config.RabbitMQ.Host, config.RabbitMQ.Port, config.RabbitMQ.Vhost,
	)
}

func GetApp() *APP {
	return &config.APP
}
