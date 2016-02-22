# frozen_string_literal: true
require 'securerandom'
require 'rainforest'
require 'parallel'
require 'ruby-progressbar'

class RainforestCli::TestImporter
  attr_reader :options, :client, :test_files
  THREADS = 32

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

  def export
    tests = Rainforest::Test.all(page_size: 1000)
    p = ProgressBar.create(title: 'Rows', total: tests.count, format: '%a %B %p%% %t')
    Parallel.each(tests, in_threads: THREADS, finish: lambda { |_item, _i, _result| p.increment }) do |test|

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

  def upload
    upload_groups = make_test_priority_groups

    logger.info 'Uploading tests...'

    # Upload in parallel if order doesn't matter
    if upload_groups.count > 1
      upload_groups_sequentially(upload_groups)
    else
      upload_group_in_parallel(upload_groups.first)
    end
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
    name = File.join([@test_files.test_folder, name])

    File.open(name, 'w') { |file| file.write(sprintf(SAMPLE_FILE, uuid)) }

    logger.info "Created #{name}" if file_name.nil?
    name
  end

  private

  def upload_groups_sequentially(upload_groups)
    progress_bar = ProgressBar.create(title: 'Rows', total: test_files.count, format: '%a %B %p%% %t')
    upload_groups.each_with_index do |rfml_tests, idx|
      if idx == (rfml_tests.length - 1)
        upload_group_in_parallel(rfml_tests, progress_bar)
      else
        rfml_tests.each { |rfml_test| upload_test(rfml_test) }
        progress_bar.increment
      end
    end
  end

  def upload_group_in_parallel(rfml_tests, progress_bar = nil)
    progress_bar ||= ProgressBar.create(title: 'Rows', total: rfml_tests.count, format: '%a %B %p%% %t')
    Parallel.each(
      rfml_tests,
      in_threads: THREADS,
      finish: lambda { |_item, _i, _result| progress_bar.increment }
    ) do |rfml_test|
      upload_test(rfml_test)
    end
  end

  def upload_test(rfml_test)
    return unless rfml_test.steps.count > 0

    if @options.debug
      logger.debug "Starting: #{rfml_test.rfml_id}"
      logger.debug "\t#{rfml_test.start_uri || "/"}"
    end

    test_obj = create_test_obj(rfml_test)
    # Upload the test
    begin
      if rfml_id_mappings[rfml_test.rfml_id]
        t = Rainforest::Test.update(rfml_id_mappings[rfml_test.rfml_id], test_obj)

        logger.info "\tUpdated #{rfml_test.rfml_id} -- ##{t.id}" if @options.debug
      else
        t = Rainforest::Test.create(test_obj)

        logger.info "\tCreated #{rfml_test.rfml_id} -- ##{t.id}" if @options.debug
        rfml_id_mappings[rfml_test.rfml_id] = t.id
      end
    rescue => e
      logger.fatal "Error: #{rfml_test.rfml_id}: #{e}"
      exit 2
    end
  end

  def rfml_id_mappings
    if @_id_mappings.nil?
      @_id_mappings = {}.tap do |id_mappings|
        Rainforest::Test.all(page_size: 1000, rfml_ids: test_files.rfml_ids).each do |rf_test|
          rfml_id = rf_test.rfml_id
          next if rfml_id.nil?

          id_mappings[rfml_id] = rf_test.id
        end
      end
    end
    @_id_mappings
  end

  def create_test_obj(rfml_test)
    test_obj = {
      start_uri: rfml_test.start_uri || '/',
      title: rfml_test.title,
      description: rfml_test.description,
      tags: (['ro'] + rfml_test.tags).uniq,
      rfml_id: rfml_test.rfml_id,
      elements: rfml_test.steps.map do |step|
        case step.type
        when :step
          {
            type: 'step',
            redirection: true,
            element: {
              action: step.action,
              response: step.response
            }
          }
        when :test
          {
            type: 'test',
            redirection: true,
            element: {
              id: rfml_id_mappings[step.rfml_id]
            }
          }
        end
      end
    }

    unless rfml_test.browsers.empty?
      test_obj[:browsers] = rfml_test.browsers.map do|b|
        {'state' => 'enabled', 'name' => b}
      end
    end

    test_obj
  end

  def make_test_priority_groups
    # Prioritize embedded tests before other tests
    upload_groups = []
    unordered_tests = []
    queued_tests = test_files.test_data.dup

    until queued_tests.empty?
      new_ordered_group = []
      ordered_ids = upload_groups.flatten.map(&:rfml_id)

      queued_tests.each do |rfml_test|
        if (rfml_test.embedded_ids - ordered_ids).empty?
          new_ordered_group << rfml_test
        else
          unordered_tests << rfml_test
        end
      end

      # If all the queued tests make it to the unordered tests group, then
      # they contain non-existent RFML ids.
      if queued_tests.length == unordered_tests.length
        misconfigured_tests = filter_misconfigured_tests(queued_tests)
        raise TestNotFound.new(misconfigured_tests.map(&:file_name))
      end

      upload_groups << new_ordered_group
      queued_tests = unordered_tests
      unordered_tests = []
    end

    upload_groups
  end

  # Filter out tests that depend on the actual misconfigured tests
  def filter_misconfigured_tests(unfiltered_tests)
    all_ids = unfiltered_tests.map(&:rfml_id)
    unfiltered_tests.reject { |test| (test.embedded_ids - all_ids).empty? }
  end

  class TestNotFound < RuntimeError
    def initialize(file_names)
      super("The following tests contain embedded tests not found in test directory:\n\t#{file_names.join("\n\t")}\n\n")
    end
  end
end
