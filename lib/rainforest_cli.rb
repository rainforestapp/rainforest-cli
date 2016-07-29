# frozen_string_literal: true
require 'erb'
require 'json'
require 'logger'
require 'rainforest_cli/version'
require 'rainforest_cli/constants'
require 'rainforest_cli/options'
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
    OptionParser.new(['--help']) if args.size == 0

    begin
      options.validate!
    rescue OptionParser::ValidationError => e
      logger.fatal e.message
      exit 2
    end

    case options.command
    when 'run' then Runner.new(options).run
    when 'new' then TestFiles.new(options).create_file
    when 'validate' then Validator.new(options).validate
    when 'upload' then Uploader.new(options).upload
    when 'rm' then Deleter.new(options).delete
    when 'export' then Exporter.new(options).export
    when 'report' then Reporter.new(options).report
    when 'sites', 'folders', 'browsers'
      Resources.new(options).public_send(options.command)
    else
      logger.fatal 'Unknown command'
      exit 2
    end

    true
  end

  def self.logger
    @logger ||= Logger.new(STDOUT)
  end

  def self.logger=(logger)
    @logger = logger
  end
end
