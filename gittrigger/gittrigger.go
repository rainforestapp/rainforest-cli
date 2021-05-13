package gitTrigger

import (
	"bytes"
	"os/exec"
	"regexp"
	"strings"
)

const gitTriggerString = "@rainforest"

type gitTrigger struct {
	Trigger    string
	LastCommit string
}

func NewGitTrigger() (gitTrigger, error) {
	newGit := gitTrigger{Trigger: gitTriggerString}
	err := newGit.GetLatestCommit()
	if err != nil {
		return gitTrigger{}, err
	}
	return newGit, nil
}

func (g *gitTrigger) GetLatestCommit() error {
	var out bytes.Buffer
	cmd := exec.Command("git", "log", "-1", "--pretty=%B")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	g.LastCommit = strings.TrimSpace(out.String())
	return nil
}

func (g *gitTrigger) GetRemote() (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("git", "remote")
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		return "", err
	}
	remote_name := strings.TrimSpace(out.String())

	out.Reset()
	cmd.Stdout = &out

	cmd = exec.Command("git", "remote", "get-url", "--all", remote_name)
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(strings.Split(strings.ReplaceAll(out.String(), "\r\n", "\n"), "\n")[0]), nil
}

func (g gitTrigger) CheckTrigger() bool {
	shouldTrigger := strings.Contains(g.LastCommit, g.Trigger)
	return shouldTrigger
}

func (g gitTrigger) GetTags() []string {
	tagRegex, _ := regexp.Compile("#([\\w_-]+)")
	foundTags := tagRegex.FindAllString(g.LastCommit, -1)
	strippedTags := make([]string, len(foundTags))
	for i, tag := range foundTags {
		strippedTags[i] = strings.Replace(tag, "#", "", 1)
	}
	return strippedTags
}
