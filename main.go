package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "fatal: need ELF file path")
		os.Exit(1)
	}

	var content string
	yamlData, err := loadFile("pwnhand.yaml")
	if err == nil {
		content = yamlData.Disas
	} else {
		elfPath := os.Args[1]
		cmd := exec.Command("objdump", "-d", "-M", "intel", elfPath)
		var objdumpRes strings.Builder
		cmd.Stdout = &objdumpRes

		err := cmd.Run()
		if err != nil {
			fmt.Fprintln(os.Stderr, "fatal: objdump command failed")
			os.Exit(1)
		}

		content = objdumpRes.String()
	}

	linesStr := strings.Split(content, "\n")
	nLines := len(linesStr)
	lines := make([]line, nLines)
	for i := range lines {
		lines[i] = line{
			idx:     i,
			content: linesStr[i],
		}
	}
	if yamlData != nil {
		for _, c := range yamlData.Comments {
			fields := strings.Split(c, ":")
			lineIdx, err := strconv.Atoi(fields[0])
			if err != nil {
				continue
			}

			switch fields[1] {
			case "eol":
				lines[lineIdx].eolComm = fields[2]
			case "plate":
				lines[lineIdx].plateComm = fields[2]
			}

		}
	}

	m := &model{
		content: content,
		lines:   lines,
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
