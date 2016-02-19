# frozen_string_literal: true
# coding: utf-8
lib = File.expand_path('../lib', __FILE__)
$LOAD_PATH.unshift(lib) unless $LOAD_PATH.include?(lib)
require 'rainforest/cli/version'

Gem::Specification.new do |spec|
  spec.name          = 'rainforest-cli'
  spec.version       = RainforestCli::VERSION
  spec.authors       = ['Simon Mathieu', 'Russell Smith']
  spec.email         = ['simon@rainforestqa.com', 'russ@rainforestqa.com']
  spec.description   = %q{Command line utility for Rainforest QA}
  spec.summary       = %q{Command line utility for Rainforest QA}
  spec.homepage      = 'https://www.rainforestqa.com/'
  spec.license       = 'MIT'

  spec.files         = `git ls-files`.split($/)
  spec.executables   = spec.files.grep(%r{^bin/}) { |f| File.basename(f) }
  spec.test_files    = spec.files.grep(%r{^(test|spec|features)/})
  spec.require_paths = ['lib']

  spec.add_dependency 'httparty', '~> 0.13.7'
  spec.add_dependency 'parallel', '~> 1.6.1'
  spec.add_dependency 'ruby-progressbar', '~> 1.7.5'
  spec.add_dependency 'rainforest', '~> 2.1.0'
  spec.add_development_dependency 'bundler', '~> 1.3'
  spec.add_development_dependency 'rake', '~> 10.4.2'
end
