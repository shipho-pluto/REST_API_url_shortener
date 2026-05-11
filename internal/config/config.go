package config

import (
	"os"
	"time"

	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env     string  `yaml:"env" env:"ENV" env-default:"local"`
	Storage Storage `yaml:"storage" env:"STORAGE"`
	Cache   Cache   `yaml:"cache" env:"CACHE"`
	Server  Server  `yaml:"http_server" env:"HTTP_SERVER"`
}

type Cache struct {
	Addr        string        `yaml:"addr"`
	Password    string        `yaml:"password"`
	DB          int           `yaml:"db"`
	MaxRetries  int           `yaml:"max_retries"`
	DialTimeout time.Duration `yaml:"dial_timeout"`
	Timeout     time.Duration `yaml:"timeout" env:"TIMEOUT"`
}

type Storage struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port" env:"PORT"`
	User     string `yaml:"user" env:"USER"`
	DBName   string `yaml:"dbname" env:"DBNAME"`
	Password string `yaml:"password" env:"PASSWORD"`
	SSLMode  string `yaml:"sslmode"`
}

type Server struct {
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env:"PASSWORD" env-required:"true"`
	Address     string        `yaml:"address" env:"Address"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func Mustload() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH environment is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("File does not exist: %s", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Cannot read config settings from file: %s", configPath)
	}

	return &cfg
}
