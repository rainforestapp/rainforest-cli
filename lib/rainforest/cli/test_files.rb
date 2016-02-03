class RainforestCli::TestFiles
  DEFAULT_TEST_FOLDER = '/spec/rainforest'.freeze
  EXT = ".rfml".freeze

  attr_reader :test_paths

  def initialize(test_folder = nil)
    @test_paths = "#{test_folder || DEFAULT_TEST_FOLDER}/**/*#{EXT}"
  end

  def file_extension
    EXT
  end
end
