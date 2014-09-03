[![Build Status](https://travis-ci.org/rainforestapp/rainforest-cli.png?branch=master)](https://travis-ci.org/rainforestapp/rainforest-cli)

# Rainforest-cli

A command line interface to interact with RainforestQA.

## Installation

    $ gem install rainforest-cli

## Basic Usage
To use the cli client, you'll need your API token from a test settings page from inside [Rainforest](https://app.rainforestqa.com/).

Run all of your tests

    rainforest run all --token YOUR_TOKEN_HERE

Run all in the foreground and report

    rainforest run all --fg --token YOUR_TOKEN_HERE

Run all tests with tag 'run-me' and abort previous in-progress runs.

    rainforest run --tag run-me --fg --conflict abort --token YOUR_TOKEN_HERE 


## Options

Required:
- `--token <your-rainforest-token>` - you must supply your token (get it from any tests API tab)

The options are:

- `--browsers ie8` or `--browsers ie8,chrome` - specficy the browsers you wish to run against. This overrides the test own settings. Valid browsers are ie8, ie9, chrome, firefox and safari.
- `--tag run-me` - only run tests which have this tag (recommended if you have lots of [test-steps](http://docs.rainforestqa.com/pages/example-test-suite.html#test_steps))!)
- `--conflict abort` - if you trigger rainforest more than once, anything running will be aborted and a fresh run started
- `--fail-fast` - fail the build as soon as the first failed result comes in. If you don't pass this it will wait until 100% of the run is done
- `--fg` - results in the foreground - this is what you want to make the build pass / fail dependent on rainforest results 
- `--site-id` - only run tests for a specific site
- `--custom-url` - use a custom url for this run. Example use case: an ad-hoc QA environment with [Fourchette](https://github.com/rainforestapp/fourchette). You will need to specify a `site_id` too for this to work. Please note that we will be creating environments under the hood and will not affect your test permanently.


## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
