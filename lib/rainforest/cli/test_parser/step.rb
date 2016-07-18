# frozen_string_literal: true
class RainforestCli::TestParser::Step < Struct.new(:action, :response, :redirect)
  UPLOADABLE_REGEX = /{{ ?file\.(download|screenshot)\((.+)\)+ ?}}/

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
        response: response
      }
    }
  end

  def has_uploadable?
    !!uploadable_in_action || !!uploadable_in_response
  end

  def uploadable_in_action
    action.match(UPLOADABLE_REGEX)
  end

  def uploadable_in_response
    response.match(UPLOADABLE_REGEX)
  end
end
