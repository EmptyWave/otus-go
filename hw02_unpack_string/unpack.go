package hw02unpackstring

import (
	"errors"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var shield bool
	var prevVal rune
	var builder, tempBuilder strings.Builder

	for i, curVal := range str {
		if unicode.IsDigit(curVal) && i == 0 {
			return "", ErrInvalidString
		}

		if unicode.IsLetter(curVal) {
			if shield {
				shield = false

				curVal = prevVal + curVal
			} else if unicode.IsLetter(prevVal) {
				builder.WriteString(tempBuilder.String())
				tempBuilder.Reset()
			}

			tempBuilder.WriteRune(curVal)
			prevVal = curVal
		} else if unicode.IsDigit(curVal) {
			if shield {
				shield = false

				if unicode.IsLetter(prevVal) {
					builder.WriteString(tempBuilder.String())
					tempBuilder.Reset()
				}

				tempBuilder.WriteRune(curVal)
				prevVal = curVal
			} else {
				if tempBuilder.Len() == 0 {
					return "", ErrInvalidString
				}

				builder.WriteString(strings.Repeat(tempBuilder.String(), int(curVal-'0')))
				tempBuilder.Reset()
			}
		} else if curVal == '\\' {
			if tempBuilder.Len() > 0 {
				builder.WriteString(tempBuilder.String())
				tempBuilder.Reset()
			}

			if prevVal == '\\' {
				shield = false
				tempBuilder.WriteRune(curVal)
				prevVal = rune(0)
			} else {
				shield = true
				prevVal = curVal
			}
		}
	}

	if tempBuilder.Len() > 0 {
		builder.WriteString(tempBuilder.String())
	}

	return builder.String(), nil
}
