package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
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

	valuesMap := make(map[string]string)
	game := Game{}
	for scanner.Scan() {
		key := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if key == "[" || key == "]" {
			continue
		}
		if scanner.Scan() {
			val := strings.TrimSpace(scanner.Text())
			if val == "]" {
				fmt.Printf("missing value for key: %s\n", key)
				break
			}
			valuesMap[key] = val
		}
	}

	for key, val := range valuesMap {
		switch key {
		case "event":
			words := strings.Split(val, " ")
			if len(words) == 3 {
				game.Rated = strings.ToLower(words[0]) == "rated"
				game.Speed = strings.ToLower(words[1])
			}
		case "site":
		// Derivied value, not needed
		case "date":
			dateStrings := strings.Split(val, ".")
			if len(dateStrings) == 3 {
				year, err := strconv.Atoi(dateStrings[0])
				if err != nil {
					log.Printf("Unable to parse date as int: %v\n", err)
					break
				}
				month, err := strconv.Atoi(dateStrings[1])
				if err != nil {
					log.Printf("Unable to parse date as int: %v\n", err)
					break
				}
				day, err := strconv.Atoi(dateStrings[2])
				if err != nil {
					log.Printf("Unable to parse date as int: %v\n", err)
					break
				}
				game.CreatedAt = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).UnixMilli()
			}
		case "white":
			game.Players.White.User.Name = strings.TrimSpace(val)
		case "black":
			game.Players.Black.User.Name = strings.TrimSpace(val)
		case "result":
			result := GameResult(strings.TrimSpace(val))
			switch result {
			case BlackWins:
				game.Winner = "black"
			case WhiteWins:
				game.Winner = "white"
			case Draw:
				game.Winner = "draw"
			default:
			}
		case "gameid":
			game.ID = strings.TrimSpace(val)
		case "opening":
			game.Opening.Name = strings.TrimSpace(val)
		case "whiteelo":
			elo, err := strconv.Atoi(val)
			if err != nil {
				log.Printf("Could not parse rating as int: %v\n", err)
				break
			}
			game.Players.White.Rating = elo
		case "blackelo":
			elo, err := strconv.Atoi(val)
			if err != nil {
				log.Printf("Could not parse rating as int: %v\n", err)
				break
			}
			game.Players.Black.Rating = elo
		case "timecontrol":
			timeStrings := strings.Split(val, "+")
			if len(timeStrings) == 2 {
				initial, err := strconv.Atoi(strings.TrimSpace(timeStrings[0]))
				if err != nil {
					log.Printf("Could not parse time control: %v\n", err)
					break
				}
				increment, err := strconv.Atoi(strings.TrimSpace(timeStrings[1]))
				if err != nil {
					log.Printf("Could not parse time control: %v\n", err)
					break
				}
				game.Clock.Initial = initial
				game.Clock.Increment = increment
			}
		default:
		}
	}
	return game, nil
}

func tokenizerPGN(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return
	}
	for i := 0; i < len(data); i++ {
		advance++
		nextByte := string(data[i])
		if nextByte == "[" || nextByte == "]" {
			token = append(token, data[i])
			break
		}

		if nextByte == "\"" {
			break
		}
		
		if nextByte == "\n" {
			continue
		}

		token = append(token, data[i])
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
