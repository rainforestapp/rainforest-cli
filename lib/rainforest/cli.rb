require "rainforest/cli/version"
require "rainforest/cli/options"
require "rainforest/cli/runner"
require "rainforest/cli/http_client"
require "rainforest/cli/git_trigger"
require "rainforest/cli/csv_importer"
require "httparty"
require "json"
require "logger"

module Rainforest
  module Cli
    def self.start(args)
      options = OptionParser.new(args)
      OptionParser.new(['--help']) if args.size == 0

      validate_options!(options)

      runner = Runner.new(options)
      runner.run

      true
    end

    def self.validate_options!(options)
      unless options.token
        logger.fatal "You must pass your API token using: --token TOKEN"
        exit 2
      end

      if options.custom_url && options.site_id.nil?
        logger.fatal "The site-id and custom-url options are both required."
        exit 2
      end

      if options.import_file_name && options.import_name
        unless File.exists?(options.import_file_name)
          logger.fatal "Input file: #{options.import_file_name} not found"
          exit 2
        end

      elsif options.import_file_name || options.import_name
        logger.fatal "You must pass both --import-variable-csv-file and --import-variable-name"
        exit 2
      end
    end

    def self.logger
      @logger ||= Logger.new(STDOUT)
    end

    def self.logger=(logger)
      @logger = logger
    end
  end
end
