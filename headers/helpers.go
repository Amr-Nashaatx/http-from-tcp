package headers

import (
	"slices"
)

func isDigit(c rune) bool {
	digs := []rune{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	return slices.Contains(digs, c)
}

func isLowerCaseChar(c rune) bool {
	for s := 'a'; s <= 'z'; s++ {
		if c == s {
			return true
		}
	}
	return false
}
func isUpperCaseChar(c rune) bool {
	for s := 'A'; s <= 'Z'; s++ {
		if c == s {
			return true
		}
	}
	return false
}
func IsValidToken(token string) bool {
	for _, c := range token {
		isChar := isLowerCaseChar(c) || isUpperCaseChar(c)
		isDigit := isDigit(c)
		isSpecial := isSpecialChar(c)

		if !isChar && !isDigit && !isSpecial {
			return false
		}
	}

	return true
}

func isSpecialChar(c rune) bool {
	spcChars := []rune{'!', '#', '$', '%', '&', ',', '*', '+', '-', '.', '^', '_', '`', '|', '~'}

	return slices.Contains(spcChars, c)
}
