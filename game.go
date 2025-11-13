package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"
)
var fileTokens []string = []string{"a", "b", "c", "d", "e", "f", "g", "h"}
var rankTokens []string = []string{"1", "2", "3", "4", "5", "6", "7", "8"}
var pieceTokens []string = []string{"K", "Q", "R", "B", "N"}
var longCastle string = "O-O-O"
var shortCastle string = "O-O"
var check string = "+"
var mate string = "#"
var capture string = "x"
var promote string = "="

var standardStartingFEN string = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

var discriminatorRE = regexp.MustCompile(`^[a-h]?[1-8]?$`)
var squareRE = regexp.MustCompile(`^[a-h][1-8]$`)

type PlayerColor int
const (
	White PlayerColor = iota
	Black 
)

type PieceType string
const (
	King PieceType = "K"
	Queen PieceType = "Q"
	Rook PieceType = "R"
	Bishop PieceType = "B"
	Knight PieceType = "N"
	Pawn PieceType = ""
)
var IsValidPieceType = map[PieceType]struct{} {
	King: {},
	Queen: {},
	Rook: {},
	Bishop: {},
	Knight: {},
	Pawn: {},
}
type piece struct {
	PieceType PieceType
	PlayerColor PlayerColor
	Square string
}

type GameState struct {
	PlayerTurn PlayerColor
	Pieces map[string]piece
}

type Move struct {
	PieceType, PromoteTo PieceType
	Target, Discriminator string
	IsCheck, IsCheckmate,
	IsCapture, IsLongCastle,
	IsShortCastle bool
}

type GameResult string
const (
	WhiteWins GameResult = "1-0"
	BlackWins GameResult = "0-1"
	Draw      GameResult = "1/2-1/2"
	Unknown   GameResult = "*"
)

var IsValidGameResult = map[GameResult]struct{} {
	WhiteWins: {},
	BlackWins: {},
	Draw: {},
	Unknown: {},
}

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
	InitalFEN string `json:"initialFen"`
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

func GameFromPGN(data []byte) (*Game, error) {
	if bytes.ContainsAny(data, "\t") {
		return &Game{}, errors.New("PGN format must not contain any tabs")
	}
	scanner := bufio.NewScanner(bytes.NewReader(data))
	scanner.Split(tokenizerPGN)

	valuesMap := make(map[string]string)
	game := Game{}
	for scanner.Scan() {
		// Move strings should not be lower case, but PGN keys should be for matching
		cleanKey := strings.TrimSpace(scanner.Text())
		key := strings.ToLower(cleanKey)
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
		} else {
			// If none of the keys match, then it is assumed this is the move string.
			// This assumption seems likely to be wrong in edge cases. As of now, I don't
			// have an assertive way of identifying a move string.
			// Assuming move strings has caused this many bugs: 1
			valuesMap["moves"] = cleanKey
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
		case "fen":
			game.InitalFEN = strings.TrimSpace(val)
		case "moves":
			moveNumberRE, err := regexp.Compile(`^\d+\.$`)
			if err != nil {
				log.Printf("Bad regexp: %v\n", err)
				break
			}
			var gameMoves string
			for token := range strings.SplitSeq(val, " ") {
				if moveNumberRE.MatchString(token) {
					continue
				}
				if _, ok := IsValidGameResult[GameResult(token)]; ok {
					continue
				}
				gameMoves = fmt.Sprintf("%s %s", gameMoves, token)
			}
			game.Moves = strings.TrimSpace(gameMoves)
		default:
		}
	}
	return &game, nil
}

func tokenizerMoveString(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return
	}
	if len(data) >= len(longCastle) && string(data[:len(longCastle)]) == longCastle {
		token = slices.Concat(token, data[:len(longCastle)])
		advance = len(longCastle)
		return
	}

	if len(data) >= len(shortCastle) && string(data[:len(shortCastle)]) == shortCastle {
		token = data[:len(shortCastle)]
		advance = len(shortCastle)
		return
	}

	if slices.Contains(pieceTokens, string(data[0])) {
		advance++
		token = append(token, data[0])
		return
	}

	if len(data) == 1 {
		token = data
		advance++
		switch t := string(data); t {
		case check, mate:
		default:
			err = errors.New("Unknown symbol")
		}
		return
	}

	if slices.Contains(rankTokens, string(data[0])) {
		token = append(token, data[0])
		advance++
		return
	}

	if slices.Contains(fileTokens, string(data[0])) {
		token = append(token, data[0])
		advance++
		if slices.Contains(rankTokens, string(data[1])) {
			token = append(token, data[1])
			advance++
		}
		return
	}

	if string(data[0]) == capture {
		token = append(token, data[0])
		advance++
		return
	}
	if string(data[0]) == promote {
		token = append(token, data[0])
		advance++
		return
	}
	return
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
[TimeControl "%s"]%s

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
	moveString = fmt.Sprintf("%s %s", moveString, result)

	var fen string
	if game.InitalFEN == "" {
		fen = fmt.Sprintf("\n[FEN \"%s\"]", standardStartingFEN)
	} else {
		fen = fmt.Sprintf("\n[FEN \"%s\"]", game.InitalFEN)
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
				fen,
				moveString,
	)
	return gamePGN, nil
}

