# frozen_string_literal: true
class RainforestCli::TestParser::EmbeddedTest < Struct.new(:rfml_id, :redirect)
  def type
    :test
  end

  def to_s
    "--> embed: #{rfml_id}"
  end

  def redirection
    redirect || 'true'
  end

  def to_element(primary_key_id)
    {
      type: 'test',
      redirection: redirection,
      element: {
        id: primary_key_id,
      },
    }
  end

  def has_uploadable?
    false
  end
end
