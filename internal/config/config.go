package config

import (
	"log"
	"os"

	toml "github.com/pelletier/go-toml"
)

type Config struct {
	Username      []string
	GameDirectory string
	PAT           string
}

func ReadConfig(configPath string) (*Config, error) {
	configData, err := os.ReadFile(configPath)
	if err != nil {
		log.Printf("Error reading config file: %v\n", err)
		return nil, err
	}

	var config Config
	err = toml.Unmarshal(configData, &config)
	if err != nil {
		log.Printf("Error processing config: %v\n", err)
		return nil, err
	}

	return &config, nil
}
