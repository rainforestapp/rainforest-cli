class RainforestCli::TestFiles
  DEFAULT_TEST_FOLDER = './spec/rainforest'.freeze
  EXT = ".rfml".freeze

  attr_reader :test_folder, :test_paths, :test_data

  def initialize(test_folder = nil)
    @test_folder = test_folder || DEFAULT_TEST_FOLDER

    unless Dir.exists?(@test_folder)
      Dir.mkdir(@test_folder)
    end
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
end
