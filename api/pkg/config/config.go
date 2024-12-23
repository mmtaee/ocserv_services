package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strings"
)

type Config struct {
	APP      APP
	DB       DB
	RabbitMQ RabbitMQ
}

type APP struct {
	Debug        bool
	SecretKey    string
	Host         string
	Port         string
	AllowOrigins []string
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

var config Config

func Set(debug bool) {
	if debug {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}
	secretKey := os.Getenv("SECRET_KEY")
	if secretKey == "" {
		secretKey = "SECRET_KEY122456"
		log.Println("SECRET_KEY environment variable not set. set default secret key to: " + secretKey)
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
		Debug:     debug,
		SecretKey: secretKey,
		Host:      host,
		Port:      port,
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
	log.Println("Configuration applied successfully")
}

func GetDSN() string {
	if config.APP.Debug {
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.Name)
	}
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.DB.Host, config.DB.Port, config.DB.User, config.DB.Password, config.DB.Name,
	)
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
