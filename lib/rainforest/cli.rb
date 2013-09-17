require "rainforest/cli/version"
require "rainforest/cli/options"
require "httparty"
require "json"

module Rainforest
  module Cli 
    API_URL = 'https://app.rainforestqa.com/api/1/runs'

    def self.start(args)
      @options = OptionParser.new(args)

      post_opts = {}
      if !@options.tags.empty?
        post_opts[:tags] = @options.tags
      else
        post_opts[:tests] = @options.tests
      end

      post_opts[:conflict] = @options.conflict if @options.conflict
      post_opts[:browsers] = @options.browsers if @options.browsers
      post_opts[:gem_version] = Rainforest::Cli::VERSION

      puts "Issuing run"

      response = post(API_URL, post_opts)

      if response['error']
        puts "Error starting your run: #{response['error']}"
        exit
      end

      run_id = response["id"]
      running = true

      return unless @options.foreground?

      while running 
        sleep 5
        response = get "#{API_URL}/#{run_id}?gem_version=#{Rainforest::Cli::VERSION}"
        if %w(queued in_progress sending_webhook waiting_for_callback).include?(response["state"])
          puts "Run #{run_id} is #{response['state']} and is #{response['current_progress']['percent']}% complete"
          running = false if response["result"] == 'failed' && @options.failfast?
        else
          puts "Run #{run_id} is now #{response["state"]} and has #{response["result"]}"
          running = false
        end
      end

      if response["result"] != "passed"
        exit 1
      end
    end

    def self.post(url, body = {})
      response = HTTParty.post url, {
        body: body, 
        headers: {"CLIENT_TOKEN" => @options.token}
      }

      JSON.parse(response.body)
    end

    def self.get(url, body = {})
      response = HTTParty.get url, {
        body: body, 
        headers: {"CLIENT_TOKEN" => @options.token}
      }

      JSON.parse(response.body)
    end
  end
end
