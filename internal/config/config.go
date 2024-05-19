package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

const configPath = "config/main.yaml"

type Config struct {
	Port int    `yaml:"port"`
	Env  string `yaml:"env"`
	Salt string `yaml:"salt"`
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
