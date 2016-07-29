# frozen_string_literal: true

require 'httparty'

require 'rainforest/cli/file_uploader/param_types'

class RainforestCli::FileUploader::MultiFormPostRequest
  BOUNDARY = 'RainforestCli'

  class << self
    def request(url, params)
      HTTParty.post(url, body: make_body(params), headers: headers)
    end

    def make_body(params)
      fp = []
      params.each do |k, v|
        if v.respond_to?(:read)
          fp.push(FileParam.new(k, v))
        else
          fp.push(Param.new(k, v))
        end
      end
      fp.map { |p| "--#{BOUNDARY}\n#{p.to_multipart}" }.join + "--#{BOUNDARY}--"
    end

    def headers
      { 'Content-type' => "multipart/form-data, boundary=#{BOUNDARY}" }
    end
  end
end
