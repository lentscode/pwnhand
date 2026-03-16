package main

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	content     string
	lines       []line
	actualLines []string
	rows        int
	columns     int
	currentLine int

	lineNumberFigures int
	availableColumns  int

	y           int
	renderStart int
	renderEnd   int

	currentMode mode

	err error
}

type mode int

const (
	normalMode mode = iota
	eolCommentMode
	plateCommentMode
)

type lineComments struct {
	eolComm   string
	plateComm string
}

type commType int

const (
	eolComm commType = iota
	plateComm
)

type line struct {
	idx       int
	content   string
	size      int
	plateSize int
	lineComments
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if cmd := m.handleKeyPress(msg); cmd != nil {
			return m, cmd
		}
	case tea.WindowSizeMsg:
		m.columns = msg.Width
		m.rows = msg.Height - 1
		m.availableColumns = m.columns - m.lineNumberFigures - 3
		m.renderEnd = m.renderStart + msg.Height - 1

	case tea.MouseMotionMsg:
		m.currentLine = min(max(m.currentLine+msg.Y, 0), len(m.actualLines))
	}

	lines := make([]string, 0)
	for i := range m.lines {
		lines = append(lines, m.renderLine(&m.lines[i])...)
	}
	m.actualLines = lines

	if len(m.actualLines) == 0 {
		return m, nil
	}
	return m, nil
}

func (m model) View() tea.View {
	lines := m.actualLines[m.renderStart:m.renderEnd]
	var statusBar string
	if m.currentMode == normalMode {
		statusBar = "NORMAL"
	} else {
		statusBar = "COMMENT"
	}
	v := tea.NewView(strings.Join(lines, "\n") + "\n" + bold(statusBar))
	return v
}

func (m *model) renderLine(l *line) []string {
	lines := make([]string, 0)
	l.size = 1
	l.plateSize = 0
	content := l.content
	if m.availableColumns == 0 {
		return lines
	}
	if l.lineComments.eolComm != "" || m.currentLine == l.idx && m.currentMode == eolCommentMode {
		content += strings.Repeat(" ", 8) + "# " + l.lineComments.eolComm
	}

	nSpaces := m.columns - m.availableColumns
	if plateComm := l.lineComments.plateComm; plateComm != "" || m.currentLine == l.idx && m.currentMode == plateCommentMode {
		plateCommentLines := strings.SplitSeq(plateComm, "\n")
		for pl := range plateCommentLines {
			lines = append(lines, strings.Repeat(" ", nSpaces)+"# "+pl)
			l.plateSize++
		}

		if len(lines) == 0 {
			lines = append(lines, strings.Repeat(" ", nSpaces)+"# ")
		}
	}

	contentLength := len(content)
	newLine := strings.Builder{}
	format := fmt.Sprintf("%%%dd | ", m.lineNumberFigures)
	fmt.Fprintf(&newLine, format, l.idx)
	if contentLength == 0 {
		if l.idx == m.currentLine {
			lines = append(lines, blackOnWhite(newLine.String()))
		} else {
			lines = append(lines, newLine.String())
		}
		return lines
	}

	i := 0
	for i < contentLength {
		if i > 0 {
			newLine.WriteString(strings.Repeat(" ", nSpaces))
			l.size++
		}
		end := min(i+m.availableColumns, contentLength)
		newLine.WriteString(content[i:end])
		if l.idx == m.currentLine {
			lines = append(lines, blackOnWhite(newLine.String()))
		} else {
			lines = append(lines, newLine.String())
		}
		newLine.Reset()
		i = end
	}

	return lines
}

func (m *model) saveState() {
	yamlData := &yamlData{
		Disas:    m.content,
		Comments: make([]string, 0),
	}

	for _, l := range m.lines {
		if l.eolComm != "" {
			commentEntry := fmt.Sprintf("%d:eol:%s", l.idx, l.eolComm)
			yamlData.Comments = append(yamlData.Comments, commentEntry)
		}
		if l.plateComm != "" {
			commentEntry := fmt.Sprintf("%d:plate:%s", l.idx, l.plateComm)
			yamlData.Comments = append(yamlData.Comments, commentEntry)
		}
	}

	m.err = saveFile("pwnhand.yaml", yamlData)
}

func (m *model) writeComment(str string) {
	if str == "" {
		return
	}

	switch m.currentMode {
	case eolCommentMode:
		if str == "backspace" {
			deleteCharacterFromString(&m.lines[m.currentLine].eolComm)
		} else {
			m.lines[m.currentLine].eolComm += decodeSpecialCharacters(str)
		}
	case plateCommentMode:
		if str == "backspace" {
			deleteCharacterFromString(&m.lines[m.currentLine].plateComm)
		} else {
			m.lines[m.currentLine].plateComm += decodeSpecialCharacters(str)
		}
	}
}
