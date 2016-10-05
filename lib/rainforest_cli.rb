# frozen_string_literal: true
require 'erb'
require 'json'
require 'logger'
require 'rainforest_cli/version'
require 'rainforest_cli/constants'
require 'rainforest_cli/options'
require 'rainforest_cli/commands'
require 'rainforest_cli/runner'
require 'rainforest_cli/http_client'
require 'rainforest_cli/git_trigger'
require 'rainforest_cli/csv_importer'
require 'rainforest_cli/test_parser'
require 'rainforest_cli/test_files'
require 'rainforest_cli/remote_tests'
require 'rainforest_cli/validator'
require 'rainforest_cli/exporter'
require 'rainforest_cli/deleter'
require 'rainforest_cli/uploader'
require 'rainforest_cli/resources'
require 'rainforest_cli/junit_outputter'
require 'rainforest_cli/reporter'

module RainforestCli
  def self.start(args)
    options = OptionParser.new(args)
    commands = Commands.new do |c|
      c.add('run', 'Run your tests on Rainforest') { Runner.new(options).run }
      c.add('new', 'Create a new RFML test') { TestFiles.new(options).create_file }
      c.add('validate', 'Validate your RFML tests') { Validator.new(options).validate }
      c.add('upload', 'Upload your RFML tests') { Uploader.new(options).upload }
      c.add('rm', 'Remove an RFML test locally and remotely') { Deleter.new(options).delete }
      c.add('export', 'Export your remote Rainforest tests to RFML') { Exporter.new(options).export }
      c.add('csv-upload', 'Upload a new tabular variable from a CSV file') { CSVImporter.new(options).import }
      c.add('report', 'Create a JUnit report from your run results') { Reporter.new(options).report }
    end

    @http_client = HttpClient.new(token: options.token)
    ::Rainforest.api_key = options.token

    if args.size == 0
      commands.print_documentation
      options.print_documentation
    end

    begin
      options.validate!
    rescue OptionParser::ValidationError => e
      logger.fatal e.message
      exit 2
    end

    commands.call(options.command) if options.command
    true
  end

  def self.logger
    @logger ||= Logger.new(STDOUT)
  end

  def self.logger=(logger)
    @logger = logger
  end

  def self.http_client
    @http_client
  end
end
