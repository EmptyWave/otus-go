package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var builder strings.Builder
	var i int

	shieldedBackSlash := '\\'

	explodedStr := []rune(str)
	explodedStrLength := len(explodedStr)

	for i < explodedStrLength {
		curVal := explodedStr[i]

		if unicode.IsDigit(curVal) {
			return "", ErrInvalidString
		}

		if curVal == shieldedBackSlash && i+1 < explodedStrLength {
			i++
			curVal = explodedStr[i]

			if curVal != shieldedBackSlash && !unicode.IsDigit(explodedStr[i]) {
				return "", ErrInvalidString
			}
		}

		if i+1 < explodedStrLength && unicode.IsDigit(explodedStr[i+1]) {
			i++
			builder.WriteString(strings.Repeat(string(curVal), int(explodedStr[i]-'0')))
		} else {
			builder.WriteRune(curVal)
		}

		i++
	}

	return builder.String(), nil
}
