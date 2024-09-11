//go:build !solution

package varfmt

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func Sprintf(format string, args ...interface{}) string {
	aStr := make([]string, 0, len(args))
	aSize := 0
	for _, arg := range args {
		s := fmt.Sprint(arg)
		aStr = append(aStr, s)
		aSize += len(s)
	}

	var sb strings.Builder
	sb.Grow(len(format) + aSize)

	pos := 0
	n := 0
	nLen := -1
	for len(format) > 0 {
		r, size := utf8.DecodeRuneInString(format)
		format = format[size:]
		switch true {
		case r == '{':
			nLen = 0
		case r == '}':
			write := n
			if nLen <= 0 {
				write = pos
			}
			sb.WriteString(aStr[write])
			n = 0
			nLen = -1
			pos++
		case nLen >= 0:
			n = n*10 + (int(r) - '0')
			nLen++
		default:
			sb.WriteRune(r)
		}
	}

	return sb.String()
}
