require 'securerandom'
require 'rainforest'
require 'parallel'
require 'ruby-progressbar'

class RainforestCli::TestImporter
  attr_reader :options, :client, :test_files
  THREADS = 32.freeze

  SAMPLE_FILE = <<EOF
#! %s (this is the ID, don't edit it)
# title: New test
#
# 1. steps:
#   a) pairs of lines are steps (first line = action, second = response)
#   b) second line must have a ?
#   c) second line must not be blank
# 2. comments:
#   a) lines starting # are comments
#

EOF

  def initialize(options)
    # FIXME: Temporarily switching the api to local. Do not keep it this way!
    ::Rainforest.api_base = 'http://app.rainforest.dev/api/1'

    @options = options
    ::Rainforest.api_key = @options.token
    @test_files = RainforestCli::TestFiles.new(@options.test_spec_folder)
  end

  def logger
    RainforestCli.logger
  end

  def export
    tests = Rainforest::Test.all(page_size: 1000)
    p = ProgressBar.create(title: 'Rows', total: tests.count, format: '%a %B %p%% %t')
    Parallel.each(tests, in_threads: THREADS, finish: lambda { |item, i, result| p.increment }) do |test|

      # File name
      file_name = sprintf('%010d', test.id) + "_" + test.title.strip.gsub(/[^a-z0-9 ]+/i, '').gsub(/ +/, '_').downcase
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
      file.puts "" unless index == 0
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
    test.description.to_s.strip.lines.map(&:chomp).each_with_index do |line, line_no|
      line = line.gsub(/\#+$/, '').strip

      # make sure the test has an ID
      has_id = true if line[0] == "!"

      out << "#" + line
    end

    unless has_id
      browsers = test.browsers.map {|b| b[:name] if b[:state] == "enabled" }.compact
      out = ["#! #{SecureRandom.uuid}", "# title: #{test.title}", "# start_uri: #{test.start_uri}", "# tags: #{test.tags.join(", ")}", "# browsers: #{browsers.join(", ")}", "#", " "] + out
    end

    out.compact.join("\n")
  end

  def _get_id test
    id = nil
    test.description.to_s.strip.lines.map(&:chomp).each_with_index do |line, line_no|
      line = line.gsub(/\#+$/, '').strip
      if line[0] == "!"
        id = line[1..-1].split(' ').first
        break
      end
    end
    id
  end

  def upload
    logger.info "Uploading tests..."
    p = ProgressBar.create(title: 'Rows', total: test_files.count, format: '%a %B %p%% %t')

    Parallel.each(test_files.test_data, in_threads: THREADS, finish: lambda { |item, i, result| p.increment }) do |rfml_test|
      next unless rfml_test.steps.count > 0

      if @options.debug
        logger.debug "Starting: #{rfml_test.rfml_id}"
        logger.debug "\t#{rfml_test.start_uri || "/"}"
      end

      test_obj = create_test_obj(rfml_test)

      # Upload the test
      begin
        if rfml_id_mappings[rfml_test.rfml_id]
          t = Rainforest::Test.update(rfml_id_mappings[rfml_test.rfml_id], test_obj)

          logger.info "\tUpdated #{rfml_test.id} -- ##{t.id}" if @options.debug
        else
          t = Rainforest::Test.create(test_obj)

          logger.info "\tCreated #{rfml_test.id} -- ##{t.id}" if @options.debug
        end
      rescue => e
        logger.fatal "Error: #{rfml_test.rfml_id}: #{e}"
        exit 2
      end
    end
  end

  def validate
    tests = {}
    has_errors = []

    Dir.glob(test_files.test_paths).each do |file_name|
      out = RainforestCli::TestParser::Parser.new(File.read(file_name)).process

      tests[file_name] = out
      has_errors << file_name if out.errors != {}
    end

    if !has_errors.empty?
      logger.error "Parsing errors:"
      logger.error ""
      has_errors.each do |file_name|
        logger.error " " + file_name
        tests[file_name].errors.each do |line, error|
          logger.error "\t#{error.to_s}"
        end
      end

      exit 2
    end

    if @options.debug
      tests.each do |file_name,test|
        logger.debug test.inspect
        logger.debug "#{file_name}"
        logger.debug test.description
        test.steps.each do |step|
          logger.debug "\t#{step}"
        end
      end
    else
      logger.info "[VALID]"
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

    File.open(name, "w") { |file| file.write(sprintf(SAMPLE_FILE, uuid)) }

    logger.info "Created #{name}" if file_name.nil?
    name
  end

  private

  def create_rfml_id_mappings
    @_id_mappings ||= {}.tap do |id_mappings|
      logger.info "Syncing tests"
      Rainforest::Test.all(page_size: 1000, rfml_ids: test_files.rfml_ids).each do |rf_test|
        rfml_id = rf_test.rfml_id
        next if rfml_id.nil?

        id_mappings[rfml_id] = rf_test.id
      end
    end
  end

  def create_test_obj(rfml_test)
    test_obj = {
      start_uri: rfml_test.start_uri || "/",
      title: rfml_test.title,
      description: rfml_test.description,
      tags: (["ro"] + rfml_test.tags).uniq,
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
              id: rf_ids[step.rfml_id]
            }
          }
        end
      end
    }

    unless rfml_test.browsers.empty?
      test_obj[:browsers] = rfml_test.browsers.map {|b|
        {'state' => 'enabled', 'name' => b}
      }
    end

    test_obj
  end
end
