# frozen_string_literal: true
require 'securerandom'
require 'rainforest'
require 'parallel'
require 'ruby-progressbar'

class RainforestCli::TestImporter
  attr_reader :options, :client, :test_files

  SAMPLE_FILE = <<EOF
#! %s (Test ID - only edit if this test has not yet been uploaded)
# title: New test
# start_uri: /
#
# Lines starting with # are test attributes or comments
# Possible attributes: #{RainforestCli::TestParser::Parser::TEXT_FIELDS.join(', ')}
#
# Steps are composed of two lines: an action and a question. Example:
#
# This is the step action.
# This is the step question?
#

EOF

  def initialize(options)
    @options = options
    ::Rainforest.api_key = @options.token
    @test_files = RainforestCli::TestFiles.new(@options.test_folder)
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
      file_name = create_new(file_name)
      File.truncate(file_name, 0)

      # Get the full test from the API
      test = Rainforest::Test.retrieve(test.id)

      File.open(file_name, 'a') do |file|
        file.puts _get_header(test)

        index = 0
        test.elements.each do |element|
          index = _process_element(file, element, index)
        end
      end
    end
  end

  def _process_element file, element, index
    case element[:type]
    when 'test'
      element[:element][:elements].each do |sub_element|
        index = _process_element(file, sub_element, index)
      end
    when 'step'
      file.puts '' unless index == 0
      file.puts "# step #{index + 1}" if @options.debug
      file.puts element[:element][:action]
      file.puts element[:element][:response]
    else
      raise "Unknown element type: #{element[:type]}"
    end

    index += 1
    index
  end

  # add comments if not already present
  def _get_header test
    out = []

    has_id = false
    test.description.to_s.strip.lines.map(&:chomp).each_with_index do |line, _line_no|
      line = line.gsub(/\#+$/, '').strip

      # make sure the test has an ID
      has_id = true if line[0] == '!'

      out << '#' + line
    end

    unless has_id
      browsers = test.browsers.map {|b| b[:name] if b[:state] == 'enabled' }.compact
      out = [
        "#! #{SecureRandom.uuid}",
        "# title: #{test.title}",
        "# start_uri: #{test.start_uri}",
        "# tags: #{test.tags.join(", ")}",
        "# browsers: #{browsers.join(", ")}",
        '#',
        ' ',
      ] + out
    end

    out.compact.join("\n")
  end

  def _get_id test
    id = nil
    test.description.to_s.strip.lines.map(&:chomp).each_with_index do |line, _line_no|
      line = line.gsub(/\#+$/, '').strip
      if line[0] == '!'
        id = line[1..-1].split(' ').first
        break
      end
    end
    id
  end

  def validate
    tests = {}
    has_errors = []

    Dir.glob(test_files.test_paths).each do |file_name|
      out = RainforestCli::TestParser::Parser.new(file_name).process

      tests[file_name] = out
      has_errors << file_name if out.errors != {}
    end

    if !has_errors.empty?
      logger.error 'Parsing errors:'
      logger.error ''
      has_errors.each do |file_name|
        logger.error ' ' + file_name
        tests[file_name].errors.each do |_line, error|
          logger.error "\t#{error}"
        end
      end

      exit 2
    end

    if @options.debug
      tests.each do |file_name, test|
        logger.debug test.inspect
        logger.debug "#{file_name}"
        logger.debug test.description
        test.steps.each do |step|
          logger.debug "\t#{step}"
        end
      end
    else
      logger.info '[VALID]'
    end

    return tests
  end

  def create_new file_name = nil
    name = @options.file_name if @options.file_name
    name = file_name if !file_name.nil?
    ext = test_files.file_extension

    uuid = SecureRandom.uuid
    name = "#{uuid}#{ext}" unless name
    name += ext unless name[-ext.length..-1] == ext

    FileUtils.mkdir_p(test_files.test_folder) unless Dir.exist?(test_files.test_folder)
    name = File.join([test_files.test_folder, name])

    File.open(name, 'w') { |file| file.write(sprintf(SAMPLE_FILE, uuid)) }

    logger.info "Created #{name}" if file_name.nil?
    name
  end
end
