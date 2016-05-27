# frozen_string_literal: true
class RainforestCli::TestFiles
  DEFAULT_TEST_FOLDER = './spec/rainforest'
  FILE_EXTENSION = '.rfml'
  SAMPLE_FILE = <<EOF
#! %s
# title: New test
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
    if @test_data.nil?
      @test_data = []
      if Dir.exist?(@test_folder)
        Dir.glob(test_paths) do |file_name|
          @test_data << RainforestCli::TestParser::Parser.new(file_name).process
        end
      end
    end
    @test_data
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

    uuid = SecureRandom.uuid

    name = file_name || uuid.to_s
    name += file_extension unless name[-file_extension.length..-1] == file_extension
    name = File.join(test_folder, name)

    File.open(name, 'w') { |file| file.write(sprintf(SAMPLE_FILE, uuid)) }

    logger.info "Created #{name}"
    name
  end

  private

  def logger
    RainforestCli.logger
  end
end
