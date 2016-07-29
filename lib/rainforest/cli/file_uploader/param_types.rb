# frozen_string_literal: true

require 'mimemagic'

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
    file.pos = 0
    content_type = MimeMagic.by_magic(file) || 'text/plain'
    file.pos = 0
    <<-EOS
Content-Disposition: form-data; name="#{param_name}"; filename="#{File.basename(file.path)}"\r
Content-Type: #{content_type}\r
Content-Transfer-Encoding: binary\r
\r
#{file.read.strip}\r
    EOS
  end
end
