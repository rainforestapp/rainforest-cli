[![CircleCI](https://circleci.com/gh/rainforestapp/rainforest-cli.svg?style=shield)](https://circleci.com/gh/rainforestapp/rainforest-cli)

# Rainforest-cli

A command line interface to interact with [Rainforest QA](https://www.rainforestqa.com/).

This is the easiest way to integrate Rainforest with your deploy scripts or CI server. See [our documentation](https://intercom.help/rainforest/developer-tools/cli-docs/rainforest-cli-for-continuous-integration) on the subject.

## Installation

If you are on a mac and use brew, you can run:

```bash
brew install rainforestapp/public/rainforest-cli
```

If not, follow the directions on our [download page](https://dl.equinox.io/rainforest_qa/rainforest-cli/stable). For non-homebrew mac users, please make sure you use the "Install via the command line" instructions.

The CLI will check for updates and automatically update itself on every use unless the `--skip-update` global flag is given.

For migration directions from 1.x, please read our [migration guide](./migration.md).

## Basic Usage
To use the cli client, you'll need your API token from your [integrations settings page](https://app.rainforestqa.com/settings/integrations).

In order to access your account from the CLI, set the `RAINFOREST_API_TOKEN` environment variable
to your API token. Alternatively, you may pass your token with the global `--token` flag.

CLI Commands are formatted as follows:
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

Run all tests with tag 'run-me' and abort previous in-progress runs.

```bash
rainforest run --tag run-me --conflict abort
```

Run all tests and generate a junit xml report.

```bash
rainforest run all --junit-file results.xml
```

Run individual tests in the foreground and report.

```bash
rainforest run <test_id1> <test_id2>
```

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

Run-level setting options (`--browsers`, `--environment_id`, etc) behave the same for `run -f`. Other test filtering options (such as `--run-group`, `--site`, etc) cannot be used in conjunction with `run -f`.

#### Viewing Account Specific Information

See a list of all of your sites and their IDs
```bash
rainforest sites
```

See a list of all of your smart folders and their IDs
```bash
rainforest folders
```

See a list of all of your browsers and their IDs
```bash
rainforest browsers
```

See a list of all of your features and their IDs
```bash
rainforest features
```

See a list of all of your run groups and their IDs
```bash
rainforest run-groups
```

To generate a junit xml report for a test run which has already completed
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

## Options

### Global

- `--token <your-rainforest-token>` - supply your token (get it from any tests API tab), if not set in `RAINFOREST_API_TOKEN` environment variable
- `--skip-update` - Do not automatically check for CLI updates.

### Writing Tests
Rainforest Tests written using RFML have the following format

```
#! [RFML ID]
# title: [TITLE]
# start_uri: [START_URI]
# tags: [TAGS]
# site_id: [SITE ID]
# browsers: [BROWSER IDS]
# feature_id: [FEATURE_ID]
# state: [STATE]
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
- `ACTION 1`, `ACTION 2`, ... - The directions for your tester to follow in this
step. You must have at least one step in your test.
- `QUESTION 1`, `QUESTION 2`, ... - The question you would like your tester to
answer in this step. You must have at least one step in your test.

Optional Fields:
- `SITE ID` - Site ID for the site this test is for. You can find your available
site IDs with the `sites` command. Sites can be configured at
https://app.rainforestqa.com/settings/sites.
- `BROWSER IDS` - Comma separated list of browsers for this test. You can reference
your available browsers with the `browsers` command. If left empty or omitted,
your test will default to using your account's default browsers.
- `TAGS` - Comma separated list of your desired tags for this test.
- `FEATURE_ID` - Feature ID for the feature that this test is a part of. You can
find your available feature IDs with the `features` command.
- `STATE` - State of the test. Valid states are `enabled` and `disabled`.
- `OTHER COMMENTS` - Any comments you'd like to save to this test. All lines beginning with
`#` will be ignored by Rainforest unless they begin with a supported data field,
such as `tags` or `start_uri`.
- `REDIRECT FLAG` - A `true` or `false` flag to designate whether the tester should be
redirected. The default value is `true`. This flag is only applicable for embedded
tests and the first step of a test.
- `EMBEDDED TEST RFML ID` - Embed the steps of another test within the current test
using the embedded test's RFML ID.

For more information on embedding inline screenshots and file downloads,
[see our examples](./examples/inline_files.md)

For more information on test writing, please visit our [documentation](http://support.rainforestqa.com/hc/en-us/sections/200585603-Writing-Tests).

### Command Line Options

Popular command line options are:
- `--browsers ie8` or `--browsers ie8,chrome` - specify the browsers you wish to run against. This overrides the test own settings. Valid browsers can be found in your account settings.
- `--tag TAG_NAME` - filter tests by tag. Can be used multiple times for filtering by multiple tags.
- `--site-id SITE_ID` - filter tests by a specific site. You can see a list of your site IDs with `rainforest sites`.
- `--folder ID/--filter ID` - filter tests in specified folder.
- `--feature ID` - filter tests in a feature.
- `--run-group ID` - run/filter based on a run group. When used with `run`, this trigger a run from the run group; it can't be used in conjunction with other test filters.
- `--environment-id` - run your tests using this environment. Otherwise it will use your default environment
- `--conflict OPTION` - use the `abort` option to abort any runs in progress in the same environment as your new run. use the `abort-all` option to abort all runs in progress.
- `--bg` - creates a run in the background and rainforest-cli exits immediately after. Do not use if you want rainforest-cli to track your run and exit with an error code upon run failure (ie: using Rainforest in your CI environment).
- `--crowd [default|on_premise_crowd]` - select your crowd of testers for clients with on premise testers. For more information, contact us at help@rainforestqa.com.
- `--wait RUN_ID` - wait for an existing run to finish instead of starting a new one, and exit with a non-0 code if the run fails. rainforest-cli will exit immediately if the run is already complete.
- `--fail-fast` - fail the build as soon as the first failed result comes in. If you don't pass this it will wait until 100% of the run is done. Has no effect with `--bg`.
- `--custom-url` - use a custom url for this run to use an ad-hoc QA environment on all tests. You will need to specify a `site_id` too for this to work. Note that we will be creating a new environment for your account for this particular run.
- `--git-trigger` - only trigger a run when the last commit (for a git repo in the current working directory) has contains `@rainforest` and a list of one or more tags. E.g. "Fix checkout process. @rainforest #checkout" would trigger a run for everything tagged `checkout`. This over-rides `--tag` and any tests specified. If no `@rainforest` is detected it will exit 0.
- `--description "CI automatic run"` - add an arbitrary description for the run.
- `--flatten-steps` - Use with `rainforest download` to download your tests with steps extracted from embedded tests.
- `--test-folder /path/to/directory` - Use with `rainforest [new, upload, export]`. If this option is not provided, rainforest-cli will, in the case of 'new' create a directory, or in the case of 'upload' and 'export' use the directory, at the default path `./spec/rainforest/`.
- `--junit-file` - Create a junit xml report file with the specified name.  Must be run in foreground mode, or with the report command. Uses the rainforest
api to construct a junit report.  This is useful to track tests in CI such as Jenkins or Bamboo.
- `--run-id` - Only used with the report command.  Specify a past rainforest run by ID number to generate a report for.
- `--import-variable-csv-file /path/to/csv/file.csv` - Use with `run` and `--import-variable-name` to upload new tabular variable values before your run to specify the path to your CSV file.
- `--import-variable-name NAME` - Use with `run` and `--import-variable-csv-file` to upload new tabular variable values before your run to specify the name of your tabular variable. You may also use this with the `csv-upload` command to update your variable without starting a run.
- `--single-use` - Use with `run` or `csv-upload` to flag your variable upload as `single-use`. See `--import-variable-csv-file` and `--import-variable-name` options as well.

## Support

Email [help@rainforestqa.com](mailto:help@rainforestqa.com) if you're having trouble using this gem or need help to integrate Rainforest in your CI or deployment flow.

## Contributing

1. Fork it
2. Make sure you init the submodules (`git submodule init && git submodule update`)
3. Create your feature branch (`git checkout -b my-new-feature`)
4. Commit your changes (`git commit -am 'Add some feature'`)
5. Push to the branch (`git push origin my-new-feature`)
6. Create new Pull Request

## Release process

Check the `circle.yml` for the latest, but currently merging to master will build and deploy to the following Equinox channels:

Tag                             | Channels
--------------------------------|-------------
No tag                          | dev
vX.Y.Z-alpha.N or vX.Y.Z-beta.N | beta, dev
vX.Y.Z                          | stable, beta, dev

Development + release process is:

1. Branch from master
2. Do work
3. Open PR against master
4. Merge to master
5. Branch from master to update `CHANGELOG.md` to include the commit hashes and release date
6. Update the `version` constant in `rainforest-cli.go` following [semvar](http://semver.org/)
7. Merge to master
8. Tag the master branch with the release:
```bash
   git tag vX.Y.Z or vX.Y.Z-alpha.N or vX.Y.Z-beta.N
   git push origin vX.Y.Z
```
9. Merge to master to release to stable/beta/dev
10. Add release to Github [release page](https://github.com/rainforestapp/rainforest-cli/releases)
