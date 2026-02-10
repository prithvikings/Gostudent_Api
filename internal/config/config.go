package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type HTTPServerConfig struct {
	Addr string `yaml:"address"`
}

type Config struct {
	Env              string `yaml:"env" env:"ENV" env-required:"true"`
	StoragePath      string `yaml:"storage_path" env-required:"true"`
	HTTPServerConfig `yaml:"http_server"`
}

func MustLoad() *Config {
	var configPath string

	configPath = os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configFlag := flag.String("config", "", "path to the configuration file")
		flag.Parse()

		configPath = *configFlag

		if configPath == "" {
			log.Fatal("config path is not provided")
		}
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("failed to read config file: %s", err.Error())
	}

	return &cfg
}
