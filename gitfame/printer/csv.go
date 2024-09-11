package printer

import (
	"encoding/csv"
	"fmt"
	"io"

	"gitlab.com/slon/shad-go/gitfame/gitstats"
)

type csvprinter struct{}

func (p *csvprinter) Print(w io.Writer, u []gitstats.UserInfo) error {
	writer := csv.NewWriter(w)
	header := []string{"Name", "Lines", "Commits", "Files"}
	if err := writer.Write(header); err != nil {
		return err
	}
	for _, info := range u {
		if info.Files == 0 {
			continue
		}
		row := []string{
			info.Name,
			fmt.Sprintf("%d", info.Lines),
			fmt.Sprintf("%d", info.Commits),
			fmt.Sprintf("%d", info.Files),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	writer.Flush()

	return nil
}
