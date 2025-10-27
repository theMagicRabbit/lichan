package main

import (
	"log"
	"os"
	"path"

	"github.com/theMagicRabbit/lichan/internal/config"
)

type state struct {
	config *config.Config
}

func main() {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Could not locate user config directory: %v\n", err)
	}

	configFile := path.Join(userConfigDir, "lichan", "config.toml")
	config, err := config.ReadConfig(configFile)
	if err != nil {
		log.Fatalf("Error reading config: %v\n", err)
	}
	
	state := state{
		config: config,
	}

	err = state.config.CreateDirs()
	if err != nil {
		log.Fatal(err)
	}
}
