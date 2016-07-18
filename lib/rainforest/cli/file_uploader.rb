# frozen_string_literal: true
class RainforestCli::FileUploader
  def initialize(options)
    @test_files = RainforestCli::TestFiles.new(options)
    @remote_tests = RainforestCli::RemoteTests.new(options.token)
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
      @remote_tests.primary_key_dictionary # fetch Test IDs
      logger.info "Starting file uploads for #{tests_with_uploadables.count} tests:"
      tests_with_uploadables.each do |rfml_test|
        logger.info "\t#{rfml_test.title}"
        steps_with_uploadables = rfml_test.steps.select(&:has_uploadable?)
        steps_with_uploadables.each do |step|
          upload_from_step(step, rfml_test)
        end
      end
    end
  end

  def upload_from_step(step, rfml_test)
    upload_from_match_data(step.uploadable_in_action, rfml_test) if step.uploadable_in_action
    upload_from_match_data(step.uploadable_in_response, rfml_test) if step.uploadable_in_response
    # recursive for cases with multiple uploads
    # upload_from_step(step, rfml_test) if step.has_uploadable?
  end

  def upload_from_match_data(match_data, rfml_test)
    test_id = @remote_tests.primary_key_dictionary[rfml_test.rfml_id]
    file_dir = File.dirname(rfml_test.file_name)
    file_name = File.expand_path(File.join(file_dir, match_data[2]))

    if File.exist?(file_name)
      puts "/tests/#{test_id}/files"
    else
      logger.warn "\t\tNo such file exists: #{file_name}"
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
