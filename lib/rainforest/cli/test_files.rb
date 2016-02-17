# frozen_string_literal: true
class RainforestCli::TestFiles
  DEFAULT_TEST_FOLDER = './spec/rainforest'.freeze
  EXT = '.rfml'.freeze

  attr_reader :test_folder, :test_paths, :test_data

  def initialize(test_folder = nil)
    @test_folder = test_folder || DEFAULT_TEST_FOLDER
    create_test_folder unless Dir.exist?(@test_folder)

    @test_paths = "#{@test_folder}/**/*#{EXT}"
    @test_data = [].tap do |all_tests|
      Dir.glob(@test_paths) do |file_name|
        all_tests << RainforestCli::TestParser::Parser.new(File.read(file_name)).process
      end
    end
  end

  def file_extension
    EXT
  end

  def rfml_ids
    @test_data.map(&:rfml_id)
  end

  def count
    @test_data.count
  end

  private

  def create_test_folder
    test_folder_string = test_folder.dup

    # preserve leading slash for absolute paths
    leading_slash = ''
    if test_folder_string[0] == '/'
      leading_slash = '/'
      test_folder_string = test_folder_string[1..-1]
    end

    # Remove leading `./` for simplicity
    test_folder_string = test_folder_string[0..1] == './' ? test_folder_string[2..-1] : test_folder_string

    folders = test_folder_string.split('/')
    folders.first.prepend(leading_slash)

    (0...folders.length).each do |idx|
      dir_name = folders[0..idx].join('/')
      Dir.mkdir(dir_name) unless Dir.exist?(dir_name)
    end
  end
end
