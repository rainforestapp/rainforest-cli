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

      response = wrap_exceptions(options[:retries_on_failures]) do
        HTTParty.send(method, url, { body: body, headers: headers, verify: false })
      end

      if response.code.between?(200, 299)
        JSON.parse(response.body)
      elsif options[:attempts].to_i == 0
        logger.fatal "Non 200 code received for request to #{url}"
        logger.debug response.body if options[:debug]
        exit 1
      else
        logger.warn("HTTP request was unsuccessful. URL: #{url}. Status: #{response.code}")
        logger.warn("Retrying HTTP request #{remaining_attempts} more times")
        options[:attempts] -= 1 unless options[:attempts].nil?
        request(method, path, body, options)
      end
    end

    def api_token_set?
      !@token.nil?
    end

    private

    def wrap_exceptions(retries_on_failures)
      @retry_delay = 0
      @waiting_on_retries = false
      loop do
        begin
          # Suspend tries until wait period is over
          if @waiting_on_retries
            Kernel.sleep 5
          else
            return Http::Exceptions.wrap_exception { yield }
          end
        rescue Http::Exceptions::HttpException, Timeout::Error => e
          raise e unless retries_on_failures

          unless @waiting_on_retries
            @waiting_on_retries = true
            @retry_delay += RETRY_INTERVAL

            RainforestCli.logger.warn 'Exception Encountered while trying to contact Rainforest API:'
            RainforestCli.logger.warn "\t\t#{e.message}"
            RainforestCli.logger.warn "Retrying again in #{@retry_delay} seconds..."

            Kernel.sleep @retry_delay
            @waiting_on_retries = false
          end
        end
      end
    end

    def make_url(url)
      File.join(API_URL, url)
    end

    def headers
      {
        'CLIENT_TOKEN' => @token,
        'User-Agent' => "Rainforest-cli-#{RainforestCli::VERSION}",
      }
    end

    def logger
      RainforestCli.logger
    end
  end
end
