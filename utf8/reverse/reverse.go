//go:build !solution

package reverse

import (
	"strings"
	"unicode/utf8"
)

func Reverse(input string) string {
	var answer strings.Builder
	answer.Grow(len(input))

	for len(input) > 0 {
		currentRune, currentLen := utf8.DecodeLastRuneInString(input)
		answer.WriteRune(currentRune)
		input = input[:len(input)-currentLen]
	}
	return answer.String()
}
