require 'optparse'

module Rainforest
  module Cli
    class OptionParser
      attr_reader :command, :token, :tags, :conflict

      def initialize(args)
        @args = args.dup
        @tags = []

        @parsed = ::OptionParser.new do |opts|
          opts.on("--fg", "Run the tests in foreground.") do |value|
            @foreground = value
          end

          opts.on("--fail-fast", String, "Fail as soon as there is a failure (don't wait for completion)") do |value|
            @failfast = true
          end

          opts.on("--token TOKEN", String, "Your rainforest API token.") do |value|
            @token = value
          end

          opts.on("--tag TOKEN", String, "A tag to run the tests with") do |value|
            @tags << value
          end

          opts.on("--conflict MODE", String, "How should Rainforest handle existing in progress runs?") do |value|
            @conflict = value
          end
        end.parse!(@args)

        @command = @args.shift
        @tests = @args.dup
      end

      def tests
        @tests
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

