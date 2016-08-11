# frozen_string_literal: true
require 'optparse'

module RainforestCli
  class OptionParser
    attr_writer :file_name, :tags
    attr_reader :command, :token, :tags, :conflict, :browsers, :site_id, :environment_id,
                :import_file_name, :import_name, :custom_url, :description, :folder,
                :debug, :file_name, :test_folder, :embed_tests, :app_source_url, :crowd, :run_id,
                :junit_file

    TOKEN_NOT_REQUIRED = %w{new validate}.freeze

    def initialize(args)
      @args = args.dup
      @tags = []
      @browsers = nil
      @debug = false
      @junit_file = nil
      @run_id = nil
      @token = ENV['RAINFOREST_API_TOKEN']

      # NOTE: Disabling line length cop to allow for consistency of syntax
      # rubocop:disable Metrics/LineLength
      @parsed = ::OptionParser.new do |opts|
        opts.set_program_name 'Rainforest CLI'
        opts.version = RainforestCli::VERSION

        opts.on('--debug') do
          @debug = true
        end

        opts.on('--file') do |value|
          @file_name = value
        end

        opts.on('--test-folder FILE_PATH', 'Specify the test folder. Defaults to spec/rainforest if not set.') do |value|
          @test_folder = value
        end

        opts.on('--crowd TYPE', 'Specify the crowd type, defaults to "default". Options are default or on_premise_crowd if available for your account.') do |value|
          @crowd = value
        end

        opts.on('--import-variable-csv-file FILE', 'Import step variables; CSV data') do |value|
          @import_file_name = value
        end

        opts.on('--import-variable-name NAME', 'Import step variables; Name of variable (note, will be replaced if exists)') do |value|
          @import_name = value
        end

        opts.on('--git-trigger', 'Only run if the last commit contains @rainforestapp') do |_value|
          @git_trigger = true
        end

        opts.on('--fg', 'Run the tests in foreground.') do |value|
          @foreground = value
        end

        opts.on('--wait ID', Integer, 'Wait for run #ID to complete') do |value|
          @wait = true
          @run_id = value
        end

        opts.on('--fail-fast', String, "Fail as soon as there is a failure (don't wait for completion)") do |_value|
          @failfast = true
        end

        opts.on('--token API_TOKEN', String, 'Your rainforest API token.') do |value|
          @token = value
        end

        opts.on('--tag TAG', String, 'A tag to run the tests with') do |value|
          @tags << value
        end

        opts.on('--folder ID', 'Run tests in the specified folders') do |value|
          @folder = value
        end

        opts.on('--browsers LIST', 'Run against the specified browsers') do |value|
          @browsers = value.split(',').map {|x| x.strip.downcase }.uniq
        end

        opts.on('--conflict MODE', String, 'How should Rainforest handle existing in progress runs?') do |value|
          @conflict = value
        end

        opts.on('--environment-id ID', Integer, 'Run using this environment. If excluded, will use your default') do |value|
          @environment_id = value
        end

        opts.on('--site-id ID', Integer, 'Only run tests for a specific site') do |value|
          @site_id = value
        end

        opts.on('--custom-url URL', String, 'Use a custom url for this run. You will need to specify a site_id too for this to work.') do |value|
          @custom_url = value
        end

        opts.on('--description DESCRIPTION', 'Add a description for the run.') do |value|
          @description = value
        end

        opts.on('--embed-tests', 'Export tests without expanding embedded test steps') do |_value|
          @embed_tests = true
        end

        opts.on('--app-source-url FILE', 'URL for mobile app download (in beta)') do |value|
          @app_source_url = value
        end

        opts.on('--junit-file FILE', 'Gather the results of a run and create junit output in FILE.xml, must be run with --fg') do |value|
          @junit_file = value
        end

        opts.on('--run-id ID', 'Gather the results of a completed run, must be run with export and --junit-file') do |value|
          @run_id = value
        end

        opts.on_tail('--help', 'Display help message and exit') do |_value|
          puts opts
          exit 0
        end

        opts.on_tail('--version', 'Display gem version') do
          puts opts.ver
          exit 0
        end

      end.parse!(@args)
      # rubocop:enable Metrics/LineLength
      # NOTE: end Rubocop exception

      @command = @args.shift

      if ['new', 'rm'].include?(@command)
        @file_name = @args.shift

        if @file_name && @command == 'rm'
          @file_name = File.expand_path(@file_name) if @file_name
        end
      end

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

    def wait?
      !!@wait
    end

    def foreground?
      @foreground
    end

    def junit_file?
      @junit_file
    end

    def validate!
      if !TOKEN_NOT_REQUIRED.include?(command)
        unless token
          raise ValidationError, 'You must pass your API token using: --token TOKEN'
        end
      end

      if custom_url && site_id.nil?
        raise ValidationError, 'The site-id and custom-url options are both required.'
      end

      if import_file_name && import_name
        unless File.exist?(import_file_name)
          raise ValidationError, "Input file: #{import_file_name} not found"
        end
      elsif import_file_name || import_name
        raise ValidationError, 'You must pass both --import-variable-csv-file and --import-variable-name'
      end

      if command == 'rm' && file_name.nil?
        raise ValidationError, 'You must include a file name'
      end

      if command == 'run' && junit_file?
        unless foreground?
          raise ValidationError, 'You can only generate junit test output in foreground mode'
        end
      end

      if command == 'report'
        if junit_file.nil?
          raise ValidationError, 'You must specify a junit ouptut filename'
        end
        if run_id.nil?
          raise ValidationError, 'You must specify a run-id to generate a report for'
        end
      end

      true
    end

    class ValidationError < RuntimeError
    end
  end
end
