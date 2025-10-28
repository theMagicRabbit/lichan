package config

import (
	"errors"
	"log"
	"os"
	"path"
	"strings"

	toml "github.com/pelletier/go-toml/v2"
)

type Config struct {
	Username      []string `toml:"username"`
	GameDirectory string   `toml:"game_directory"`
	PAT           string   `toml:"token"`
	LastGameTime  int64    `toml:"last_run"`
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

	if len(config.Username) < 1 {
		log.Println("No usernames provided")
		return nil, errors.New("No usernames provided")
	}

	if config.GameDirectory == "" {
		log.Println("No game directory provided")
		return nil, errors.New("No game directory provided")
	}

	newPath, err := replaceTilde(config.GameDirectory)
	if err != nil {
		log.Printf("Unable to clean path: %v\n", err)
		return nil, err
	}

	config.GameDirectory = newPath

	return &config, nil
}

func replaceTilde(p string) (string, error) {
	if !strings.HasPrefix(p, "~") {
		return p, nil
	}
	
	userHome, err := os.UserHomeDir()
	if err != nil {
		return p, err
	}

	replacementPath := strings.Replace(p, "~", userHome+"/", 1)
	replacementPath = path.Clean(replacementPath)

	return replacementPath, nil
}

func (C *Config) WriteConfig(configPath string) error {
	configBytes, err := toml.Marshal(C)
	if err != nil {
		return err
	}
	err = os.WriteFile(configPath, configBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (C *Config) CreateDirs() error {
	for _, user := range C.Username {
		p := path.Join(C.GameDirectory, user)
		err := os.MkdirAll(p, 0755)
		if err != nil {
			log.Printf("Error creating directory for games: %v\n", err)
			return err
		}
	}
	return nil
}
