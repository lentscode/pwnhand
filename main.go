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
			idx: i,
			content: linesStr[i],
		}
	}
	m := model{
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

type model struct {
	lines    []line
	rows     int
	columns  int
	y        int

	lineNumberFigures int
	availableRows int
}

type comment struct {
	content  string
	line     int
	commType commType
}

type commType int

const (
	commEOL commType = iota
)

type line struct {
	idx int
	content string
	comments []comment
}

func (m model) renderLine(l line) []string {
	lines := make([]string, 0)
	contentLength := len(l.content)	
	if m.availableRows == 0 {
		return lines
	}
	
	i := 0
	newLine := strings.Builder{}
	format := fmt.Sprintf("%%%dd | ", m.lineNumberFigures)
	fmt.Fprintf(&newLine, format, l.idx)
	if contentLength == 0 {
		lines = append(lines, newLine.String())
		return lines
	}

	nSpaces := m.rows - m.availableRows
	for i < contentLength {
		if i > 0 {
			newLine.WriteString(strings.Repeat(" ", nSpaces))
		}
		end := min(i + m.availableRows, contentLength)
		newLine.WriteString(l.content[i:end])
		lines = append(lines,  newLine.String())
		newLine.Reset()
		i = end
	}

	return lines
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "k":
			m.y = max(m.y-1, 0)
		case "j":
			m.y = min(m.y+1, len(m.lines))
		case "ctrl+u":
			m.y = max(m.y-m.columns/2, 0)
		case "ctrl+d":
			m.y = min(m.y+m.columns/2, len(m.lines))
		}
	case tea.WindowSizeMsg:
		m.columns = msg.Height
		m.rows = msg.Width
		m.availableRows = m.rows - m.lineNumberFigures - 3
	case tea.MouseMotionMsg:
		m.y = min(max(m.y+msg.Y, 0), len(m.lines))
	}
	return m, nil
}

func (m model) View() tea.View {
	lines := make([]string, 0)
	for _, l := range m.lines {
		lines = append(lines, m.renderLine(l)...)
	}
	lines = lines[m.y:min(m.y+m.columns, len(lines))]
	v := tea.NewView(strings.Join(lines, "\n") + "\n")
	return v
}
