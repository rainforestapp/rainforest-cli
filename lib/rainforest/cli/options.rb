require 'optparse'

module Rainforest
  module Cli
    class BrowserException < Exception
      def initialize browsers
        invalid_browsers = browsers - OptionParser::VALID_BROWSERS
        super "#{invalid_browsers.join(', ')} is not valid. Valid browsers: #{OptionParser::VALID_BROWSERS.join(', ')}"
      end
    end

    class OptionParser
      attr_reader :command, :token, :tags, :conflict, :browsers, :site_id,
                  :import_file_name, :import_name, :custom_url, :list_tests,
                  :list_sites, :list_runs

      VALID_BROWSERS = %w{chrome firefox safari ie8 ie9}.freeze

      def initialize(args)
        @args = args.dup
        @tags = []
        @browsers = nil

        @parsed = ::OptionParser.new do |opts|
          opts.on("--import-variable-csv-file FILE", "Import step variables; CSV data") do |value|
            @import_file_name = value
          end

          opts.on("--import-variable-name NAME", "Import step variables; Name of variable (note, will be replaced if exists)") do |value|
            @import_name = value
          end

          opts.on("--git-trigger", "Only run if the last commit contains @rainforestapp") do |value|
            @git_trigger = true
          end

          opts.on("--fg", "Run the tests in foreground.") do |value|
            @foreground = value
          end

          opts.on("--fail-fast", String, "Fail as soon as there is a failure (don't wait for completion)") do |value|
            @failfast = true
          end

          opts.on("--token TOKEN", String, "Your rainforest API token.") do |value|
            @token = value
          end

          opts.on("--tag TAG", String, "A tag to run the tests with") do |value|
            @tags << value
          end

          opts.on("--browsers LIST", "Run against the specified browsers") do |value|
            @browsers = value.split(',').map{|x| x.strip.downcase }.uniq

            raise BrowserException, @browsers unless (@browsers - VALID_BROWSERS).empty?
          end

          opts.on("--conflict MODE", String, "How should Rainforest handle existing in progress runs?") do |value|
            @conflict = value
          end

          opts.on("--site-id ID", Integer, "Only run tests for a specific site") do |value|
            @site_id = value
          end

          opts.on("--custom-url URL", String, "Use a custom url for this run. You will need to specify a site_id too for this to work.") do |value|
            @custom_url = value
          end

          opts.on("--list-tests", String, "List tests and exit") do |value|
            @list_tests = true
          end

          opts.on("--list-sites", String, "List sites and exit") do |value|
            @list_sites = true
          end

          opts.on("--list-runs [STATE]", String, "List all runs [with optional STATE] and exit") do |value|
            @list_runs = value ? value : true
          end
        end.parse!(@args)

        @command = @args.shift
        @tests = @args.dup
      end

      def tests
        @tests
      end

      def git_trigger?
        @git_trigger
      end

      def failfast?
        @failfast
      end

      def foreground?
        @foreground
      end
    end
  end
end

