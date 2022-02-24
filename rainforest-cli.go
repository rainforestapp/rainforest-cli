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
	version = "2.28.0"
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

	// CircleCI Orb version
	orbVersion string
	// GitHub Action version
	ghActionVersion string

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
// cli.Context, the other is fakeContext which is used for testing.
type cliContext interface {
	String(flag string) (val string)
	GlobalString(flag string) (val string)
	StringSlice(flag string) (vals []string)
	GlobalStringSlice(flag string) (vals []string)
	Bool(flag string) (val bool)
	GlobalBool(flag string) (val bool)
	Int(flag string) (val int)
	GlobalInt(flag string) (val int)
	Uint(flag string) (val uint)
	GlobalUint(flag string) (val uint)

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
	app.Name = "Rainforest CLI"
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

		api = rainforest.NewClient(c.String("token"), c.Bool("debug"))

		// Set the User-Agent that will be used for api calls
		api.UserAgent = "rainforest-cli/" + version
		orbVersion = os.Getenv("ORB_VERSION")
		ghActionVersion = os.Getenv("GH_ACTION_VERSION")
		if orbVersion != "" {
			api.UserAgent += " rainforest-orb/" + orbVersion
		} else if ghActionVersion != "" {
			api.UserAgent += " rainforest-gh-action/" + ghActionVersion
		}
		if build != "" {
			api.UserAgent += " build: " + build
		}
		api.SendTelemetry = !c.Bool("disable-telemetry")

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
		cli.BoolFlag{
			Name:  "disable-telemetry",
			Usage: "Stops the cli sharing information about which CI system you may be using, and where you host your git repo (i.e. your git remote). Rainforest uses this to better integrate with CI tooling, and code hosting companies, it is not sold or shared. Disabling this may affect your Rainforest experience.",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Output http request header information for debug purposes",
		},
	}
	app.OnUsageError = func(c *cli.Context, err error, isSubcommand bool) error {
		return cli.NewExitError("Unknown argument", 1)
	}

	app.Commands = []cli.Command{
		{
			Name:         "run",
			Aliases:      []string{"r"},
			Usage:        "Run your tests on Rainforest",
			OnUsageError: onCommandUsageErrorHandler("run"),
			Action:       func(c *cli.Context) error {
				return startRun(c)
			},
			Description: "Runs your tests on Rainforest platform. " +
				"You need to specify list of test IDs to run or use keyword 'all'. " +
				"Alternatively you can use one of the filtering options.",
			ArgsUsage: "[test IDs]",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "f, files",
					Usage: "Run local tests specified by `FILES or FOLDERS`",
				},
				cli.StringSliceFlag{
					Name:  "tag",
					Usage: "Filter tests by `TAG`. Can be used multiple times for filtering by multiple tags.",
				},
				cli.StringSliceFlag{
					Name:  "exclude",
					Usage: "Don't execute test specified by `FILE`. Can be used multiple times for specifying multiple files.",
				},
				cli.StringSliceFlag{
					Name:  "force-execute",
					Usage: "Execute test specified by `FILE` even if execute: false is specified. Can be used multiple times for specifying multiple files.",
				},
				cli.StringFlag{
					Name:  "site, site-id",
					Usage: "Filter tests by a specific site. You can see a list of your `SITE-ID`s with the sites command.",
				},
				cli.StringFlag{
					Name:  "folder, folder-id, filter, filter-id",
					Usage: "Filter tests by a specific folder. You can see a list of your `FOLDER-ID`s with the folders command.",
				},
				cli.IntFlag{
					Name:  "feature, feature-id",
					Usage: "Filter tests by a specific feature. You can see a list of your `FEATURE-ID`s with the features command.",
				},
				cli.IntFlag{
					Name:  "run-group, run-group-id",
					Usage: "Start a run using a run group. You can see a list of your `RUN-GROUP-ID`s with the run-groups command. This option cannot be used in conjunction with other filtering options.",
				},
				cli.StringSliceFlag{
					Name: "browser, browsers",
					Usage: "Specify the `BROWSER` you wish to run against. This overrides test level settings." +
						"Can be used multiple times to run against multiple browsers.",
				},
				cli.StringFlag{
					Name:  "environment-id",
					Usage: "Run your tests using specified `ENVIRONMENT`. Otherwise it will use your default one.",
				},
				cli.StringFlag{
					Name: "crowd",
					Usage: "Run your tests using specified `CROWD`. Available choices are: default, automation, automation_and_crowd " +
						"or on_premise_crowd. Contact your CSM for more details.",
				},
				cli.StringFlag{
					Name: "conflict",
					Usage: "Use the abort option to abort any runs in the same environment or " +
						"use the abort-all option to abort all runs in progress.",
				},
				cli.BoolFlag{
					Name: "bg, background",
					Usage: "Run in the background. This option makes cli return after successfully starting a run, " +
						"without waiting for the run results.",
				},
				cli.BoolFlag{
					Name: "fail-fast, ff",
					Usage: "Fail the build as soon as the first failed result comes in. " +
						"If you don't pass this it will wait until 100% of the run is done. Use with --fg.",
				},
				cli.StringFlag{
					Name: "custom-url",
					Usage: "Specify the URL for the run to use when testing against an ephemeral environment. " +
						"This will create a new temporary environment for the run.",
				},
				cli.BoolFlag{
					Name: "git-trigger",
					Usage: "Only trigger a run when the last commit (for a git repo in the current working directory) " +
						"contains @rainforest and a list of one or more tags. rainforest-cli exits with 0 otherwise.",
				},
				cli.StringFlag{
					Name:  "description",
					Usage: "Add arbitrary `DESCRIPTION` to the run.",
				},
				cli.StringFlag{
					Name: "release",
					Usage: "Adds a `RELEASE` ID that is associated with this run. You can use any string, but commonly used " +
						"IDs are commit SHAs, build IDs, branch names, etc.",
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
					Usage: "Monitor existing run with `RUN_ID` instead of starting a new one.",
				},
				cli.UintFlag{
					Name:  "max-reruns",
					Usage: "Rerun `max-reruns` times before reporting failure.",
				},
				cli.UintFlag{
					Name:  "automation-max-retries",
					Usage: "Try to pass a test `automation-max-retries` times within the same run before reporting run failure.",
				},
				cli.StringFlag{
					Name:  "save-run-id",
					Usage: "Save the created run's ID to `FILE`",
				},
			},
		},
		{
			Name:         "rerun",
			Aliases:      []string{"rr"},
			Usage:        "Rerun failed tests from a previous run",
			OnUsageError: onCommandUsageErrorHandler("rerun"),
			Action:       func(c *cli.Context) error {
				return rerunRun(c)
			},
			Description: "Reruns the failed tests from a previous run on Rainforest platform. " +
				"Parameters such as 'environment', 'crowd', 'release', etc. are copied from the previous run.",
			ArgsUsage: "[run ID]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name: "conflict",
					Usage: "Use the abort option to abort any runs in the same environment or " +
						"use the abort-all option to abort all runs in progress.",
				},
				cli.BoolFlag{
					Name: "bg, background",
					Usage: "Run in the background. This option makes cli return after successfully starting a run, " +
						"without waiting for the run results.",
				},
				cli.BoolFlag{
					Name: "fail-fast, ff",
					Usage: "Fail the build as soon as the first failed result comes in. " +
						"If you don't pass this it will wait until 100% of the run is done. Use with --fg.",
				},
				cli.StringFlag{
					Name:  "junit-file",
					Usage: "Create a JUnit XML report `FILE` with the specified name. Must be run in foreground mode.",
				},
				cli.UintFlag{
					Name:  "max-reruns",
					Usage: "Rerun `max-reruns` times before reporting failure.",
				},
				cli.UintFlag{
					Name:  "rerun-attempt",
					Usage: "Which rerun attempt this is.",
				},
				cli.StringFlag{
					Name:  "save-run-id",
					Usage: "Save the created run's ID to `FILE`",
					Value: ".rainforest_run_id",
				},
			},
		},
		{
			Name:         "new",
			Usage:        "Create a new RFML test",
			OnUsageError: onCommandUsageErrorHandler("new"),
			ArgsUsage:    "[name]",
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
			Action: func(c *cli.Context) error {
				return newRFMLTest(c)
			},
		},
		{
			Name:         "validate",
			Usage:        "Validate your RFML tests",
			OnUsageError: onCommandUsageErrorHandler("validate"),
			ArgsUsage:    "[path to RFML file]",
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
			Action: func(c *cli.Context) error {
				return validateRFML(c, api)
			},
		},
		{
			Name:         "upload",
			Usage:        "Upload your tests",
			OnUsageError: onCommandUsageErrorHandler("upload"),
			ArgsUsage:    "[path to file]",
			Description: "Uploads specified test to Rainforest. " +
				"If no filepath is given it uploads all tests.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:   "test-folder",
					Value:  "./spec/rainforest/",
					Usage:  "`PATH` where to look for a tests to upload.",
					EnvVar: "RAINFOREST_TEST_FOLDER",
				},
				cli.BoolFlag{
					Name:  "synchronous-upload",
					Usage: "Uploads your test in a synchronous manner i.e. not using concurrency.",
				},
			},
			Action: func(c *cli.Context) error {
				return uploadTests(c, api)
			},
		},
		{
			Name:         "rm",
			Usage:        "Remove an RFML test locally and remotely",
			OnUsageError: onCommandUsageErrorHandler("rm"),
			ArgsUsage:    "[path to RFML file]",
			Description:  "Remove RFML file and remove test from Rainforest test suite.",
			Action:       func(c *cli.Context) error {
				return deleteRFML(c)
			},
		},
		{
			Name: "download",
			// Left for legacy reason, should nuke?
			Aliases:      []string{"export"},
			Usage:        "Download your Rainforest tests",
			OnUsageError: onCommandUsageErrorHandler("download"),
			ArgsUsage:    "[test IDs]",
			Description: "Download your Rainforest tests. " +
				"You may specify list of test IDs or download all tests by default. " +
				"Alternatively you can use one of the filtering options.",
			Flags: []cli.Flag{
				cli.StringSliceFlag{
					Name:  "tag",
					Usage: "Filter tests by `TAG`. Can be used multiple times for filtering by multiple tags.",
				},
				cli.IntFlag{
					Name:  "site, site-id",
					Usage: "Filter tests by a specific site. You can see a list of your `SITE-ID`s with the sites command.",
				},
				cli.IntFlag{
					Name:  "folder, folder-id, filter, filter-id",
					Usage: "Filter tests by a specific folder. You can see a list of your `FOLDER-ID`s with the folders command.",
				},
				cli.IntFlag{
					Name:  "feature, feature-id",
					Usage: "Filter tests by a specific feature. You can see a list of your `FEATURE-ID`s with the features command.",
				},
				cli.IntFlag{
					Name:  "run-group, run-group-id",
					Usage: "Filter tests by a specific run group. You can see a list of your `RUN-GROUP-ID`s with the run-groups command.",
				},
				cli.StringFlag{
					Name:   "test-folder",
					Value:  "./spec/rainforest/",
					Usage:  "`PATH` at which to save all the downloaded tests.",
					EnvVar: "RAINFOREST_TEST_FOLDER",
				},
				cli.BoolFlag{
					Name:  "flatten-steps",
					Usage: "Download your tests with steps extracted from embedded tests.",
				},
			},
			Action: func(c *cli.Context) error {
				return downloadTests(c, api)
			},
		},
		{
			Name:         "csv-upload",
			Usage:        "Create or update tabular var from CSV.",
			OnUsageError: onCommandUsageErrorHandler("csv-upload"),
			Description:  "Upload a CSV file to create or update tabular variables.",
			ArgsUsage:    "[path to CSV file]",
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
			Name:         "mobile-upload",
			Usage:        "Upload your mobile app to Rainforest.",
			OnUsageError: onCommandUsageErrorHandler("mobile-upload"),
			Description:  "Upload a mobile app file to Rainforest.",
			ArgsUsage:    "[path to mobile app file]",
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "site-id",
					Usage: "The site-id of the app you are uploading. You can see a list of your `SITE-ID`s with the sites command.",
				},
				cli.IntFlag{
					Name:  "environment-id",
					Usage: "The environment-id of the app you are uploading. You can see a list of your `ENVIRONMENT-ID`s with the environment command.",
				},
				cli.StringFlag{
					Name:  "app-slot",
					Usage: "An optional flag for specifying the app slot (1-100) of your app, if your site-environment contains multiple apps. Default is 1.",
				},
			},
			Action: func(c *cli.Context) error {
				return mobileAppUpload(c, api)
			},
		},
		{
			Name:         "report",
			Usage:        "Create a report from your run results",
			OnUsageError: onCommandUsageErrorHandler("report"),
			Description: "Creates a report from your specified run." +
				"You can specify output file using options, otherwise report will be generated to STDOUT",
			ArgsUsage: "[run ID]",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "junit-file",
					Usage: "`PATH` of file to which write a JUnit report for the specified run.",
				},
			},
			Action: func(c *cli.Context) error {
				return writeJunit(c, api, 0)
			},
		},
		{
			Name:         "sites",
			Usage:        "Lists available sites",
			OnUsageError: onCommandUsageErrorHandler("sites"),
			Action: func(c *cli.Context) error {
				return printSites(api)
			},
		},
		{
			Name:         "environments",
			Usage:        "Lists available environments",
			OnUsageError: onCommandUsageErrorHandler("environments"),
			Action: func(c *cli.Context) error {
				return printEnvironments(api)
			},
		},
		{
			Name:         "folders",
			Usage:        "Lists available folders",
			OnUsageError: onCommandUsageErrorHandler("folders"),
			Action: func(c *cli.Context) error {
				return printFolders(api)
			},
		},
		{
			Name:         "filters",
			Usage:        "Lists available saved filters",
			OnUsageError: onCommandUsageErrorHandler("filters"),
			Action: func(c *cli.Context) error {
				return printFolders(api)
			},
		},
		{
			Name:         "browsers",
			Usage:        "Lists available browsers",
			OnUsageError: onCommandUsageErrorHandler("browsers"),
			Action: func(c *cli.Context) error {
				return printBrowsers(api)
			},
		},
		{
			Name:         "features",
			Usage:        "Lists available features",
			OnUsageError: onCommandUsageErrorHandler("features"),
			Action: func(c *cli.Context) error {
				return printFeatures(api)
			},
		},
		{
			Name:         "run-groups",
			Usage:        "Lists available run groups",
			OnUsageError: onCommandUsageErrorHandler("run-groups"),
			Action: func(c *cli.Context) error {
				return printRunGroups(api)
			},
		},
		{
			Name:         "update",
			Usage:        "Updates application to the latest version",
			OnUsageError: onCommandUsageErrorHandler("update"),
			Action:       func(c *cli.Context) error {
				return updateCmd(c)
			},
		},
	}

	app.Run(shuffleFlags(os.Args))
}

