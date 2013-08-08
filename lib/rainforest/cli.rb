require "rainforest/cli/version"
require "rainforest/cli/options"
require 'httparty'
require 'json'

module Rainforest
  module Cli 
    API_URL = 'https://app.rainforestqa.com/api/1/runs'

    def self.start(args)
      @options = OptionParser.new(args)

      response = post(API_URL, tests: @options.tests)
      run_id = response["id"]
      running = true

      while running
        sleep 5
        response = get "#{API_URL}/#{run_id}"
        if %w(queued in_progress sending_webhook waiting_for_callback).include?(response["state"])
          puts "Run #{run_id} is still running"
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
