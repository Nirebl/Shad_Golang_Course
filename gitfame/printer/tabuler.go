package printer

import (
	"fmt"
	"io"
	"text/tabwriter"

	"gitlab.com/slon/shad-go/gitfame/gitstats"
)

type tabularprinter struct{}

func (p *tabularprinter) Print(w io.Writer, u []gitstats.UserInfo) error {
	myWriter := tabwriter.NewWriter(w, 1, 1, 1, ' ', 0)
	_, err := fmt.Fprint(myWriter, "Name\tLines\tCommits\tFiles\n")
	if err != nil {
		return err
	}
	for _, info := range u {
		_, err := fmt.Fprintf(myWriter, "%s\t%d\t%d\t%d\n", info.Name, info.Lines, info.Commits, info.Files)
		if err != nil {
			return err
		}
	}
	if err := myWriter.Flush(); err != nil {
		return err
	}
	return nil
}
