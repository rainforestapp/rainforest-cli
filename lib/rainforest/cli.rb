# frozen_string_literal: true
require 'erb'
require 'json'
require 'logger'
require 'rainforest/cli/version'
require 'rainforest/cli/constants'
require 'rainforest/cli/options'
require 'rainforest/cli/runner'
require 'rainforest/cli/http_client'
require 'rainforest/cli/git_trigger'
require 'rainforest/cli/csv_importer'
require 'rainforest/cli/test_parser'
require 'rainforest/cli/test_files'
require 'rainforest/cli/remote_tests'
require 'rainforest/cli/validator'
require 'rainforest/cli/exporter'
require 'rainforest/cli/deleter'
require 'rainforest/cli/uploader'
require 'rainforest/cli/resources'
require 'rainforest/cli/junit_outputter'

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
