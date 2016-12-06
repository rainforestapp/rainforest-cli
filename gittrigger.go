package main

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

func newGitTrigger() (gitTrigger, error) {
	newGit := gitTrigger{Trigger: gitTriggerString}
	err := newGit.getLatestCommit()
	if err != nil {
		return gitTrigger{}, err
	}
	return newGit, nil
}

func (g *gitTrigger) getLatestCommit() error {
	cmd := exec.Command("git", "log", "-1", "--pretty=%B")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return err
	}
	g.LastCommit = strings.TrimSpace(out.String())
	return nil
}

func (g gitTrigger) checkTrigger() bool {
	shouldTrigger := strings.Contains(g.LastCommit, g.Trigger)
	return shouldTrigger
}

func (g gitTrigger) getTags() []string {
	tagRegex, _ := regexp.Compile("#([\\w_-]+)")
	foundTags := tagRegex.FindAllString(g.LastCommit, -1)
	strippedTags := make([]string, len(foundTags))
	for i, tag := range foundTags {
		strippedTags[i] = strings.Replace(tag, "#", "", 1)
	}
	return strippedTags
}
