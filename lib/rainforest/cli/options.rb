require 'optparse'

module Rainforest
  module Cli
    class OptionParser
      attr_reader :command, :token

      def initialize(args)
        @args = args.dup

        @parsed = ::OptionParser.new do |opts|
          opts.on("--fg", "Run the tests in foreground.") do |value|
            @foreground = value
          end

          opts.on("--token TOKEN", String, "Your rainforest API token.") do |value|
            @token = value
          end
        end.parse!(@args)

        @command = @args.shift
        @tests = @args.dup
      end

      def tests
        @tests
      end

      def foreground?
        @foreground
      end
    end
  end
end

