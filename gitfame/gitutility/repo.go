package gitutility

import (
	"os/exec"
	"strings"
)

type RepoHandler string

func (repo *RepoHandler) GetFilesInRevision(revision string) ([]string, error) {
	cmd := exec.Command("git", "ls-tree", "-r", "--name-only", revision)
	cmd.Dir = string(*repo)
	fileNames, err := cmd.Output()

	if err != nil {
		return nil, err
	}

	if len(fileNames) == 0 {
		return nil, nil
	}

	output := strings.Split(strings.TrimSpace(string(fileNames)), "\n")

	return output, nil
}

func (repo *RepoHandler) GetBlameInfo(fileName string, revision string) (string, error) {

	cmd := exec.Command("git", "blame", fileName, "--porcelain", revision)
	cmd.Dir = string(*repo)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	//	lines := strings.Split(string(output), "\n")

	return string(output), nil
}

func (repo *RepoHandler) GetGitLog(fileName string, revision string) (string, error) {

	cmd := exec.Command("git", "log", revision, "-1", "--pretty=format:%H %an", "--", fileName)
	cmd.Dir = string(*repo)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	//lines := strings.Split(string(output), "\n")

	return string(output), nil
}
