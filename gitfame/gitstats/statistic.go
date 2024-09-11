package gitstats

type Statistic struct {
	Commits int `json:"commits"`
	Files   int `json:"files"`
	Lines   int `json:"lines"`
}

type UserInfo struct {
	Name string `json:"name"`
	Statistic
}

func (currentStat *Statistic) Increase(increaseStat *Statistic) {
	currentStat.Commits += increaseStat.Commits
	currentStat.Files += increaseStat.Files
	currentStat.Lines += increaseStat.Lines
}
