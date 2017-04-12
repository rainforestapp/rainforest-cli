package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/rainforestapp/rainforest-cli/rainforest"
	"github.com/urfave/cli"
)

const (
	// Version of the app in SemVer
	version = "2.0.1"
	// This is the default spec folder for RFML tests
	defaultSpecFolder = "./spec/rainforest"
)

var (
	// Build info to be set while building using:
	// go build -ldflags "-X main.build=build"
	build string

	// Release channel to be set while building using:
	// go build -ldflags "-X main.releaseChannel=channel"
	releaseChannel string

	// Rainforest API client
	api *rainforest.Client

	// default output for printing resource tables
	tablesOut io.Writer = os.Stdout

	// Run status polling interval
	runStatusPollInterval = time.Second * 5

	// Batch size (number of rows) for tabular var upload
	tabularBatchSize = 50
	// Concurrent connections when uploading CSV rows
	tabularConcurrency = 1
	// Maximum concurrent connections with Rainforest server
	rfmlDownloadConcurrency = 4
	// Concurrent connections when uploading RFML files
	rfmlUploadConcurrency = 4
)

// cliContext is an interface providing context of running application
// i.e. command line options and flags. One of the types that provides the interface is
// cli.Context, the other is fakeCLIContext which is used for testing.
type cliContext interface {
	String(flag string) (val string)
	StringSlice(flag string) (vals []string)
	Bool(flag string) (val bool)
	Int(flag string) (val int)

	Args() (args cli.Args)
}

// Create custom writer which will use timestamps
type logWriter struct{}

func (l *logWriter) Write(p []byte) (int, error) {
	log.Printf("%s", p)
	return len(p), nil
}

