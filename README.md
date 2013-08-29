# Rainforest-cli

A command line interface to interact with RainforestQA.

## Installation

    $ gem install rainforest-cli

## Usage
To use the cli client, you'll need your API token from a test settings page from inside [Rainforest](https://app.rainforestqa.com/).

Run all of your tests

    rainforest run all --token YOUR_TOKEN_HERE

Run and report

    rainforest run --fg all --token YOUR_TOKEN_HERE

Run all tests with tag 'run-me' and abort previous in-progress runs.

    rainforest run --tag run-me --conflict abort --token YOUR_TOKEN_HERE 

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
