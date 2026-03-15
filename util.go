package main

const (
	ansiEscape       = "\x1b"
	ansiColorReset   = ansiEscape + "[0m"
	ansiColorBlackFg = ansiEscape + "[0;30m"
	ansiColorWhiteBg = ansiEscape + "[0;100m"
	ansiColorBold    = ansiEscape + "[1m"
)

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

func bold(str string) string {
	return ansiColorBold + str + ansiColorReset
}
