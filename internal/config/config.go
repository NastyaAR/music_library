package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Logger `yaml:"logger"`
	Db     `yaml:"postgres"`
}

type Logger struct {
	LogLevel string `yaml:"log_level"`
	LogFile  string `yaml:"log_file"`
}

type Db struct {
	Host         string `yaml:"host" env:"POSTGRES_HOST"`
	Port         int    `yaml:"port"`
	User         string `yaml:"user" env:"POSTGRES_USER"`
	Password     string `yaml:"password" env:"POSTGRES_PASSWORD"`
	Name         string `yaml:"name" env:"POSTGRES_DB"`
	DbTimeoutSec int    `yaml:"db_timeout_sec"`
}

func ReadConfig(configPath string) (*Config, error) {
	cfg := Config{}

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("read config error: %v", err.Error())
	}

	err = cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("read config error: %v", err.Error())
	}

	return &cfg, nil
}
