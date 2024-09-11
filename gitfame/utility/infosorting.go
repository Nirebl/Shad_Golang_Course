package utility

import (
	"os"
	"sort"

	"gitlab.com/slon/shad-go/gitfame/gitstats"
)

func SortUserInfos(infos []gitstats.UserInfo, orderBy string) []gitstats.UserInfo {

	sorted := infos

	switch orderBy {
	case "lines":
		sort.SliceStable(sorted, func(i, j int) bool {
			if sorted[i].Lines == sorted[j].Lines {
				if sorted[i].Commits == sorted[j].Commits {
					if sorted[i].Files == sorted[j].Files {
						return sorted[i].Name < sorted[j].Name
					} else {
						return sorted[i].Files > sorted[j].Files
					}
				} else {
					return sorted[i].Commits > sorted[j].Commits
				}
			} else {
				return sorted[i].Lines > sorted[j].Lines
			}
		})
	case "commits":
		sort.SliceStable(sorted, func(i, j int) bool {
			if sorted[i].Commits == sorted[j].Commits {
				if sorted[i].Lines == sorted[j].Lines {
					if sorted[i].Files == sorted[j].Files {
						return sorted[i].Name < sorted[j].Name
					} else {
						return sorted[i].Files > sorted[j].Files
					}
				} else {
					return sorted[i].Lines > sorted[j].Lines
				}
			} else {
				return sorted[i].Commits > sorted[j].Commits
			}
		})
	case "files":
		sort.SliceStable(sorted, func(i, j int) bool {
			if sorted[i].Files == sorted[j].Files {
				if sorted[i].Lines == sorted[j].Lines {
					if sorted[i].Commits == sorted[j].Commits {
						return sorted[i].Name < sorted[j].Name
					} else {
						return sorted[i].Commits > sorted[j].Commits
					}
				} else {
					return sorted[i].Lines > sorted[j].Lines
				}
			} else {
				return sorted[i].Files > sorted[j].Files
			}
		})
	default:
		os.Exit(1)
	}

	return sorted
}
