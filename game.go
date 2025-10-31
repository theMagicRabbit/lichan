package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type GameResult string
const (
	WhiteWins GameResult = "1-0"
	BlackWins GameResult = "0-1"
	Draw      GameResult = "1/2-1/2"
	Unknown   GameResult = "*"
)

type Game struct {
	ID         string `json:"id"`
	Rated      bool   `json:"rated"`
	Variant    string `json:"variant"`
	Speed      string `json:"speed"`
	Perf       string `json:"perf"`
	CreatedAt  int64  `json:"createdAt"`
	LastMoveAt int64  `json:"lastMoveAt"`
	Status     string `json:"status"`
	Source     string `json:"source"`
	Players    struct {
		White struct {
			User struct {
				Name string `json:"name"`
				ID   string `json:"id"`
			} `json:"user"`
			Rating      int  `json:"rating"`
			RatingDiff  int  `json:"ratingDiff"`
			Provisional bool `json:"provisional"`
		} `json:"white"`
		Black struct {
			User struct {
				Name string `json:"name"`
				ID   string `json:"id"`
			} `json:"user"`
			Rating     int `json:"rating"`
			RatingDiff int `json:"ratingDiff"`
		} `json:"black"`
	} `json:"players"`
	FullID  string `json:"fullId"`
	Winner  string `json:"winner"`
	Opening struct {
		Eco  string `json:"eco"`
		Name string `json:"name"`
		Ply  int    `json:"ply"`
	} `json:"opening"`
	Moves string `json:"moves"`
	Clock struct {
		Initial   int `json:"initial"`
		Increment int `json:"increment"`
		TotalTime int `json:"totalTime"`
	} `json:"clock"`
}

func (g *Game) WriteGame(s *state, outputDir string) error {
	gameString, err := GameToPGN(g, s.SiteUrl)
	if err != nil {
		return err
	}

	gameYear, gameMonth, gameDay := time.UnixMilli(g.CreatedAt).Date()
	gameDate := fmt.Sprintf("%d.%d.%d", gameYear, gameMonth, gameDay)

	fileTitle := fmt.Sprintf("%s_%s.pgn", gameDate, g.ID)
	gameFilePath := fmt.Sprintf("%s/%s", outputDir, fileTitle)
	err = os.WriteFile(gameFilePath, []byte(gameString), 0644)
	if err != nil {
		return err
	}
	log.Printf("Wrote %s\n", gameFilePath)
	return nil
}

func GameFromPGN(data []byte) (Game, error) {
	if bytes.ContainsAny(data, "\t") {
		return Game{}, errors.New("PGN format must not contain any tabs")
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(tokenizerPGN)

	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}
	game := Game{}
	return game, nil
}

func tokenizerPGN(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return
	}
	for advance = 1; advance < len(data); advance++ {
		nextByte := string(data[advance])
		if nextByte == "[" || nextByte == "]" {
			token = append(token, data[advance])
			break
		}

		if nextByte == "\"" {
			break
		}
		
		if nextByte == "\n" {
			continue
		}

		token = append(token, data[advance])
	}
	if len(token) > 0 && string(token[0]) == "\"" {
		err = errors.New("Malformed PGN: Unmatch quotation mark")
	}
	return
}

func GameToPGN(game *Game, url string) (string, error) {
	pgnTemplate := `[Event "%s"]
[Site "%s/%s"]
[Date "%s"]
[Round "-"]
[White "%s"]
[Black "%s"]
[Result "%s"]
[GameId "%s"]
[WhiteElo "%d"]
[BlackElo "%d"]
[Opening "%s"]
[TimeControl "%s"]

%s
`
	var event string
	if game.Rated {
		event = fmt.Sprintf("%s %s game", "rated", game.Speed)
	} else {
		event = fmt.Sprintf("%s %s game", "unrated", game.Speed)
	}

	gameYear, gameMonth, gameDay := time.UnixMilli(game.CreatedAt).Date()
	gameDate := fmt.Sprintf("%d.%d.%d", gameYear, gameMonth, gameDay)

	var result GameResult
	switch game.Winner {
	case "black":
		result = BlackWins
	case "white":
		result = WhiteWins
	case "draw":
		result = Draw
	default:
		result = Unknown
	}

	gameTimeControl := fmt.Sprintf("%d +%d", game.Clock.Initial, game.Clock.Increment)

	moveSlice := strings.Split(game.Moves, " ")
	whiteMove := true
	var moveString string 
	moveCounter := 1
	for _, move := range moveSlice {
		if !whiteMove {
			moveString = fmt.Sprintf("%s %s", moveString, move)
			moveCounter++
		} else {
			if moveCounter != 1 {
				moveString = fmt.Sprintf("%s %d. %s", moveString, moveCounter, move)
			} else {
				moveString = fmt.Sprintf("%d. %s", moveCounter, move)
			}
		}
		whiteMove = !whiteMove
	}


	gamePGN := fmt.Sprintf(pgnTemplate,
				event,
				url,
				game.ID,
				gameDate,
				game.Players.White.User.Name,
				game.Players.Black.User.Name,
				result,
				game.ID,
				game.Players.White.Rating,
				game.Players.Black.Rating,
				game.Opening.Name,
				gameTimeControl,
				moveString,
	)
	return gamePGN, nil
}
