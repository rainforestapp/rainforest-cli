# frozen_string_literal: true
class RainforestCli::TestParser::Step < Struct.new(:action, :response, :redirect)
  UPLOADABLE_REGEX = /{{ *file\.(download|screenshot)\(([^\)]+)\) *}}/

  def type
    :step
  end

  def to_s
    "#{action} --> #{response}"
  end

  def has_uploadable_files?
    uploadable_in_action.any? || uploadable_in_response.any?
  end

  def uploadable_in_action
    action.scan(UPLOADABLE_REGEX).select do |match|
      needs_parameterization?(match)
    end
  end

  def uploadable_in_response
    response.scan(UPLOADABLE_REGEX).select do |match|
      needs_parameterization?(match)
    end
  end

  private

  def needs_parameterization?(match)
    argument = match[1]
    parameters = argument.split(',').map(&:strip)
    if parameters.length >= 2
      has_file_id = parameters[0].to_i > 0
      has_file_sig = parameters[1].length == 6
      !(has_file_id && has_file_sig)
    else
      true
    end
  end
end
