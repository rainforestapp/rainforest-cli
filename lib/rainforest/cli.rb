require "rainforest/cli/version"
require "rainforest/cli/options"
require "rainforest/cli/runner"
require "rainforest/cli/http_client"
require "rainforest/cli/git_trigger"
require "rainforest/cli/csv_importer"
require "rainforest/cli/test_importer"
require "rainforest/cli/test_parser"
require "erb"
require "httparty"
require "json"
require "logger"

module Rainforest
  module Cli
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
      when 'run'
        runner = Runner.new(options)
        runner.run
      when 'new'
        t = TestImporter.new(options)
        t.create_new
      when 'validate'
        t = TestImporter.new(options)
        t.validate
      else
        logger.fatal "Unknown command"
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
end
