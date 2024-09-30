package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	DB                      `yaml:"db"`
	HTTPServer              `yaml:"http_server"`
	LogLevel                int           `yaml:"log_level" env-default:"0"`
	GracefulShutdownTimeout time.Duration `yaml:"graceful_shutdown_timeout" env-default:"30s"`
}

type DB struct {
	Host        string `yaml:"host" env-default:"localhost"`
	Port        int    `yaml:"port" env-default:"3306"`
	User        string `yaml:"user" env-default:"root"`
	Password    string `yaml:"password" env-default:""`
	DBName      string `yaml:"db_name" env-default:"pizza_store"`
	Type        string `yaml:"type" env-required:"true"`
	StoragePath string `yaml:"storage_path" env-default:""`
}

type HTTPServer struct {
	Address      string        `yaml:"address" env-default:"localhost:8080"`
	ReadTimeout  time.Duration `yaml:"read_header_timeout" env-default:"5s"`
	WriteTimeout time.Duration `yaml:"write_timeout" env-default:"10s"`
	IdleTimeout  time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
