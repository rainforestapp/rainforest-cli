# frozen_string_literal: true
require 'httparty'
require 'http/exceptions'

module RainforestCli
  class HttpClient
    API_URL = ENV.fetch('RAINFOREST_API_URL') do
      'https://app.rainforestqa.com/api/1'
    end.freeze
    RETRY_INTERVAL = 10

    def initialize(options)
      @token = options.fetch(:token)
    end

    def delete(path, body = {}, options = {})
      request(:delete, path, body, options)
    end

    def post(path, body = {}, options = {})
      request(:post, path, body, options)
    end

    def get(path, body = {}, options = {})
      request(:get, path, body, options)
    end

    def request(method, path, body, options)
      url = File.join(API_URL, path)

      loop do
        begin
          response = Http::Exceptions.wrap_exception do
            HTTParty.send(method, url, { body: body, headers: headers, verify: false })
          end

          if response.code.between?(200, 299)
            return JSON.parse(response.body)
          elsif options[:retries_on_failures] && response.code >= 500
            delay = retry_delay
            logger.warn("HTTP request was unsuccessful. URL: #{url}. Status: #{response.code}")
            logger.warn "Retrying again in #{delay} seconds..."
            Kernel.sleep delay
          else
            logger.fatal "Non 200 code received for request to #{url}"
            logger.fatal "Server response: #{response.body}"
            exit 1
          end
        rescue Http::Exceptions::HttpException, Timeout::Error => e
          raise e unless options[:retries_on_failures]

          delay = retry_delay
          logger.warn 'Exception Encountered while trying to contact Rainforest API:'
          logger.warn "\t\t#{e.message}"
          logger.warn "Retrying again in #{delay} seconds..."

          Kernel.sleep delay
        end
      end
    end

    def api_token_set?
      !@token.nil?
    end

    private

    def make_url(url)
      File.join(API_URL, url)
    end

    def headers
      {
        'CLIENT_TOKEN' => @token,
        'User-Agent' => "Rainforest-cli-#{RainforestCli::VERSION}",
      }
    end

    def retry_delay
      # make retry delay random to desynchronize multiple threads
      RETRY_INTERVAL + rand(RETRY_INTERVAL)
    end

    def logger
      RainforestCli.logger
    end
  end
end
