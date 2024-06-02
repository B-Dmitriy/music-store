package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const configPath = "config/main.yaml"

type Config struct {
	Port      int    `yaml:"port"`
	Env       string `yaml:"env"`
	PassCost  int    `yaml:"password_cost"`
	SecretKey string `yaml:"secret_key"`
}

func MustReadConfig() *Config {
	config := Config{}

	file, err := os.ReadFile(configPath)
	if err != nil {
		panic(fmt.Errorf("config reading error: %s", err))
	}

	err = yaml.Unmarshal(file, &config)
	if err != nil {
		panic(fmt.Errorf("config reading error: %s", err))
	}

	return &config
}
