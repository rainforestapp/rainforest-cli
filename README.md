[![Build Status](https://travis-ci.org/rainforestapp/rainforest-cli.png?branch=master)](https://travis-ci.org/rainforestapp/rainforest-cli)

[![Gem Version](https://badge.fury.io/rb/rainforest-cli.svg)](https://badge.fury.io/rb/rainforest-cli)

# Rainforest-cli

A command line interface to interact with RainforestQA.

This is the easiest way to integrate Rainforest with your deploy scripts or CI server. See [our documentation](http://support.rainforestqa.com/hc/en-us/articles/205876128-Continuous-Integration) on the subject.

## Installation

You can install rainforest-cli with the [gem](https://rubygems.org/) utility.

```bash
gem install rainforest-cli
```

Alternatively, you can add to your Gemfile if you're in a ruby project. This is *not recommended* for most users. The reason being that we update this gem frequently and you usually want to ensure you have the latest version.

```ruby
gem "rainforest-cli", require: false
```

## Basic Usage
To use the cli client, you'll need your API token from a test settings page from inside [Rainforest](https://app.rainforestqa.com/).

You can either pass the token with `--token YOUR_TOKEN_HERE` CLI option, or put it in the `RAINFOREST_API_TOKEN` environment variable.

You can customize the amount of threads to use when making HTTP requests by setting the
`RAINFOREST_THREADS` to the number of threads you'd like to use. If you are noticing
API errors when fetching or updating multiple tests, lowering this value can help.
The default value is `16`.

## Options

#### Running Tests

Run all tests.

```bash
rainforest run all
```

Run all in the foreground and report.

```bash
rainforest run all --fg
```

Run all tests with tag 'run-me' and abort previous in-progress runs.

```bash
rainforest run --tag run-me --fg --conflict abort
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

Upload tests to Rainforest

```bash
rainforest upload
```

Remove RFML file and remove test from Rainforest test suite.

```bash
rainforest rm /path/to/test/file.rfml
```

Export all tests from Rainforest

```bash
rainforest export
```

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

## Options

### General

- `--token <your-rainforest-token>` - supply your token (get it from any tests API tab), if not set in `RAINFOREST_API_TOKEN` environment variable

### Writing Tests
Rainforest Tests written using RFML have the following format

```
#! [RFML ID]
# title: [TITLE]
# start_uri: [START_URI]
# tags: [TAGS]
# site_id: [SITE ID]
# browsers: [BROWSER IDS]
# [OTHER COMMENTS]

[ACTION 1]
[QUESTION 1]

# redirect: [REDIRECT FLAG]
- [EMBEDDED TEST RFML ID]

[ACTION 2]
[QUESTION 2]

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
your available browsers with the `browsers` command.
- `TAGS` - Comma separated list of your desired tags for this test.
- `OTHER COMMENTS` - Any comments you'd like to save to this test. All lines beginning with
`#` will be ignored by Rainforest unless they begin with a supported data field,
such as `tags` or `start_uri`.
- `REDIRECT FLAG` - A `true` or `false` flag to designate whether the tester should be
redirected. The default value is `true`. This flag is only applicable for embedded
tests and the first step of a test.
- `EMBEDDED TEST RFML ID` - Embed the steps of another test within the current test
using the embedded test's RFML ID.

For more information on test writing, please visit our [documentation](http://support.rainforestqa.com/hc/en-us/sections/200585603-Writing-Tests).

### Command Line Options

Popular command line options are:
- `--browsers ie8` or `--browsers ie8,chrome` - specify the browsers you wish to run against. This overrides the test own settings. Valid browsers can be found in your account settings.
- `--tag run-me` - only run tests which have this tag (recommended if you have lots of [test-steps](http://docs.rainforestqa.com/pages/example-test-suite.html#test_steps))!)
- `--site-id` - only run tests for a specific site. Get in touch with us for help on getting that you site id if you are unable to.
- `--folder ID` - run tests in specified folder.
- `--environment-id` - run your tests using this environment. Otherwise it will use your default environment
- `--conflict OPTION` - use the `abort` option to abort any runs in progress in the same environment as your new run. use the `abort-all` option to abort all runs in progress.
- `--fg` - results in the foreground - rainforest-cli will not return until the run is complete. This is what you want to make the build pass / fail dependent on rainforest results
- `--wait RUN_ID` - wait for an existing run to finish instead of starting a new one, and exit with a non-0 code if the run fails. rainforest-cli will exit immediately if the run is already complete.
- `--fail-fast` - fail the build as soon as the first failed result comes in. If you don't pass this it will wait until 100% of the run is done. Use with `--fg`.
- `--custom-url` - use a custom url for this run. Example use case: an ad-hoc QA environment with [Fourchette](https://github.com/rainforestapp/fourchette). You will need to specify a `site_id` too for this to work. Note that we will be creating a new environment for this particular run.
- `--git-trigger` - only trigger a run when the last commit (for a git repo in the current working directory) has contains `@rainforest` and a list of one or more tags. E.g. "Fix checkout process. @rainforest #checkout" would trigger a run for everything tagged `checkout`. This over-rides `--tag` and any tests specified. If no `@rainforest` is detected it will exit 0.
- `--description "CI automatic run"` - add an arbitrary description for the run.
- `--embed-tests` - Use with `rainforest export` to export your tests without extracting the
steps of an embedded test.
- `--test-folder /path/to/directory` - Use with `rainforest [new, upload, export]`. If this option is not provided, rainforest-cli will, in the case of 'new' create a directory, or in the case of 'upload' and 'export' use the directory, at the default path `./spec/rainforest/`.
- `--import-variable-csv-file /path/to/csv/file.csv` - Use with `run` and `--import-variable-name` to upload new tabular variable values before your run to specify the path to your CSV file.
- `--import-variable-name NAME` - Use with `run` and `import-variable-csv-file` to upload
new tabular variable values before your run to specify the name of your tabular variable.


#### Specifying Test IDs
Any integers input as arguments in the command line arguments are treated as
test IDs taken from the Rainforest dashboard. ie:

`rainforest run --token $TOKEN 1232 3212` - will export only tests
1232 and 3212. The `export` and `run` commands and are otherwise ignored.

All other argument types should be specified as seen above.


More detailed info on options can be [found here](https://github.com/rainforestapp/rainforest-cli/blob/master/lib/rainforest/cli/options.rb#L23-L74).

## Support

Email [help@rainforestqa.com](mailto:help@rainforestqa.com) if you're having trouble using this gem or need help to integrate Rainforest in your CI or deployment flow.

## Contributing

1. Fork it
2. Make sure you init the submodules (`git submodule init && git submodule update`)
3. Create your feature branch (`git checkout -b my-new-feature`)
4. Commit your changes (`git commit -am 'Add some feature'`)
5. Push to the branch (`git push origin my-new-feature`)
6. Create new Pull Request
