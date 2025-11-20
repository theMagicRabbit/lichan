package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

type StockfishProc struct {
	Cmd *exec.Cmd
	Stdin io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser
	Ready chan bool
	Moves string
}



func InitStockfish() (proc *StockfishProc, err error) {
	proc = &StockfishProc{
		Cmd: exec.CommandContext(context.Background(), "stockfish"),
		Ready: make(chan bool),
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
		command = "position startpos\n"
	} else {
		command = fmt.Sprintf("position fen %s\n", fen)
	}
	sp.Moves = fmt.Sprintf("%s moves", command)
	_, err = sp.Stdin.Write([]byte(command))
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
			fmt.Println("Stockfish best move:", tokens[1])
			sp.Ready <- true
		default:
			fmt.Println(tokens)
		}
	}
}
