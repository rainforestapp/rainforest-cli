# Rainforest-cli

A command line interface to interact with RainforestQA.

## Installation

    $ gem install rainforest-cli

## Usage

Run all of your tests

    rainforest run all

Run and report

    rainforest run --fg all

Run all tests with tag 'run-me' and abort previous in-progress runs.

    rainforest run --token a8b2fe2dd7360ec72aaef0a0312fa7fa --tag run-me --conflict abort

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request
