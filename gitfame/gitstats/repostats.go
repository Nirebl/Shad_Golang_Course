package gitstats

import (
	"path"
	"path/filepath"

	"gitlab.com/slon/shad-go/gitfame/gitutility"
	"gitlab.com/slon/shad-go/gitfame/language"
)

type RepoRequestData struct {
	RepoHandler  gitutility.RepoHandler
	Revision     string
	Extensions   []string
	Languages    []string
	Exclude      []string
	RestrictTo   []string
	UseCommitter bool
}

type reopStat struct {
	commits map[string]struct{}
	files   int
	lines   int
}

func (repoInfo *RepoRequestData) isExtensionValid(fileName string) bool {
	if len(repoInfo.Extensions) == 0 {
		return true
	}

	fileExtension := filepath.Ext(fileName)

	for _, ext := range repoInfo.Extensions {
		if ext == fileExtension {
			return true
		}
	}

	return false
}

func (repoInfo *RepoRequestData) isLanguageValid(fileName string) bool {
	if len(repoInfo.Languages) == 0 {
		return true
	}

	fileExtension := filepath.Ext(fileName)

	for _, lang := range repoInfo.Languages {
		langExts := language.GetLanguageExtensions(lang)

		for _, ext := range langExts {
			if ext == fileExtension {
				return true
			}
		}
	}

	return false
}

func (repoInfo *RepoRequestData) isFileNameExclude(fileName string) bool {
	if len(repoInfo.Exclude) == 0 {
		return true
	}

	for _, pattern := range repoInfo.Exclude {
		match, err := path.Match(pattern, fileName)
		if err != nil {
			return false
		}
		if match {
			return false
		}
	}

	return true
}

func (repoInfo *RepoRequestData) isRestrictTo(fileName string) bool {
	if len(repoInfo.RestrictTo) == 0 {
		return true
	}

	for _, pattern := range repoInfo.RestrictTo {
		match, err := path.Match(pattern, fileName)
		if err != nil {
			return false
		}
		if match {
			return true
		}
	}

	return false
}

func (repoInfo *RepoRequestData) isFileNameValid(fileName string) bool {
	if !repoInfo.isFileNameExclude(fileName) {
		return false
	}

	if !repoInfo.isExtensionValid(fileName) {
		return false
	}

	if !repoInfo.isLanguageValid(fileName) {
		return false
	}
	if !repoInfo.isRestrictTo(fileName) {
		return false
	}

	return true
}

func (repoInfo *RepoRequestData) GetRepoStats() ([]UserInfo, error) {

	files, err := repoInfo.RepoHandler.GetFilesInRevision(repoInfo.Revision)
	if err != nil {
		return nil, err
	}

	userStats := make(map[string]reopStat)

	for _, fileName := range files {
		if !repoInfo.isFileNameValid(fileName) {
			continue
		}

		fileStats, err := calculate(fileName, repoInfo)
		if err != nil {
			return nil, err
		}

		for name, fileStat := range fileStats {
			userStat, ok := userStats[name]
			if !ok {
				userStat.commits = make(map[string]struct{})
			}
			userStat.files += 1
			userStat.lines += fileStat.lines

			for hash := range fileStat.commits {
				userStat.commits[hash] = struct{}{}
			}

			userStats[name] = userStat
		}
	}

	var userInfo []UserInfo

	for name, stat := range userStats {
		var info UserInfo
		info.Name = name
		info.Files = stat.files
		info.Lines = stat.lines
		info.Commits = len(stat.commits)

		userInfo = append(userInfo, info)
	}

	return userInfo, nil
}
