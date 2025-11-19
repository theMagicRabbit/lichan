package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os/exec"
)

type StockfishProc struct {
	Cmd *exec.Cmd
	Stdin *io.PipeWriter
	Stdout *io.PipeReader
	Stderr *io.PipeReader
}



func InitStockfish() (proc *StockfishProc, err error) {
	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()
	stderrReader, stderrWriter := io.Pipe()
	proc = &StockfishProc{
		Cmd: exec.CommandContext(context.Background(), "stockfish"),
		Stdin: stdinWriter,
		Stdout: stdoutReader,
		Stderr: stderrReader,
	}
	proc.Cmd.Stdin = stdinReader
	proc.Cmd.Stdout = stdoutWriter
	proc.Cmd.Stderr = stderrWriter
	return
}

func (sp *StockfishProc) ProcessOutput() {
	sfScanner := bufio.NewScanner(sp.Stdout)
	for {
		if sfScanner.Scan() {
			fmt.Println("Output:", sfScanner.Text())
		}
	}
}
