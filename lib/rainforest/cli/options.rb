require 'optparse'

module RainforestCli
  class BrowserException < Exception
    def initialize browsers
      invalid_browsers = browsers - OptionParser::VALID_BROWSERS
      super "#{invalid_browsers.join(', ')} is not valid. Valid browsers: #{OptionParser::VALID_BROWSERS.join(', ')}"
    end
  end

  class OptionParser
    attr_reader :command, :token, :tags, :conflict, :browsers

    VALID_BROWSERS = %w{chrome firefox safari ie8 ie9}.freeze

    def initialize(args)
      @args = args.dup
      @tags = []
      @browsers = nil

      @parsed = ::OptionParser.new do |opts|
        puts opts.inspect
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

