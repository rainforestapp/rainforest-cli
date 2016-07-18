# frozen_string_literal: true
class RainforestCli::TestParser::Test < Struct.new(
  :file_name,
  :rfml_id,
  :description,
  :title,
  :start_uri,
  :site_id,
  :steps,
  :errors,
  :tags,
  :browsers
)
  def embedded_ids
    steps.inject([]) { |embeds, step| step.type == :test ? embeds + [step.rfml_id] : embeds }
  end

  def has_uploadable?
    steps.any?(&:has_uploadable?)
  end
end
