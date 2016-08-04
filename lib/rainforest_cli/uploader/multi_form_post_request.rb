# frozen_string_literal: true
require 'httparty'

class RainforestCli::Uploader::MultiFormPostRequest
  BOUNDARY = 'RainforestCli'

  class Param < Struct.new(:param_name, :value)
    def to_multipart
      <<-EOS
Content-Disposition: form-data; name="#{param_name}"\r
\r
#{value}\r
      EOS
    end
  end

  class FileParam < Struct.new(:param_name, :file)
    def to_multipart
      <<-EOS
Content-Disposition: form-data; name="#{param_name}"; filename="#{File.basename(file.path)}"\r
Content-Type: #{MimeMagic.by_path(file.path)}\r
Content-Transfer-Encoding: binary\r
\r
#{file.read.strip}\r
      EOS
    end
  end

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
