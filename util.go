package main

import (
	"regexp"
)

const (
	ansiEscape       = "\x1b"
	ansiColorReset   = ansiEscape + "[0m"
	ansiColorBlackFg = ansiEscape + "[0;30m"
	ansiColorWhiteBg = ansiEscape + "[0;100m"
)

var regex = regexp.MustCompile(`^[ ]*\d+`)

func deleteCharacterFromString(str *string) {
	if len(*str) == 0 {
		return
	}

	*str = (*str)[:len(*str)-1]
}

func decodeSpecialCharacters(str string) string {
	switch str {
	case "enter":
		return "\n"
	case "space":
		return " "
	default:
		return str
	}
}

func blackOnWhite(str string) string {
	return ansiColorBlackFg + ansiColorWhiteBg + str + ansiColorReset
}