func NewGameState(fen string) (gs *GameState, err error) {
	if strings.TrimSpace(fen) == standardStartingFEN {
		gs = initalGameState()
		return
	}

	if fen == "" {
		err = errors.New("Empty FEN string")
		return
	}

	fenFields := strings.Split(fen, " ")
	if len(fenFields) != 6 {
		err = errors.New("Invalid FEN string")
		return
	}

	fenRanks := strings.Split(fenFields[0], "/")
	if len(fenRanks) != 8 {
		err = errors.New("Invalid move string")
		return
	}

	gs = &GameState{}

	if fenFields[1] == "w" {
		gs.PlayerTurn = White
	} else {
		gs.PlayerTurn = Black
	}

	squareTracker := 0
	for ch := range fen[0] {
		if unicode.IsNumber(rune(ch)) {
			num, _ := strconv.Atoi(string(ch))
			squareTracker = squareTracker + num
			continue
		}
		newPiece := piece{}
		if unicode.IsUpper(rune(ch)) {
			newPiece.PlayerColor = White
		} else {
			newPiece.PlayerColor = Black
		}
		switch strings.ToLower(string(ch)) {
		case "p":
			newPiece.PieceType = Pawn
		case "n":
			newPiece.PieceType = Knight
		case "b":
			newPiece.PieceType = Bishop
		case "r":
			newPiece.PieceType = Rook
		case "q":
			newPiece.PieceType = Queen
		case "k":
			newPiece.PieceType = King
		}
		newPiece.Square = fenBoardOrder[squareTracker]
		gs.Pieces[newPiece.Square] = newPiece
		squareTracker++
	}

	return
}

func ParseMoveString(ms string) (move *Move, err error) {
	scanner := bufio.NewScanner(strings.NewReader(strings.TrimSpace(ms)))
	scanner.Split(tokenizerMoveString)
	var tokens []string
	for scanner.Scan() {
		tokens = append(tokens, scanner.Text())
	}
	if len(tokens) < 1 {
		err = fmt.Errorf("No tokens found in string: %v\n", ms)
		return
	}

	move = &Move{}
	switch tokens[0] {
	case "K":
		move.PieceType = King
	case "Q":
		move.PieceType = Queen
	case "B":
		move.PieceType = Bishop
	case "R":
		move.PieceType = Rook
	case "N":
		move.PieceType = Knight
	default:
		move.PieceType = Pawn
	}
	var i int
	if move.PieceType == Pawn {
		i = 0
	} else {
		i = 1
	}

	for ; i < len(tokens); i++ {
		t := tokens[i]
		if squareRE.MatchString(t) || discriminatorRE.MatchString(t) {
			if move.Target != "" {
				move.Discriminator = move.Target
			}
			move.Target = t
			continue
		}
		switch t {
		case longCastle:
			move.Target = t
			move.IsLongCastle = true
		case shortCastle:
			move.Target = t
			move.IsShortCastle = true
		case capture:
			move.IsCapture = true
		case check:
			move.IsCheck = true
		case mate:
			move.IsCheckmate = true
		case promote:
			if move.PromoteTo != Pawn {
				err = errors.New("Promotion indicated on non-pawn piece")
				return
			}
			i++
			if i >= len(tokens) {
				// Ensure we don't have a slice out of range error
				err = errors.New("Promotion indicated with no piece type provided")
				return
			}
			promotePiece := tokens[i]
			if !slices.Contains(pieceTokens, promotePiece) {
				err = errors.New("Promotion indicated with no piece type provided")
				return
			}
			if promotePiece == "" {
				err = errors.New("Cannot promote pawn to pawn")
				return
			}
			move.PromoteTo = PieceType(promotePiece)
		default:
			err = errors.New("Unknown token")
		}
	}
	return
}

