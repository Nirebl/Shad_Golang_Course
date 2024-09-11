package printer

import (
	"io"

	"gitlab.com/slon/shad-go/gitfame/gitstats"
)

type Printer interface {
	Print(w io.Writer, u []gitstats.UserInfo) error
}
