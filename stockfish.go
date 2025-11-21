package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"
)

type StockfishProc struct {
	Cmd    *exec.Cmd
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser
	Info   struct {
		Mu    *sync.Mutex
		Value []string
	}
	Ready    chan bool
	Bestmove chan string
	Moves    string
}

var extendedMoveRe *regexp.Regexp = regexp.MustCompile(`([a-h][1-8]){2}[qrbn]?`)

func InitStockfish() (proc *StockfishProc, err error) {
	proc = &StockfishProc{
		Cmd:      exec.CommandContext(context.Background(), "stockfish"),
		Ready:    make(chan bool),
		Bestmove: make(chan string),
		Info: struct {
			Mu    *sync.Mutex
			Value []string
		}{
			Mu: &sync.Mutex{},
		},
	}

	stdin, err := proc.Cmd.StdinPipe()
	if err != nil {
		return
	}

	stdout, err := proc.Cmd.StdoutPipe()
	if err != nil {
		return
	}

	stderr, err := proc.Cmd.StderrPipe()
	if err != nil {
		return
	}
	proc.Stdin = stdin
	proc.Stdout = stdout
	proc.Stderr = stderr
	return
}

func (sp *StockfishProc) SetupGame(fen string) (err error) {
	var command string
	if fen == standardStartingFEN {
		command = "position startpos"
	} else {
		command = fmt.Sprintf("position fen %s", fen)
	}
	sp.Moves = fmt.Sprintf("%s moves", command)
	_, err = sp.Stdin.Write([]byte(command + "\n"))
	if err != nil {
		return
	}

	err = sp.IsReady()
	return
}

func (sp *StockfishProc) IsReady() (err error) {
	_, err = sp.Stdin.Write([]byte("isready\n"))
	return
}

func (sp *StockfishProc) SearchMove(move string) (err error) {
	sp.Moves = fmt.Sprintf("%s %s", sp.Moves, move)
	thinkTime := int(time.Minute / time.Millisecond)

	err = sp.IsReady()
	if err != nil {
		return
	}
	<-sp.Ready

	_, err = sp.Stdin.Write([]byte(sp.Moves + "\n"))
	if err != nil {
		return
	}

	err = sp.IsReady()
	if err != nil {
		return
	}
	<-sp.Ready

	command := fmt.Sprintf("go depth 245 movetime %d\n", thinkTime)
	_, err = sp.Stdin.Write([]byte(command))
	return
}

func (sp *StockfishProc) ProcessOutput() {
	sfScanner := bufio.NewScanner(sp.Stdout)
	for sfScanner.Scan() {
		text := sfScanner.Text()
		if text == "" {
			continue
		}

		tokens := strings.Split(text, " ")
		switch tokens[0] {
		case "Stockfish", "option", "id":
			continue
		case "uciok", "readyok":
			sp.Ready <- true
		case "bestmove":
			sp.Bestmove <- tokens[1]
		case "info":
			if slices.Contains(tokens, "currmovenumber") {
				break
			}
			sp.Info.Mu.Lock()
			sp.Info.Value = tokens
			sp.Info.Mu.Unlock()
		default:
			fmt.Println(tokens)
		}
	}
}

func GetPVMoves(PV []string) (moves []string, err error) {
	for i, token := range PV {
		if extendedMoveRe.MatchString(token) {
			moves = PV[i:]
			break
		}
	}
	return
}