// main is an entry point of the app. It sets up the new cli app, and defines the API.
func main() {
	updateFinishedChan := make(chan struct{})
	app := cli.NewApp()
	app.Usage = "Rainforest QA CLI - https://www.rainforestqa.com/"
	app.Version = version
	if releaseChannel != "" {
		app.Version = fmt.Sprintf("%v - %v channel", app.Version, releaseChannel)
	}
	if build != "" {
		app.Version = fmt.Sprintf("%v - build: %v", app.Version, build)
	}

	// Use our custom writer to print our errors with timestamps
	cli.ErrWriter = &logWriter{}

	// Before running any of the commands we init the API Client & update
	app.Before = func(c *cli.Context) error {
		go autoUpdate(c, updateFinishedChan)

		api = rainforest.NewClient(c.String("token"))

		// Set the User-Agent that will be used for api calls
		if build != "" {
			api.UserAgent = "rainforest-cli/" + version + " build: " + build
		} else {
			api.UserAgent = "rainforest-cli/" + version
		}

		return nil
	}

	// Wait for the update to finish if it's still going on
	app.After = func(c *cli.Context) error {
		<-updateFinishedChan
		return nil
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "token, t",
			Usage:  "API token. You can find it at https://app.rainforestqa.com/settings/integrations",
			EnvVar: "RAINFOREST_API_TOKEN",
		},
		cli.BoolFlag{
			Name:  "skip-update",
			Usage: "Used to disable auto-updating of the cli",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "Run your tests on Rainforest",
			Action:  startRun,
			Description: "Runs your tests on Rainforest platform. " +
				"You need to specify list of test IDs to run or use keyword 'all'. " +
				"Alternatively you can use one of the filtering options.",
			ArgsUsage: "[test IDs]",
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "tag",
					Usage: "filter tests by `TAG`. Can be used multiple times for filtering by multiple tags.",
				},
				cli.StringFlag{
					Name:  "site, site-id",
					Usage: "filter tests by a specific site. You can see a list of your `SITE-ID`s with sites command.",
				},
				cli.StringFlag{
					Name:  "folder",
					Usage: "filter tests by a specific folder. You can see a list of your `FOLDER-ID`s with folders command.",
				},
				cli.StringSliceFlag{
					Name: "browser, browsers",
					Usage: "specify the `BROWSER` you wish to run against. This overrides test level settings." +
						"Can be used multiple times to run against multiple browsers.",
				},
				cli.StringFlag{
					Name:  "environment-id",
					Usage: "run your tests using specified `ENVIRONMENT`. Otherwise it will use your default one.",
				},
				cli.StringFlag{
					Name:  "crowd",
					Value: "default",
					Usage: "run your tests using specified `CROWD`. Available choices are: default or on_premise_crowd. " +
						"Contact your CSM for more details.",
				},
				cli.StringFlag{
					Name: "conflict",
					Usage: "use the abort option to abort any runs in the same environment or " +
						"use the abort-all option to abort all runs in progress.",
				},
				cli.BoolFlag{
					Name: "bg, background",
					Usage: "run in the background. This option makes cli return after succesfully starting a run, " +
						"without waiting for the run results.",
				},
				cli.BoolFlag{
					Name: "fail-fast, ff",
					Usage: "fail the build as soon as the first failed result comes in. " +
						"If you don't pass this it will wait until 100% of the run is done. Use with --fg.",
				},
				cli.StringFlag{
					Name: "custom-url",
					Usage: "use a custom `URL` for this run. Example use case: an ad-hoc QA environment with Fourchette. " +
						"You will need to specify a site_id too for this to work.",
				},
				cli.BoolFlag{
					Name: "git-trigger",
					Usage: "only trigger a run when the last commit (for a git repo in the current working directory) " +
						"contains @rainforest and a list of one or more tags. rainforest-cli exits with 0 otherwise.",
				},
				cli.StringFlag{
					Name:  "description",
					Usage: "add arbitrary `DESCRIPTION` to the run.",
				},
				cli.StringFlag{
					Name:  "junit-file",
					Usage: "Create a JUnit XML report `FILE` with the specified name. Must be run in foreground mode.",
				},
				cli.StringFlag{
					Name:  "import-variable-name",
					Usage: "`NAME` of the tabular variable to be created or updated.",
				},
				cli.StringFlag{
					Name:  "import-variable-csv-file",
					Usage: "`PATH` to the CSV file to be uploaded.",
				},
				cli.BoolFlag{
					Name:  "overwrite-variable",
					Usage: "If the flag is set, named variable will be updated.",
				},
				cli.BoolFlag{
					Name:  "single-use",
					Usage: "This option marks uploaded variable as single-use",
				},
				cli.StringFlag{
					Name:  "wait, reattach",
					Usage: "monitor existing run with `RUN_ID` instead of starting a new one.",
				},
			},
		},
		{
			Name:      "new",
			Usage:     "Create a new RFML test",
			ArgsUsage: "[name]",
			Description: "Create new Rainforest test in RFML format (Rainforest Markup Language). " +
				"You may also specify a custom test title or file name.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "test-folder",
					Value:  defaultSpecFolder,
					Usage:  "`PATH` at which to create new test.",
					EnvVar: "RAINFOREST_TEST_FOLDER",
				},
			},
			Action: newRFMLTest,
		},
		{
			Name:      "validate",
			Usage:     "Validate your RFML tests",
			ArgsUsage: "[path to RFML file]",
			Description: "Validate your test for syntax. " +
				"If no filepath is given it validates all RFML tests and performs additional checks for RFML ID validity and more. " +
				"If API token is set it'll validate your tests against server data as well.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "test-folder",
					Value:  "./spec/rainforest/",
					Usage:  "`PATH` where to look for a tests to validate.",
					EnvVar: "RAINFOREST_TEST_FOLDER",
				},
			},
			Action: validateRFML,
		},
		{
			Name:      "upload",
			Usage:     "Upload your RFML tests",
			ArgsUsage: "[path to RFML file]",
			Description: "Uploads specified test to Rainforest. " +
				"If no filepath is given it uploads all RFML tests.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "test-folder",
					Value:  "./spec/rainforest/",
					Usage:  "`PATH` where to look for a tests to upload.",
					EnvVar: "RAINFOREST_TEST_FOLDER",
				},
				cli.BoolFlag{
					Name:  "synchronous-upload",
					Usage: "uploads your test in a synchronous manner i.e. not using concurrency.",
				},
			},
			Action: uploadRFML,
		},
		{
			Name:        "rm",
			Usage:       "Remove an RFML test locally and remotely",
			ArgsUsage:   "[path to RFML file]",
			Description: "Remove RFML file and remove test from Rainforest test suite.",
			Action:      deleteRFML,
		},
		{
			Name: "download",
			// Left for legacy reason, should nuke?
			Aliases:   []string{"export"},
			Usage:     "Download your remote Rainforest tests to RFML",
			ArgsUsage: "[test IDs]",
			Description: "Download your remote tests from Rainforest to RFML. " +
				"You may specify list of test IDs or download all tests by default. " +
				"Alternatively you can use one of the filtering options.",
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "tag",
					Usage: "filter tests by `TAG`. Can be used multiple times for filtering by multiple tags.",
				},
				cli.IntFlag{
					Name:  "site, site-id",
					Usage: "filter tests by a specific site. You can see a list of your `SITE-ID`s with sites command.",
				},
				cli.IntFlag{
					Name:  "folder, folder-id",
					Usage: "filter tests by a specific folder. You can see a list of your `FOLDER-ID`s with folders command.",
				},
				cli.StringFlag{
					Name:   "test-folder",
					Value:  "./spec/rainforest/",
					Usage:  "`PATH` at which to save all the downloaded tests.",
					EnvVar: "RAINFOREST_TEST_FOLDER",
				},
				cli.BoolFlag{
					Name:  "embed-tests",
					Usage: "download your tests without extracting the steps of an embedded test.",
				},
			},
			Action: func(c *cli.Context) error {
				return downloadRFML(c, api)
			},
		},
		{
			Name:        "csv-upload",
			Usage:       "Create or update tabular var from CSV.",
			Description: "Upload a CSV file to create or update tabular variables.",
			ArgsUsage:   "[path to CSV file]",
			Flags: []cli.Flag{
				cli.StringFlag{
					// Alternative name left for legacy reason.
					Name:  "name, import-variable-name",
					Usage: "`NAME` of the tabular variable to be created or updated.",
				},
				cli.BoolFlag{
					Name:  "overwrite-variable, overwrite",
					Usage: "If the flag is set, named variable will be updated.",
				},
				cli.BoolFlag{
					Name:  "single-use",
					Usage: "This option marks uploaded variable as single-use",
				},
				// Left here for legacy reason, but imho we should move that to args
				cli.StringFlag{
					Name:  "csv-file, import-variable-csv-file",
					Usage: "DEPRECATED: `PATH` to the CSV file to be uploaded. Since v2 please provide the path as an argument.",
				},
			},
			Action: func(c *cli.Context) error {
				return csvUpload(c, api)
			},
		},
		{
			Name:  "report",
			Usage: "Create a report from your run results",
			Description: "Creates a report from your specified run." +
				"You can specify output file using options, otherwise report will be generated to STDOUT",
			ArgsUsage: "[run ID]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "junit-file",
					Usage: "`PATH` of file to which write a JUnit report for the specified run.",
				},
				cli.StringFlag{
					Name:  "run-id",
					Usage: "DEPRECATED: ID of a run for which to generate results. Since v2 please provide the run ID as an argument.",
				},
			},
			Action: createReport,
		},
		{
			Name:  "sites",
			Usage: "Lists available sites",
			Action: func(c *cli.Context) error {
				return printSites(api)
			},
		},
		{
			Name:  "folders",
			Usage: "Lists available folders",
			Action: func(c *cli.Context) error {
				return printFolders(api)
			},
		},
		{
			Name:  "browsers",
			Usage: "Lists available browsers",
			Action: func(c *cli.Context) error {
				return printBrowsers(api)
			},
		},
		{
			Name:      "update",
			Usage:     "Updates application to the latest version on specified release channel (stable/beta)",
			ArgsUsage: "[CHANNEL]",
			Action:    updateCmd,
		},
	}

	app.Run(shuffleFlags(os.Args))
}

// shuffleFlags moves global flags to the beginning of args array (where they are supposed to be),
// so they are picked up by the cli package, even though they are supplied as a command argument.
// might not be needed if upstream makes this change as well
func shuffleFlags(originalArgs []string) []string {
	globalOptions := []string{}
	rest := []string{}

	// We need to skip the filename as its arg[0] that's why iteration starts at 1
	// then filter out global flags and put them into separate array than the rest of arg
	for i := 1; i < len(originalArgs); i++ {
		option := originalArgs[i]
		if option == "--token" || option == "-t" {
			if i+1 < len(originalArgs) && originalArgs[i+1][:1] != "-" {
				globalOptions = append(globalOptions, originalArgs[i:i+2]...)
				i++
			} else {
				log.Fatalln("No token specified with --token flag")
			}
		} else if option == "--skip-update" {
			globalOptions = append(globalOptions, option)
		} else {
			rest = append(rest, option)
		}
	}

	shuffledFlags := []string{originalArgs[0]}
	shuffledFlags = append(shuffledFlags, globalOptions...)
	shuffledFlags = append(shuffledFlags, rest...)

	return shuffledFlags
}
