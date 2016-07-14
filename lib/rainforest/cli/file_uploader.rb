# frozen_string_literal: true
class RainforestCli::FileUploader
  def initialize(options)
    @test_files = RainforestCli::TestFiles.new(options)
  end

  # check files for the markup
  # if markup exists, get the file ids
  # if a file hasn't been created yet, create it for the id
  # upload the file for the file id and other data
  # upload to AWS
  # replace parsed string with the correct file
  def upload
    if tests_with_uploadables.empty?
      logger.info 'Nothing to upload'
    else
      # logger.info 'Upload time!'
      logger.info "Starting file uploads for #{tests_with_uploadables.count} tests:"
      tests_with_uploadables.each do |test|
        logger.info "\t#{test.title}"
      end
    end
  end

  private

  def tests_with_uploadables
    @files_with_uploadables ||= @test_files.with_uploadables
  end

  def logger
    RainforestCli.logger
  end
end
