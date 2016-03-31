# frozen_string_literal: true
require 'securerandom'
require 'rainforest'
require 'parallel'
require 'ruby-progressbar'

class RainforestCli::Exporter
  attr_reader :options, :client, :test_files

  def initialize(options)
    @options = options
    ::Rainforest.api_key = @options.token
    @test_files = RainforestCli::TestFiles.new(@options)
  end

  def logger
    RainforestCli.logger
  end

  def threads
    RainforestCli::THREADS
  end

  def export
    tests = Rainforest::Test.all(page_size: 1000)
    p = ProgressBar.create(title: 'Rows', total: tests.count, format: '%a %B %p%% %t')
    Parallel.each(tests, in_threads: threads, finish: lambda { |_item, _i, _result| p.increment }) do |test|

      # File name
      file_name = sprintf('%010d', test.id) + '_' + test.title.strip.gsub(/[^a-z0-9 ]+/i, '').gsub(/ +/, '_').downcase
      file_name = test_files.create_file(file_name)
      File.truncate(file_name, 0)

      # Get the full test from the API
      test = Rainforest::Test.retrieve(test.id)

      File.open(file_name, 'a') do |file|
        file.puts(get_header(test))

        test.elements.each_with_index do |element, index|
          process_element(file, element, index)
        end
      end
    end
  end

  private

  def process_element file, element, index
    case element[:type]
    when 'test'
      if @options.embed_tests
        file.puts '' unless index == 0
        file.puts "- #{element[:element][:rfml_id]}"
      else
        element[:element][:elements].each do |sub_element|
          index = process_element(file, sub_element, index)
        end
      end
    when 'step'
      file.puts '' unless index == 0
      file.puts "# step #{index + 1}" if @options.debug
      file.puts element[:element][:action]
      file.puts element[:element][:response]
    else
      raise "Unknown element type: #{element[:type]}"
    end
  end

  def get_header(test)
    <<-EOF
#! #{test.rfml_id}
# title: #{test.title}
# start_uri: #{test.start_uri}
# tags: #{test.tags.join(", ")}
# browsers: #{test.browsers.join(", ")}
#

    EOF
  end
end
