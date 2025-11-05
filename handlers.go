package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (s *state) handlerDownloads(username string) error {
	opts := "opening=true&sort=dateAsc"
	reqUrl := fmt.Sprintf("%s%s%s?%s", s.ApiUrl, "/api/games/user/", username, opts)

	if s.Config.LastGameTime > 0 {
		reqUrl = fmt.Sprintf("%s&since=%d", reqUrl, s.Config.LastGameTime)
	}

	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Add("authorization", fmt.Sprintf("Bearer %s", s.Config.PAT))
	req.Header.Add("accept", "application/x-ndjson")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var games []Game
	gamesScaner := bufio.NewScanner(res.Body)

	for gamesScaner.Scan() {
		gameBytes := gamesScaner.Bytes()
		game := Game{}
		err = json.Unmarshal(gameBytes, &game)
		if err != nil {
			log.Printf("Error unmarshaling game: %v\n", err)
			continue
		}
		games = append(games, game)
	}

	err = gamesScaner.Err()
	if err != nil {
		log.Printf("Error scanning for games: %v\n", err)
	}

	for _, game := range games {
		outputDir := filepath.Join(s.Config.GameDirectory, username)
		err := game.WriteGame(s, outputDir)
		if err != nil {
			return err
		}
		s.Config.LastGameTime = game.CreatedAt
	}

	return nil
}

func (s *state) handlerAnalyze(username string) error {
	log.Printf("Processing games for %s previously downloaded.", username)
	userGames := filepath.Join(s.Config.GameDirectory, username)
	engineGames := filepath.Join(s.Config.EngineDirectory, username)

	// Get existing files
	files, err := os.ReadDir(userGames)
	if err != nil {
		log.Printf("Unable to read user game directory: %v\n", err)
		return err
	}

	// START LOOP
	for _, file := range files {
		gameFile := file.Name()
		if file.IsDir() || strings.ToLower(filepath.Ext(gameFile)) != ".pgn"  {
			continue
		}
		gamePath := filepath.Join(userGames, gameFile)
		engineFile := strings.ToLower(strings.TrimSuffix(gameFile, filepath.Ext(gameFile))+"_stockfish.pgn")
		enginePath := filepath.Join(engineGames, engineFile)
		_, err := os.Stat(enginePath)
		if err == nil {
			// if the file exists, assume that the game has already been processed
			continue
		} else {
			if !errors.Is(err, fs.ErrNotExist) {
				// If the error is anything other than the file not existing, log the error and skip
				log.Printf("Error accessing engine path: %v\n", err)
				continue
			}
		}

		gamePGNBytes, err := os.ReadFile(gamePath)
		if err != nil {
			log.Printf("Error reading game PNG file: %v\n", err)
			continue
		}

		game, err := GameFromPGN(gamePGNBytes)
		if err != nil {
			return err
		}

		for ms := range strings.SplitSeq(game.Moves, " ") {
			move, err := ParseMoveString(ms)
			if err != nil {
				log.Printf("Unable to parse move: %v\n", err)
				break
			}
			fmt.Println(*move)
		}

	// Feed game to stockfish while processing results
	// Write output to processed file
	}
	// STOP LOOP
	return nil
}


