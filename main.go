package main

import (
	"log"
	"os"
	"path"

	"github.com/theMagicRabbit/lichan/internal/config"
)

type state struct {
	Config  *config.Config
	ApiUrl  string
	SiteUrl string
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
	defer config.WriteConfig(configFile)

	state := state{
		Config:  config,
		ApiUrl:  "https://lichess.org",
		SiteUrl: "https://lichess.org",
	}

	err = state.Config.CreateDirs()
	if err != nil {
		log.Fatal(err)
	}
	for _, user := range state.Config.Username {
		state.handlerDownloads(user)
		state.handlerAnalyze(user)
	}
}
