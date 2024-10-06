package config

import (
	"errors"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
)

type Config struct {
	Env        string `yaml:"env" env-default:"local" env-required:"true"`
	Postgres   `yaml:"postgres"`
	HTTPServer `yaml:"http_server"`
}

type Postgres struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     string `yaml:"port" env-required:"true"`
	Username string `yaml:"username" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DB       string `yaml:"db" env-required:"true"`
}

type HTTPServer struct {
	Address string `yaml:"address" env-default:"0.0.0.0:8080"`
}

func getActualLevelConfig(configPath *string) error {
	configDirPath := os.Getenv("CONFIG_DIRECTORY")
	if configDirPath == "" {
		log.Fatal("CONFIG_DIRECTORY is not set")
	}
	var result error = nil
	if _, err := os.Stat(configDirPath + "/settings.local.yaml"); !os.IsNotExist(err) {
		*configPath = configDirPath + `/settings.local.yaml`
		return result
	}
	if _, err := os.Stat(configDirPath + "/settings.yaml"); !os.IsNotExist(err) {
		*configPath = configDirPath + `/settings.yaml`
		return result
	}
	return errors.New("settings don't exist in path: " + configDirPath)
}

func MustLoad() *Config {
	var configPath string

	err := getActualLevelConfig(&configPath)
	if err != nil {
		log.Fatal("Config files not found", err)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}
	return &cfg
}
