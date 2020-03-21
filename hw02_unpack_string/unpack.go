package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(inputString string) (string, error) {
	//specials := "nr\\\""

	if inputString == "" {
		return "", nil
	}
	outputString := strings.Builder{}
	var lastStringElement string

	// transform string into slice
	inputSlice := []rune(inputString)

	for i := 0; i < len(inputSlice); i++ {
		stringElement := string(inputSlice[i])

		if stringElement == "\\" {
			// Go to the next rune
			i++

			// Asterisk task
			if unicode.IsDigit(inputSlice[i]) || string(inputSlice[i]) == `\` {
				lastStringElement = string(inputSlice[i])
			} else {
				return "", ErrInvalidString
			}
			outputString.WriteString(lastStringElement)
			continue
		}

		// If we found a digit - repeat the last inserted element N-1 times
		if unicode.IsDigit(inputSlice[i]) {
			if lastStringElement == "" {
				return "", ErrInvalidString
			}
			multiplier, _ := strconv.Atoi(stringElement)

			// -1 because we've already printed the first rune (on previous cycle, before digit)
			multiplier--
			outputString.WriteString(strings.Repeat(lastStringElement, multiplier))
			lastStringElement = ""
			continue
		}

		// Insert default rune without any multiplier or slashes
		lastStringElement = stringElement
		outputString.WriteString(stringElement)
	}
	return outputString.String(), nil
}
