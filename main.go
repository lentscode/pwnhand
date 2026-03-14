package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "fatal: need ELF file path")
		os.Exit(1)
	}

	elfPath := os.Args[1]
	cmd := exec.Command("objdump", "-d", "-M", "intel", elfPath)
	var objdumpRes strings.Builder
	cmd.Stdout = &objdumpRes

	err := cmd.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	linesStr := strings.Split(objdumpRes.String(), "\n")
	nLines := len(linesStr)
	lines := make([]line, nLines)
	for i := range lines {
		lines[i] = line{
			idx:     i,
			content: linesStr[i],
		}
	}
	m := &model{
		lines: lines,
	}

	for ; nLines > 0; nLines /= 10 {
		m.lineNumberFigures++
	}

	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
