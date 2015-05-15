require "rainforest/cli/version"
require "rainforest/cli/options"
require "rainforest/cli/runner"
require "rainforest/cli/http_client"
require "rainforest/cli/git_trigger"
require "rainforest/cli/csv_importer"
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

      runner = Runner.new(options)
      runner.run

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
