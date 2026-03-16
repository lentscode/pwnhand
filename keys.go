package main

import tea "charm.land/bubbletea/v2"

func (m *model) handleKeyPress(msg tea.KeyPressMsg) tea.Cmd {
	switch msg.String() {
	case "ctrl+c":
		m.saveState()
		return tea.Quit
	case "k":
		m.handlek()
	case "j":
		m.handlej()
	case "ctrl+u":
		m.handleCtrlu()
	case "ctrl+d":
		m.handleCtrld()
	case "c":
		m.handlec()
	case "C":
		m.handleC()
	case "esc":
		m.currentMode = normalMode
	case "backspace":
		if m.currentMode != normalMode {
			m.writeComment("backspace")
		}
	//TODO: better handle the enter, for now confirms the comment
	case "enter":
		m.handleEnter()
	default:
		m.writeComment(msg.String())
	}

	return nil
}

func (m *model) handlek() {
	if m.currentMode == normalMode {
		m.currentLine = max(0, m.currentLine-1)
		currentLineSize := m.lines[m.currentLine].size + m.lines[m.currentLine].plateSize
		m.y = max(0, m.y-currentLineSize)
		if m.y < m.renderStart {
			m.renderEnd -= currentLineSize
			m.renderStart -= currentLineSize
		}
	} else {
		m.writeComment("k")
	}
}

func (m *model) handlej() {
	if m.currentMode == normalMode {
		currentLineSize := m.lines[m.currentLine].size
		if m.currentLine < len(m.lines)-1 {
			currentLineSize += m.lines[m.currentLine+1].plateSize
		}
		m.y = min(len(m.actualLines)-1, m.y+currentLineSize)
		m.currentLine = min(len(m.lines)-1, m.currentLine+1)
		if m.y >= m.renderEnd {
			m.renderEnd += currentLineSize
			m.renderStart += currentLineSize
		}
	} else {
		m.writeComment("j")
	}
}

func (m *model) handlec() {
	if m.currentMode == normalMode {
		m.currentMode = eolCommentMode
	} else {
		m.writeComment("c")
	}
}

func (m *model) handleC() {
	if m.currentMode == normalMode {
		m.currentMode = plateCommentMode
	} else {
		m.writeComment("C")
	}
}

func (m *model) handleEnter() {
	if m.currentMode == plateCommentMode {
		m.writeComment("enter")
	} else {
		m.currentMode = normalMode
	}
}

func (m *model) handleCtrlu() {
	return //TODO: fix this
	if m.currentMode == normalMode {
		offset := m.rows / 2
		m.currentLine = max(m.currentLine-offset, 0)

		renderOffset := min(max(0, m.renderStart), offset)
		m.renderEnd -= renderOffset
		m.renderStart -= renderOffset
	}
}

func (m *model) handleCtrld() {
	return //TODO: fix this
	if m.currentMode == normalMode {
		offset := m.rows / 2
		m.currentLine = min(m.currentLine+offset, len(m.actualLines)-1)

		renderOffset := min(max(0, len(m.actualLines)-m.renderEnd), offset)
		m.renderEnd += renderOffset
		m.renderStart += renderOffset
	}
}
