package main

import (
	"log"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/urfave/cli"
)

func update(silent bool) error {
	if !silent {
		log.Printf("Checking for update")
	}

	v := semver.MustParse(version)
	updater, err := selfupdate.NewUpdater(selfupdate.Config{BinaryName: "rainforest"})
	if err != nil {
		return err
	}
	latest, err := updater.UpdateSelf(v, "rainforestapp/rainforest-cli")
	if err != nil {
		return err
	}
	if !silent && latest.Version.Equals(v) {
		log.Println("No update available, already at the latest version!")
	} else if latest.Version.NE(v) {
		log.Printf("Updated to new version: %s!\n", latest.Version)
	}

	return nil
}

func updateCmd(c cliContext) error {
	err := update(false)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

func autoUpdate(c cliContext, updateFinishedChan chan<- struct{}) {
	if !c.Bool("skip-update") {
		update(true)
	}
	updateFinishedChan <- struct{}{}
}
