package printer

import (
	"fmt"
)

func Select(format string) (Printer, error) {
	switch format {
	case "tabular":
		var p tabularprinter
		return &p, nil
	case "csv":
		var p csvprinter
		return &p, nil
	case "json":
		var p jsonprinter
		return &p, nil
	case "json-lines":
		var p jsonlineprinter
		return &p, nil
	default:
		return nil, fmt.Errorf("No formatter for: " + format)
	}

}
