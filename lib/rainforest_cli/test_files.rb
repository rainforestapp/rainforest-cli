# frozen_string_literal: true
require 'securerandom'

class RainforestCli::TestFiles
  DEFAULT_TEST_FOLDER = './spec/rainforest'
  FILE_EXTENSION = '.rfml'
  SAMPLE_FILE = <<EOF
#! %s
# title: %s
# start_uri: /
# tags: rfml-test
#

This is a step action.
This is a step question?

This is another step action.
This is another step question?

EOF

  attr_reader :test_folder, :test_data

  def initialize(options)
    @options = options

    if @options.command == 'rm'
      @test_folder = File.dirname(@options.file_name)
    elsif @options.test_folder.nil?
      logger.info "No test folder supplied. Using default folder: #{DEFAULT_TEST_FOLDER}"
      @test_folder = File.expand_path(DEFAULT_TEST_FOLDER)
    else
      @test_folder = File.expand_path(@options.test_folder)
    end
  end

  def test_paths
    "#{@test_folder}/**/*#{FILE_EXTENSION}"
  end

  def test_data
    @test_data ||= get_test_data
  end

  def file_extension
    FILE_EXTENSION
  end

  def rfml_ids
    test_data.map(&:rfml_id)
  end

  def count
    test_data.count
  end

  def test_dictionary
    {}.tap do |dictionary|
      test_data.each { |rfml_test| dictionary[rfml_test.rfml_id] = rfml_test }
    end
  end

  def ensure_directory_exists
    FileUtils.mkdir_p(test_folder) unless Dir.exist?(test_folder)
  end

  def create_file(file_name = @options.file_name)
    ensure_directory_exists

    title = file_name || 'Unnamed Test'
    file_path = title.dup

    if title[-file_extension.length..-1] == file_extension
      title = title[0...-file_extension.length]
    else
      file_path += file_extension
    end

    file_path = unique_path(File.join(test_folder, file_path))

    File.open(file_path, 'w') { |file| file.write(sprintf(SAMPLE_FILE, SecureRandom.uuid, title)) }

    logger.info "Created #{file_path}"
    file_path
  end

  private

  def get_test_data
    data = []
    if Dir.exist?(@test_folder)
      Dir.glob(test_paths) do |file_name|
        data << RainforestCli::TestParser::Parser.new(file_name).process
      end
    end
    filter_tests(data)
  end

  def filter_tests(tests)
    tests.select do |test|
      pass_tag_filter = (@options.tags - test.tags).empty?
      pass_site_filter = @options.site_id.nil? || @options.site_id == test.site_id
      pass_tag_filter || pass_site_filter
    end
  end

  def unique_path(file_path)
    path = file_path[0...-file_extension.length]
    identifier = 0

    loop do
      id_string = (identifier > 0) ? " (#{identifier})" : ''
      test_path = path + id_string + file_extension

      return test_path unless File.exist?(test_path)
      identifier += 1
    end
  end

  def logger
    RainforestCli.logger
  end
end
