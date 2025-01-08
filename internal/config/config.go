package config

import (
  "os"
  "github.com/joho/godotenv"
)

type Config struct {
  Port   string
  DBFile string
}

func NewConfig() (*Config, error) {

  if err := godotenv.Load(); err != nil {
    return nil, err 
  }

  cfg := &Config{}
  
  cfg.Port = os.Getenv("TODO_PORT")
  cfg.DBFile = os.Getenv("TODO_DBFILE")

  return cfg, nil
}
