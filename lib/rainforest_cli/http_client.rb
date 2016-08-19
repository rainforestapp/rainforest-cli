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

    def delete(url, body = {})
      response = HTTParty.delete make_url(url), {
        body: body,
        headers: headers,
        verify: false,
      }

      JSON.parse(response.body)
    end

    def post(url, body = {}, options = {})
      wrap_exceptions(options[:retries_on_failures]) do
        response = HTTParty.post make_url(url), {
          body: body,
          headers: headers,
          verify: false,
        }

        return JSON.parse(response.body)
      end
    end

    def get(url, body = {}, options = {})
      wrap_exceptions(options[:retries_on_failures]) do
        response = HTTParty.get make_url(url), {
          body: body,
          headers: headers,
          verify: false,
        }

        if response.code == 200
          return JSON.parse(response.body)
        else
          RainforestCli.logger.warn("Status Code: #{response.code}, #{response.body}")
          return nil
        end
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
            Http::Exceptions.wrap_exception { yield }
            break
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
  end
end