// shuffleFlags moves global flags to the beginning of args array (where they
// are supposed to be), so they are picked up by the cli package, even though
// they are supplied as a command argument.  We also do a bit of hacking to
// allow "multiple arguments" to -f.
func shuffleFlags(originalArgs []string) []string {
	globalOptions := []string{}
	fnameArgs := []string{}
	rest := []string{}

	// We need to skip the filename as its arg[0] that's why iteration starts at 1
	// then filter out global flags and put them into separate array than the rest of arg
	for i := 1; i < len(originalArgs); i++ {
		option := originalArgs[i]
		if option == "--token" || option == "-t" {
			if i+1 < len(originalArgs) && len(originalArgs[i+1]) > 0 && originalArgs[i+1][:1] != "-" {
				globalOptions = append(globalOptions, originalArgs[i:i+2]...)
				i++
			} else {
				log.Fatalln("No token specified with --token flag")
			}

		} else if option == "-f" || option == "--files" {
			rest = append(rest, option)
			i++
			for i < len(originalArgs) && originalArgs[i][0] != '-' {
				fnameArgs = append(fnameArgs, originalArgs[i])
				i++
			}
			i--
		} else if option == "--disable-telemetry" {
			globalOptions = append(globalOptions, option)
		} else if option == "--skip-update" {
			globalOptions = append(globalOptions, option)
		} else if option == "--debug" {
			globalOptions = append(globalOptions, option)
		} else {
			rest = append(rest, option)
		}
	}

	shuffledFlags := []string{originalArgs[0]}
	shuffledFlags = append(shuffledFlags, globalOptions...)
	shuffledFlags = append(shuffledFlags, rest...)
	shuffledFlags = append(shuffledFlags, fnameArgs...)

	return shuffledFlags
}

func onCommandUsageErrorHandler(command string) func(*cli.Context, error, bool) error {
	return func(c *cli.Context, err error, isSubcommand bool) error {
		fmt.Printf("Incorrect usage: %s\n\n", err.Error())

		cli.ShowCommandHelp(c, command)
		os.Exit(1)
		return nil
	}
}
