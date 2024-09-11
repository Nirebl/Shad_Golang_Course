package printer

import (
	"encoding/json"
	"io"

	"gitlab.com/slon/shad-go/gitfame/gitstats"
)

type jsonprinter struct{}

func (p *jsonprinter) Print(w io.Writer, u []gitstats.UserInfo) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(u)
}
