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
        steps_with_uploadables = test.steps.select(&:has_uploadable?)
        steps_with_uploadables.each do |step|
          upload_from_step(step, test.rfml_id)
        end
      end
    end
  end

  def upload_from_step(step, test_rfml_id)
    upload_from_match_data(step.uploadable_in_action, test_rfml_id) if step.uploadable_in_action
    upload_from_match_data(step.uploadable_in_response, test_rfml_id) if step.uploadable_in_response
    # recursive for cases with multiple uploads
    # upload_from_step(step, test_rfml_id) if step.has_uploadable?
  end

  def upload_from_match_data(match_data, test_rfml_id)
    puts "At the upload step for #{test_rfml_id}"
  end

  private

  def tests_with_uploadables
    @files_with_uploadables ||= @test_files.with_uploadables
  end

  def logger
    RainforestCli.logger
  end
end
