[![Build Status](https://travis-ci.org/rainforestapp/rainforest-cli.png?branch=master)](https://travis-ci.org/rainforestapp/rainforest-cli)

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

Run all of your tests

```bash
rainforest run all --token YOUR_TOKEN_HERE
```

Run all in the foreground and report

```bash
rainforest run all --fg --token YOUR_TOKEN_HERE
```

Run all tests with tag 'run-me' and abort previous in-progress runs.

```bash
rainforest run --tag run-me --fg --conflict abort --token YOUR_TOKEN_HERE
```


## Options

### General

Required:
- `--token <your-rainforest-token>` - you must supply your token (get it from any tests API tab)


### Running Tests
The most popular options are:

- `--browsers ie8` or `--browsers ie8,chrome` - specficy the browsers you wish to run against. This overrides the test own settings. Valid browsers are ie8, ie9, chrome, firefox and safari.
- `--tag run-me` - only run tests which have this tag (recommended if you have lots of [test-steps](http://docs.rainforestqa.com/pages/example-test-suite.html#test_steps))!)
- `--site-id` - only run tests for a specific site. Get in touch with us for help on getting that you site id if you are unable to.
- `--environment-id` - run your tests using this environment. Otherwise it will use your default environment
- `--conflict abort` - if you trigger rainforest more than once, anything running on the same environment will be aborted and a fresh run started
- `--fg` - results in the foreground - rainforest-cli will not return until the run is complete. This is what you want to make the build pass / fail dependent on rainforest results
- `--fail-fast` - fail the build as soon as the first failed result comes in. If you don't pass this it will wait until 100% of the run is done. Use with `--fg`.
- `--custom-url` - use a custom url for this run. Example use case: an ad-hoc QA environment with [Fourchette](https://github.com/rainforestapp/fourchette). You will need to specify a `site_id` too for this to work. Note that we will be creating a new environment for this particular run.
- `--git-trigger` - only trigger a run when the last commit (for a git repo in the current working directory) has contains `@rainforest` and a list of one or more tags. E.g. "Fix checkout process. @rainforest #checkout" would trigger a run for everything tagged `checkout`. This over-rides `--tag` and any tests specified. If no `@rainforest` is detected it will exit 0.
- `--description "CI automatic run"` - add an arbitrary description for the run.

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
