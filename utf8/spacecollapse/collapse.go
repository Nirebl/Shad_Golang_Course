//go:build !solution

package spacecollapse

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

func CollapseSpaces(input string) string {
	var result strings.Builder
	result.Grow(len(input))

	prevSpace := false
	for len(input) > 0 {
		r, size := utf8.DecodeRuneInString(input)

		isSpace := unicode.IsSpace(r)
		if !isSpace {
			result.WriteRune(r)
		} else if !prevSpace {
			result.WriteRune(' ')
		}
		prevSpace = isSpace

		input = input[size:]
	}
	return result.String()
}
