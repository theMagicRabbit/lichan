package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func (s *state) handlerDownloads(username string) error {
	opts := "opening=true?sort=dateAsc"
	reqUrl := fmt.Sprintf("%s%s%s?%s", s.ApiUrl, "/api/games/user/", username, opts)

	if s.Config.LastGameTime > 0 {
		reqUrl = fmt.Sprintf("%s&since=%d", reqUrl, s.Config.LastGameTime)
	}

	fmt.Println(reqUrl)
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
		err := game.WriteGame(s)
		if err != nil {
			return err
		}
		s.Config.LastGameTime = game.CreatedAt
	}

	return nil
}
