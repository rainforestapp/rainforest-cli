[![CircleCI](https://circleci.com/gh/rainforestapp/rainforest-cli.svg?style=shield)](https://circleci.com/gh/rainforestapp/rainforest-cli)

# Rainforest-cli

A command line interface to interact with [Rainforest QA](https://www.rainforestqa.com/).

This is the easiest way to integrate Rainforest with your deploy scripts or CI server. See [our documentation](https://help.rainforestqa.com/docs/rainforest-cli-for-continuous-integration) on the subject.

The CLI uses the Rainforest API which is documented at https://help.rainforestqa.com/reference.

### Wrappers

For users of CircleCI and GitHub Actions, we have platform-specific wrappers you can use: our [CircleCI Orb](https://circleci.com/developer/orbs/orb/rainforest-qa/rainforest) and our [GitHub Action](https://github.com/marketplace/actions/rainforest-qa-github-action).

## Installation

### Docker

```bash
$ docker pull gcr.io/rf-public-images/rainforest-cli
$ docker run gcr.io/rf-public-images/rainforest-cli --version
Rainforest CLI version 3.4.1 - build: docker
```

### Brew

If you are on OSX and use Brew, you can [install the CLI](https://github.com/rainforestapp/homebrew-public) with the command:

```bash
brew install rainforestapp/public/rainforest-cli
```

### Chocolatey

If you are on Windows and use Chocolatey, you can [install the CLI](https://community.chocolatey.org/packages/rainforest-cli/) with the command:

```bash
choco install rainforest-cli
```

### Binaries

Get the CLI binaries from our [Releases page](https://github.com/rainforestapp/rainforest-cli/releases).

The CLI will check for updates and automatically update itself on every use unless the global flag `--skip-update` is used.

You can download the latest (Linux) `rainforest-cli` binary with the following command (requires `curl`, `jq` and `tar`):

```bash
curl -sL $(curl -s https://api.github.com/repos/rainforestapp/rainforest-cli/releases/latest | jq -r '.assets[].browser_download_url | select(test("linux-amd64.tar.gz"))') | tar zxf - rainforest
```

### Migrating from our old CLI

The previous version of our CLI is deprecated and parts of it no longer work. To upgrade, follow our [migration guide](./migration.md).

## Basic Usage

To authenticate against Rainforest you'll need your API token which you can get from your [integrations settings page](https://app.rainforestqa.com/settings/integrations).

Pass the token into the CLI through the `RAINFOREST_API_TOKEN` environment variable or with the global flag `--token`.

CLI commands are formatted as follows:
```bash
rainforest [global flags] <command> [command-specific-flags] [arguments]
```

## Options

#### Running Tests

Run all tests in the foreground and report.

```bash
rainforest run all
```

Run all your tests in the background and exit the process immediately.

```bash
rainforest run all --bg
```

Run all tests with tag 'run-me' and cancel previous in-progress runs.

```bash
rainforest run --tag run-me --conflict cancel
```

Run all tests and generate a junit xml report.

```bash
rainforest run all --junit-file results.xml
```

Run individual tests in the foreground and report.

```bash
rainforest run <test_id1> <test_id2>
```

Run all tests on a branch.

```bash
rainforest run all --branch branch-name
```

Run a run group.

⚠️ This uses the configuration defined in the run group (environment, platforms, execution method, location). If you wish to run tests from a run group without using the run group's configuration, you will need to use the Rainforest API directly, passing a `run_group_id` parameter to [the `POST /runs` endpoint](https://help.rainforestqa.com/reference/post-runs). ⚠️

```bash
rainforest run --run-group <run_group_id>
```

#### Rerunning Failed Tests

```bash
rainforest rerun <failed_run_id>
```

The `failed_run_id` argument is optional. If none is passed in, the CLI will look for a run ID in the `RAINFOREST_RUN_ID` environment variable.

#### Creating and Managing Tests

Create new Rainforest test in RFML format (Rainforest Markup Language).

```bash
rainforest new
```

You may also specify a custom test title or file name.

```bash
rainforest new "My Awesome Title"
```

Validate your tests for syntax and correct RFML ids for embedded tests.
Use the `--token` options or `RAINFOREST_API_TOKEN` environment variable
to validate your tests against server data as well.

```bash
rainforest validate
```

Validate RFML syntax of a specified file.
This command just validates RFML syntax for more complex validation including checking
embedded tests id correctness and existence of potential circural dependiences in tests
use general command.

```bash
rainforest validate /path/to/test/file.rfml
```

Upload tests to Rainforest

```bash
rainforest upload
```

Upload a specific test to Rainforest

```bash
rainforest upload /path/to/test/file.rfml
```

Upload a specific test to Rainforest on a branch

```bash
rainforest upload --branch branch-name /path/to/test/file.rfml
```

Remove RFML file and remove test from Rainforest test suite.

```bash
rainforest rm /path/to/test/file.rfml
```

Download all tests from Rainforest

```bash
rainforest download
```

Download tests filtered by tags, site, and smart folder

```bash
rainforest download --tag foo --tag bar --site-id 123 --folder 456
```

Download specific tests based on their id on the Rainforest dashboard

```bash
rainforest download 33445 11232 1337
```

#### Running Local RFML Tests Only

If you want to run a local set of RFML files (for instance in a CI environment), use the `run -f` option:

```
rainforest run -f [FILES OR FOLDERS]
```

`run -f` accepts any number of files and folders as arguments (folders will be recursively checked for `*.rfml` files). All embedded tests must be included within `FILES OR FOLDERS`.

There is a specific metadata option in RFML files for `run -f`: `# execute: true|false`, which indicates whether a test should be run by default (defaults to `true`). Embedded tests that are not usually run directly should specify `# execute: false`.

The following options are specific to `run -f` or behave differently:

- `--tag TAG_NAME`: only run tests that are tagged with `TAG_NAME` (which can be a comma-separated list of tags). Note that this filters *within* local RFML files, not tests stored on Rainforest. Tests that are not tagged with `TAG_NAME` will not be executed but may be still be uploaded if they are embedded in another test.
- `--exclude FILE`: exclude the test in `FILE` from being run, even if `# execute: true` is specified.
- `--force-execute FILE`: execute the test in `FILE` even if `# execute: false` is specified.

Run-level setting options (`--platforms`, `--environment_id`, etc) behave the same for `run -f`. Other test filtering options (such as `--run-group`, `--site`, etc) cannot be used in conjunction with `run -f`.

#### Viewing Account Specific Information

See a list of all of your sites and their IDs
```bash
rainforest sites
```

See a list of all of your environments and their IDs
```bash
rainforest environments
```

See a list of all of your smart folders and their IDs
```bash
rainforest folders
```

See a list of all of your platforms and their IDs
```bash
rainforest platforms
```

See a list of all of your features and their IDs
```bash
rainforest features
```

See a list of all of your run groups and their IDs
```bash
rainforest run-groups
```

To fetch a junit xml report for a test run which has already completed
```bash
rainforest report <run-id> --junit-file rainforest.xml
```

#### Updating Tabular Variables

Upload a CSV to create a new tabular variables.
```bash
rainforest csv-upload --import-variable-name my_variable PATH/TO/CSV.csv
```

Upload a CSV to update an existing tabular variables.
```bash
rainforest csv-upload --import-variable-name my_variable --overwrite-variable PATH/TO/CSV.csv
```

#### Managing branches

Create a new branch.
```bash
rainforest branch new branch-name
```

Delete an existing branch.
```bash
rainforest branch delete branch-name
```

Merge an existing branch into the `main` branch.
```bash
rainforest branch merge branch-name
```

#### Uploading Mobile Apps

Upload a mobile app to Rainforest.
```bash
rainforest mobile-upload --site-id <site_id> --environment-id <environment_id> PATH/TO/mobile_app.ipa
```
- `--site-id SITE_ID` - The site ID of the app you are uploading. You can see a list of your site IDs with the `sites` command.
- `--environment-id ENVIRONMENT_ID` - The environment ID of the app you are uploading. You can see a list of your environment IDs with the `environments` command.
- `--app-slot SLOT` - An optional flag for specifying the app slot of your app, if your site-environment contains multiple apps. Valid values are from `1` to `100`, and the default value is `1`.


## Options

### Global

- `--token <your-rainforest-token>` - your API token if it's not set via the `RAINFOREST_API_TOKEN` environment variable
- `--skip-update` - Do not automatically check for CLI updates

### Writing Tests
Rainforest Tests written using RFML have the following format

```
#! [RFML ID]
# title: [TITLE]
# start_uri: [START_URI]
# tags: [TAGS]
# site_id: [SITE ID]
# platforms: [PLATFORM IDS]
# feature_id: [FEATURE_ID]
# state: [STATE]
# type: [TYPE]
# [OTHER COMMENTS]

[ACTION 1]
[QUESTION 1]

# redirect: [REDIRECT FLAG]
- [EMBEDDED TEST RFML ID]

Action with an embedded screenshot: {{ file.screenshot(./relative/path/to/screenshot.jpg) }}
Response with an embedded file download: {{ file.download(./relative/path/to/file.txt) }}

[ACTION 3]
[QUESTION 3]

... etc.
```

Required Fields:
- `RFML ID` - Unique identifier for your test. For newly generated tests, this will
be a UUID, but you are free to change it for easier reference (for example, your
login test might have the id `login_test`).
- `TITLE` - The title of your test.
- `START_URI` - The path used to direct the tester to the correct page to begin the test.
- `TYPE` - The type of test represented. Must be one of either `test` for regular, top-level
tests or `snippet` for any test that is meant to be embedded within another test. In other
words, if you're going to execute the test directly, it should be of type `test`; if you're
going to refer to it from another file (via a `- [EMBEDDED TEST RFML ID]` directive) it should
be of type `snippet`.
- `ACTION 1`, `ACTION 2`, ... - The directions for your tester to follow in this
step. You must have at least one step in your test.
- `QUESTION 1`, `QUESTION 2`, ... - The question you would like your tester to
answer in this step. You must have at least one step in your test.

Optional Fields:
- `SITE ID` - Site ID for the site this test is for. You can find your available
site IDs with the `sites` command. Sites can be configured at
https://app.rainforestqa.com/settings/sites.
- `PLATFORMS IDS` - Comma separated list of platforms for this test. You can reference
your available platforms with the `platforms` command. If left empty or omitted,
your test will default to using your account's default platforms.
- `TAGS` - Comma separated list of your desired tags for this test.
- `FEATURE_ID` - Feature ID for the feature that this test is a part of. You can
find your available feature IDs with the `features` command.
- `STATE` - State of the test. Valid states are `enabled`, `disabled` and `draft`.
- `OTHER COMMENTS` - Any comments you'd like to save to this test. All lines beginning with
`#` will be ignored by Rainforest unless they begin with a supported data field,
such as `tags` or `start_uri`.
- `REDIRECT FLAG` - A `true` or `false` flag to designate whether the tester should be
redirected. The default value is `true`. This flag is only applicable for embedded
tests and the first step of a test.
- `EMBEDDED TEST RFML ID` - Embed the steps of another test within the current test
using the embedded test's RFML ID.

For more information on embedding inline screenshots and file downloads,
[see our examples](./examples/inline_files.md).

### Command Line Options

Popular command line options are:
- `--platform ie8` or `--platforms ie8,chrome` - specify the platform(s) you wish to run against. This overrides the test level settings. Valid platforms can be found in your account settings.
- `--tag TAG_NAME` - filter tests by tag. Can be used multiple times for filtering by multiple tags.
- `--site-id SITE_ID` - filter tests by a specific site. You can see a list of your site IDs with `rainforest sites`.
- `--folder ID/--filter ID` - filter tests in specified folder.
- `--feature ID` - filter tests in a feature.
- `--run-group ID` - run/filter based on a run group. When used with `run`, this trigger a run from the run group; it can't be used in conjunction with other test filters.
- `--environment-id` - run your tests using this environment. Otherwise it will use your default environment
- `--conflict OPTION` - use the `cancel` option to cancel any runs in progress in the same environment as your new run. Use the `cancel-all` option to cancel all runs in progress.
- `--bg` - creates a run in the background and rainforest-cli exits immediately after. Do not use if you want rainforest-cli to track your run and exit with an error code upon run failure (ie: using Rainforest in your CI environment). Cannot be used together with `--max-reruns`.
- `--execution-method [crowd|automation|automation_and_crowd|on_premise]` - select how you wish your tests to be run. Your account may not have access to all methods. For more information, contact us at help@rainforestqa.com.
- `--wait RUN_ID` - wait for an existing run to finish instead of starting a new one, and exit with a non-0 code if the run fails. rainforest-cli will exit immediately if the run is already complete.
- `--fail-fast` - return an error as soon as the first failed result comes in (the run always proceeds until completion, but the CLI will return an error code early). If you don't use it, it will wait until 100% of the run is done. Has no effect with `--bg` and cannot be used together with `--max-reruns`.
- `--custom-url` - specify the URL for the run to use when testing against an ephemeral environment. This will create a new temporary environment for the run. Temporary environments will be automatically deleted 72 hours after they were last used.
- `--git-trigger` - only trigger a run when the last commit (for a git repo in the current working directory) has contains `@rainforest` and a list of one or more tags. E.g. "Fix checkout process. @rainforest #checkout" would trigger a run for everything tagged `checkout`. This over-rides `--tag` and any tests specified. If no `@rainforest` is detected it will exit 0.
- `--description "CI automatic run"` - add an arbitrary description for the run.
- `--release "1a2b3d"` - add an ID to associate the run with a release. Commonly used values are commit SHAs, build IDs, branch names, etc.
- `--flatten-steps` - Use with `rainforest download` to download your tests with steps extracted from embedded tests.
- `--test-folder /path/to/directory` - Use with `rainforest [new, upload, export]`. If this option is not provided, rainforest-cli will, in the case of 'new' create a directory, or in the case of 'upload' and 'export' use the directory, at the default path `./spec/rainforest/`.
- `--junit-file` - Create a junit xml report file with the specified name.  Must be run in foreground mode, or with the report command. Uses the rainforest
api to construct a junit report.  This is useful to track tests in CI such as Jenkins or Bamboo.
- `--import-variable-csv-file /path/to/csv/file.csv` - Use with `run` and `--import-variable-name` to upload new tabular variable values before your run to specify the path to your CSV file.
- `--import-variable-name NAME` - Use with `run` and `--import-variable-csv-file` to upload new tabular variable values before your run to specify the name of your tabular variable. You may also use this with the `csv-upload` command to update your variable without starting a run.
- `--single-use` - Use with `run` or `csv-upload` to flag your variable upload as `single-use`. See `--import-variable-csv-file` and `--import-variable-name` options as well.
- `--disable-telemetry` stops the cli sharing information about which CI system you may be using, and where you host your git repo (i.e. your git remote). Rainforest uses this to better integrate with CI tooling, and code hosting companies, it is not sold or shared. Disabling this may affect your Rainforest experience.
- `--max-reruns` - If set to a value > 0 and one or more tests fail, the CLI will create a new run with the failed tests a number of times before reporting failure. If `--junit-file <filename>` is also used, the JUnit reports of reruns will be saved under `<filename>.1`, `<filename>.2` etc. Cannot be used together with `--fail-fast`.
- `--automation-max-retries` - If set to a value > 0 and a test fails, it will be retried within the same run, up to that number of times. If all retries fail, we report failure, if there is a pass, we report that and stop retrying. The failed-then-passed attempts do not affect the final result of the run, but can be inspected in the web interface. See [our docs](https://help.rainforestqa.com/docs/test-retries) for more detail.

## Support

Email [help@rainforestqa.com](mailto:help@rainforestqa.com) if you're having trouble using the CLI or need help with integrating Rainforest in your CI or development workflow.

## Contributing

1. Fork it
1. Create a feature branch (`git checkout -b my-new-feature`)
1. Commit your changes (`git commit -am 'Add some feature'`)
1. Push to the branch (`git push origin my-new-feature`)
1. Create a new Pull Request

## Release Process

### Development PR
1. Branch from master
1. Do work
1. Open PR against master
1. Get review, and approval
1. Merge to master
### Changelog PR
1. Branch from master to update `CHANGELOG.md` to include the commit hashes and release date
1. Update the `version` constant in `rainforest-cli.go` following [semantic versioning](http://semver.org/)
1. Merge to master

### Releasing
1. **Docker** Tag `master` after merging: `git tag vX.Y.Z && git push --tags`
1. **GitHub** Wait for the CircleCI build to finish. This will create a [draft GitHub Release](https://github.com/rainforestapp/rainforest-cli/releases). Edit the description as appropriate and publish the release.
1. **Homebrew** Update https://github.com/rainforestapp/homebrew-public to use the latest URL and SHA256. Both can be found in the GitHub Release assets. Additionally, the SHA256 is output as part of the CircleCI `Release` job.
1. **Chocolatey** [Run the workflow here](https://github.com/rainforestapp/rainforest-cli-chocolatey/actions/workflows/chocolatey.yml) to build & release an updated Chocolatey package. Note, this uses the release you published earlier.

### Releasing a beta version (Docker / GitHub)
Simply tag a commit with an alpha or beta version.
```bash
git tag vX.Y.Z-alpha.N # or vX.Y.Z-beta.N
git push origin vX.Y.Z-alpha.N
```

### Rolling back
Should you have to rollback, you will need to:

1. Delete the GitHub release you need to rollback. When run without the `--skip-autoupdate` flag, the CLI will download the latest version from GitHub, thus auto-downgrading itself.
1. Go to GCP Container Registry:
    1. Delete the container you want to rollback
    1. Set the `latest` tag on the release you wish to rollback to
1. Revert the PR that caused the rollback in the first place
1. Check in Rainforest Admin who did (or could have) used the release and notify them via Support if there were any critical issues
