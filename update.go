package main

import (
	"log"

	"github.com/equinox-io/equinox"
	"github.com/urfave/cli"
)

const equinoxAppID = "app_carcVJmQBRm"

// publicKey is a ECDSA key used to sign the cli binaries
var publicKey = []byte(`
-----BEGIN ECDSA PUBLIC KEY-----
MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAEoTjUnZL6bSXdEV9LKD+0zNekogPp8lwf
8s/auLQXwCHFii7fH8bLwSW+4a9+eF8bWo8FYk4pfSo3WT5DUqGl9qHcnv22MMCK
eiFH+GffIMk09RFqkcMh5rIPu3ykm5V8
-----END ECDSA PUBLIC KEY-----
`)

func update(channel string, silent bool) error {
	opts := equinox.Options{}
	if channel == "" && releaseChannel != "" {
		opts.Channel = releaseChannel
	} else if channel != "" {
		opts.Channel = channel
	} else {
		// fallback to stable
		opts.Channel = "stable"
	}

	if err := opts.SetPublicKeyPEM(publicKey); err != nil {
		return err
	}

	// check for the update
	if !silent {
		log.Printf("Checking for update on %v channel.", opts.Channel)
	}
	resp, err := equinox.Check(equinoxAppID, opts)
	switch {
	case err == equinox.NotAvailableErr:
		if !silent {
			log.Println("No update available, already at the latest version!")
		}
		return nil
	case err != nil:
		return err
	}

	// fetch the update and apply it
	if !silent {
		log.Print("Found a cli update, applying it.")
	}
	err = resp.Apply()
	if err != nil {
		return err
	}

	if !silent {
		log.Printf("Updated to new version: %s!\n", resp.ReleaseVersion)
	}
	return nil
}

func updateCmd(c cliContext) error {
	channel := c.Args().First()
	if !(channel == "beta" || channel == "stable" || channel == "") {
		return cli.NewExitError("Invalid release channel - use 'stable' or 'beta'", 1)
	}
	err := update(channel, false)
	if err != nil {
		return cli.NewExitError(err.Error(), 1)
	}
	return nil
}

func autoUpdate(c cliContext, updateFinishedChan chan<- struct{}) {
	if !c.Bool("skip-update") {
		update("", true)
	}
	updateFinishedChan <- struct{}{}
}
