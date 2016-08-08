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

  def to_json
    {
      start_uri: start_uri || '/',
      title: title,
      site_id: site_id,
      description: description,
      source: 'rainforest-cli',
      tags: tags.uniq,
      rfml_id: rfml_id,
      browsers: browser_json,
    }
  end

  def browser_json
    browsers.map do |b|
      {'state' => 'enabled', 'name' => b}
    end
  end

  def has_uploadable_files?
    steps.any?(&:has_uploadable_files?)
  end
end
