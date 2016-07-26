# frozen_string_literal: true
class RainforestCli::TestParser::Step < Struct.new(:action, :response, :redirect)
  UPLOADABLE_REGEX = /{{ ?file\.(download|screenshot)\(([^\)]+)\) ?}}/

  def type
    :step
  end

  def redirection
    redirect || 'true'
  end

  def to_s
    "#{action} --> #{response}"
  end

  def to_element
    {
      type: 'step',
      redirection: redirection,
      element: {
        action: action,
        response: response,
      },
    }
  end

  def has_uploadable?
    uploadable_in_action.any? || uploadable_in_response.any?
  end

  def uploadable_in_action
    action.scan(UPLOADABLE_REGEX)
  end

  def uploadable_in_response
    response.scan(UPLOADABLE_REGEX)
  end
end
