package printer

import (
	"encoding/json"
	"io"

	"gitlab.com/slon/shad-go/gitfame/gitstats"
)

type jsonlineprinter struct{}

func (p *jsonlineprinter) Print(w io.Writer, u []gitstats.UserInfo) error {
	for _, info := range u {
		data, err := json.Marshal(info)
		if err != nil {
			return err
		}
		if _, err := w.Write(data); err != nil {
			return err
		}
		if _, err := w.Write([]byte("\n")); err != nil {
			return err
		}
	}
	return nil
}
