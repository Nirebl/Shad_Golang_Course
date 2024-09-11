//go:build !solution

package speller

import (
	"fmt"
	"strings"
)

var (
	ones  = []string{"", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	teens = []string{"", "eleven", "twelve", "thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen"}
	tens  = []string{"", "ten", "twenty", "thirty", "forty", "fifty", "sixty", "seventy", "eighty", "ninety"}
)

func Spell(n int64) string {
	if n == 0 {
		return "zero"
	}
	if n < 0 {
		return "minus " + Spell(-n)
	}

	var spell string

	if n/1000000000 > 0 {
		spell += fmt.Sprintf("%s billion ", Spell(n/1000000000))
		n %= 1000000000
	}
	if n/1000000 > 0 {
		spell += fmt.Sprintf("%s million ", Spell(n/1000000))
		n %= 1000000
	}
	if n/1000 > 0 {
		spell += fmt.Sprintf("%s thousand ", Spell(n/1000))
		n %= 1000
	}
	if n/100 > 0 {
		spell += fmt.Sprintf("%s hundred ", ones[n/100])
		n %= 100
	}
	if n >= 20 {
		if n >= 21 && n%10 != 0 {
			spell += tens[n/10] + "-"
		} else {
			spell += tens[n/10]

		}

		n %= 10
	}
	if n >= 11 && n <= 19 {
		spell += teens[n-10]
	} else {
		if n == 10 {
			spell += "ten"
		} else {
			spell += ones[n]
		}
	}

	return strings.TrimSpace(spell)
}
