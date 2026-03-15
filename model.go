package main

import (
	"fmt"
	"strconv"
	"strings"

	tea "charm.land/bubbletea/v2"
)

type model struct {
	content     string
	lines       []line
	actualLines []string
	rows        int
	columns     int
	y           int
	currentLine int

	lineNumberFigures int
	availableColumns  int

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
	idx     int
	content string
	lineComments
}

func (m *model) Init() tea.Cmd {
	return nil
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c":
			m.saveState()
			return m, tea.Quit
		case "k":
			if m.currentMode == normalMode {
				m.y = max(m.y-1, 0)
			} else {
				m.writeComment("k")
			}
		case "j":
			if m.currentMode == normalMode {
				m.y = min(m.y+1, len(m.actualLines)-1)
			} else {
				m.writeComment("j")
			}
		case "ctrl+u":
			if m.currentMode == normalMode {
				m.y = max(m.y-m.rows/2, 0)
			}
		case "ctrl+d":
			if m.currentMode == normalMode {
				m.y = min(m.y+m.rows/2, len(m.actualLines)-1)
			}
		case "c":
			if m.currentMode == normalMode {
				m.currentMode = eolCommentMode
			} else {
				m.writeComment("c")
			}
		case "esc":
			m.currentMode = normalMode
		case "backspace":
			if m.currentMode != normalMode {
				m.writeComment("backspace")
			}
		//TODO: better handle the enter, for now confirms the comment
		case "enter":
			m.currentMode = normalMode
		default:
			m.writeComment(msg.String())
		}

	case tea.WindowSizeMsg:
		m.columns = msg.Width
		m.rows = msg.Height - 1
		m.availableColumns = m.columns - m.lineNumberFigures - 3

	case tea.MouseMotionMsg:
		m.y = min(max(m.y+msg.Y, 0), len(m.actualLines))
	}

	lines := make([]string, 0)
	for _, l := range m.lines {
		lines = append(lines, m.renderLine(l)...)
	}
	m.actualLines = lines

	if len(m.actualLines) == 0 {
		return m, nil
	}
	m.setCurrentLine()
	return m, nil
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

func (m *model) View() tea.View {
	lines := m.actualLines[m.y:min(m.y+m.rows, len(m.actualLines))]
	var statusBar string
	if m.currentMode == normalMode {
		statusBar = "NORMAL"
	} else {
		statusBar = "COMMENT"
	}
	newLinesToAdd := max(0, m.rows-len(m.actualLines)+m.y) + 1
	v := tea.NewView(strings.Join(lines, "\n") + strings.Repeat("\n", newLinesToAdd) + statusBar)
	return v
}

func (m *model) renderLine(l line) []string {
	lines := make([]string, 0)
	if m.availableColumns == 0 {
		return lines
	}
	if l.lineComments.eolComm != "" || m.currentLine == l.idx && m.currentMode != normalMode {
		l.content += strings.Repeat(" ", 8) + "# " + l.lineComments.eolComm
	}

	contentLength := len(l.content)
	i := 0
	newLine := strings.Builder{}
	format := fmt.Sprintf("%%%dd | ", m.lineNumberFigures)
	fmt.Fprintf(&newLine, format, l.idx)
	if contentLength == 0 {
		lines = append(lines, newLine.String())
		return lines
	}

	nSpaces := m.columns - m.availableColumns
	for i < contentLength {
		if i > 0 {
			newLine.WriteString(strings.Repeat(" ", nSpaces))
		}
		end := min(i+m.availableColumns, contentLength)
		newLine.WriteString(l.content[i:end])
		lines = append(lines, newLine.String())
		newLine.Reset()
		i = end
	}

	return lines
}

func (m *model) setCurrentLine() {
	if lineStart := regex.FindString(m.actualLines[m.y]); lineStart != "" {
		lineStart = strings.TrimLeft(lineStart, " ")
		lineNumber, err := strconv.Atoi(lineStart)
		if err == nil {
			m.currentLine = lineNumber
		}
	}
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
			commentEntry := fmt.Sprintf("%d:eol:%s", l.idx, l.plateComm)
			yamlData.Comments = append(yamlData.Comments, commentEntry)
		}
	}

	m.err = saveFile("pwnhand.yaml", yamlData)
}
