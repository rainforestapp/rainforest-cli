# frozen_string_literal: true
require 'http/exceptions'

module RainforestCli
  class HttpClient
    API_URL = ENV.fetch('RAINFOREST_API_URL') do
      'https://app.rainforestqa.com/api/1'
    end.freeze

    def initialize(options)
      @token = options.fetch(:token)
    end

    def delete(url, body = {})
      response = HTTParty.delete make_url(url), {
        body: body,
        headers: headers,
        verify: false
      }

      JSON.parse(response.body)
    end

    def post(url, body = {})
      response = HTTParty.post make_url(url), {
        body: body,
        headers: headers,
        verify: false
      }

      JSON.parse(response.body)
    end

    def get(url, body = {}, max_exceptions: 0)
      with_exception_tolerance(max_exceptions) do
        response = HTTParty.get make_url(url), {
          body: body,
          headers: headers,
          verify: false
        }

        if response.code == 200
          return JSON.parse(response.body)
        else
          return nil
        end
      end
    end

    private
    def with_exception_tolerance(allowed_exceptions = 0)
      loop do
        begin
          Http::Exceptions.wrap_exception { yield }
          break
        rescue Http::Exceptions::HttpException, Timeout::Error => e
          # Give up on final attempt
          raise e if allowed_exceptions <= 0
          allowed_exceptions -= 1
        end
      end
    end

    def make_url(url)
      File.join(API_URL, url)
    end

    def headers
      {
        'CLIENT_TOKEN' => @token,
        'User-Agent' => "Rainforest-cli-#{RainforestCli::VERSION}"
      }
    end
  end
end
