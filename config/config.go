package config

import (
	"log"
	"os"
	"strconv"

	"github.com/jiaqi-yin/go-url-shortener/services"
)

type Config struct {
	ServerAddr       string
	ShortlinkService services.ShortlinkServiceInterface
}

func LoadConfig() *Config {
	serverAddr := os.Getenv("SERVER_ADDR")
	if serverAddr == "" {
		serverAddr = ":8080"
	}
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
	return &Config{
		ServerAddr:       serverAddr,
		ShortlinkService: s,
	}
}
