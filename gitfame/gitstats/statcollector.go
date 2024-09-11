package gitstats

import (
	"strconv"
	"strings"
)

type fileStat struct {
	commits map[string]struct{}
	lines   int
}

func calculate(fileName string, info *RepoRequestData) (map[string]fileStat, error) {
	output, err := info.RepoHandler.GetBlameInfo(fileName, info.Revision)
	if err != nil {
		return nil, err
	}

	if len(output) == 0 {
		return getGitLog(fileName, info)
	}

	return parseBlameLines(strings.Split(output, "\n"), info.UseCommitter)
}

func getGitLog(fileName string, info *RepoRequestData) (map[string]fileStat, error) {
	output, err := info.RepoHandler.GetGitLog(fileName, info.Revision)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(output, "\n")

	result := make(map[string]fileStat)
	for _, line := range lines {
		currentFields := strings.Fields(line)

		prefix := currentFields[0] + " "
		name := strings.TrimPrefix(line, prefix)
		var fileStats fileStat
		fileStats.commits = make(map[string]struct{})
		fileStats.commits[currentFields[0]] = struct{}{}
		result[name] = fileStats
	}

	return result, nil
}

func parseBlameLines(lines []string, useCommitter bool) (map[string]fileStat, error) {
	userPattern := "author"
	if useCommitter {
		userPattern = "committer"
	}

	lineInCommit := make(map[string]int)
	userCommit := make(map[string]string)

	var currentHash string
	for _, line := range lines {
		currentFields := strings.Fields(line)
		//_, err := uuid.FromString(currentFields[0])
		if len(currentFields) == 0 {
			continue
		}

		if len(currentFields) == 4 { // && err == nil {
			currentHash = currentFields[0]
			count, _ := strconv.Atoi(currentFields[3])
			lineInCommit[currentHash] += count
			continue
		}
		if currentFields[0] == userPattern {
			_, ok := userCommit[currentHash]
			if !ok {
				prefix := userPattern + " "
				userCommit[currentHash] = strings.TrimPrefix(line, prefix)
			}
		}
	}

	result := make(map[string]fileStat)

	for hash, user := range userCommit {
		stat, ok := result[user]
		if !ok {
			stat.commits = make(map[string]struct{})
		}
		stat.commits[hash] = struct{}{}
		count := lineInCommit[hash]
		stat.lines += count

		result[user] = stat
	}

	return result, nil
}
