# frozen_string_literal: true
require 'securerandom'
require 'rainforest'
require 'parallel'
require 'ruby-progressbar'

class RainforestCli::Exporter
  attr_reader :options, :client, :test_files

  def initialize(options)
    @options = options
    @test_files = RainforestCli::TestFiles.new(@options)
    @remote_tests = RainforestCli::RemoteTests.new(@options)
  end

  def logger
    RainforestCli.logger
  end

  def threads
    RainforestCli::THREADS
  end

  def export
    test_ids =
      if @options.tests.length > 0
        @options.tests
      else
        @remote_tests.primary_ids
      end

    p = ProgressBar.create(title: 'Tests', total: test_ids.count, format: '%a %B %p%% %t')
    Parallel.each(test_ids, in_threads: threads, finish: lambda { |_item, _i, _result| p.increment }) do |test_id|
      # Get the full test from the API
      test = http_client.get("/tests/#{test_id}")

      # File name
      file_name = sprintf('%010d', test['id']) + '_' + test['title'].strip.gsub(/[^a-z0-9 ]+/i, '').gsub(/ +/, '_').downcase
      file_name = test_files.create_file(file_name)
      File.truncate(file_name, 0)

      File.open(file_name, 'a') do |file|
        file.puts(get_header(test))

        first_step_processed = false
        test['elements'].each_with_index do |element, index|
          first_step_processed = process_element(file, element, index, first_step_processed)
        end
      end
    end
  end

  private

  def process_element(file, element, index, first_step_processed)
    case element['type']
    when 'test'
      if @options.embed_tests
        file.puts '' unless index == 0
        # no redirect if an embedded test is the first step
        file.puts "# redirect: #{element['redirection']}" if index > 0
        file.puts "- #{element['element']['rfml_id']}"
      else
        element['element']['elements'].each_with_index do |sub_element, i|
          # no redirect flags for flattened tests
          process_element(file, sub_element, i + index, true)
        end
      end
    when 'step'
      file.puts '' unless index == 0

      # add redirect for first step if preceded by an embedded test
      if index > 0 && first_step_processed == false
        file.puts "# redirect: #{element['redirection']}"
      end
      file.puts element['element']['action'].gsub("\n", ' ').strip
      file.puts element['element']['response'].gsub("\n", ' ').strip
      first_step_processed = true
    else
      raise "Unknown element type: #{element['type']}"
    end
    first_step_processed
  end

  def get_header(test)
    browsers = test['browsers'].map { |b| b['name'] if b['state'] == 'enabled' }.compact
    <<-EOF
#! #{test['rfml_id']}
# title: #{test['title']}
# start_uri: #{test['start_uri']}
# tags: #{test['tags'].join(", ")}
# browsers: #{browsers.join(", ")}
#

    EOF
  end

  def http_client
    RainforestCli.http_client
  end
end
