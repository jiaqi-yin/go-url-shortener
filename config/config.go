package config

import (
	"log"
	"os"
	"strconv"

	"github.com/jiaqi-yin/go-url-shortener/services"
)

type Config struct {
	ShortlinkService services.ShortlinkServiceInterface
}

func LoadConfig() *Config {
	addr := os.Getenv("APP_REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	password := os.Getenv("APP_REDIS_PASSWORD")
	if password == "" {
		password = ""
	}
	dbString := os.Getenv("APP_REDIS_DB")
	if dbString == "" {
		dbString = "0"
	}
	db, err := strconv.Atoi(dbString)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connecting to redis")

	s := services.NewShortlinkService(addr, password, db)
	return &Config{ShortlinkService: s}
}
