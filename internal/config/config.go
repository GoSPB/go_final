package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port   string
	DBFile string
}

func (c *Config) Init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Ошибка чтения из .env")
	}

	c.Port = strings.Join([]string{":", os.Getenv("TODO_PORT")}, "")
	c.DBFile = os.Getenv("TODO_DBFILE")
}
