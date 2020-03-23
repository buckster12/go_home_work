package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inputString string) (string, error) {
	if inputString == "" {
		return "", nil
	}
	outputString := strings.Builder{}
	var lastStringElement rune
	var isPrevSlash bool

	for _, runeElement := range inputString {
		if isPrevSlash {
			// Asterisk task
			if unicode.IsDigit(runeElement) || string(runeElement) == `\` {
				lastStringElement = runeElement
			} else {
				return "", ErrInvalidString
			}
			outputString.WriteRune(lastStringElement)
			isPrevSlash = false
			continue
		}

		if string(runeElement) == `\` {
			// Go to the next rune
			isPrevSlash = true
			continue
		}

		// If we found a digit - repeat the last inserted element N-1 times
		if unicode.IsDigit(runeElement) {
			if lastStringElement == 0 {
				return "", ErrInvalidString
			}
			multiplier, _ := strconv.Atoi(string(runeElement))

			// -1 because we've already printed the first rune (on previous cycle, before digit)
			multiplier--
			outputString.WriteString(strings.Repeat(string(lastStringElement), multiplier))
			lastStringElement = 0
			continue
		}

		// Insert default rune without any multiplier or slashes
		lastStringElement = runeElement
		outputString.WriteRune(runeElement)
	}
	if isPrevSlash {
		return "", ErrInvalidString
	}
	return outputString.String(), nil
}
