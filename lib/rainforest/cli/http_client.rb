module Rainforest
  module Cli
    class HttpClient
      API_URL = ENV.fetch("RAINFOREST_API_URL") do
        'https://app.rainforestqa.com/api/1'
      end.freeze

      def initialize(token:)
        @token = token
      end

      def delete(url, body = {})
        response = HTTParty.delete make_url(url), {
          body: body,
          headers: headers,
        }

        JSON.parse(response.body)
      end

      def post(url, body = {})
        response = HTTParty.post make_url(url), {
          body: body,
          headers: headers,
        }

        JSON.parse(response.body)
      end

      def get(url, body = {})
        response = HTTParty.get make_url(url), {
          body: body,
          headers: headers,
        }

        if response.code == 200
          JSON.parse(response.body)
        else
          nil
        end
      end

      private
      def make_url(url)
        File.join(API_URL, url)
      end

      def headers
        {
          "CLIENT_TOKEN" => @token,
          "User-Agent" => "Rainforest-cli-#{Rainforest::Cli::VERSION}"
        }
      end
    end
  end
end
