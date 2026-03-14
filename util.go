package main

import (
	"regexp"
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
