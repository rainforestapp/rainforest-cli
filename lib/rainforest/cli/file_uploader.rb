# frozen_string_literal: true

require 'mimemagic'

class RainforestCli::FileUploader
  def initialize(options)
    @http_client = RainforestCli::HttpClient.new(token: options.token)
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

  def upload_from_match_data(matches, rfml_test)
    test_id = @remote_tests.primary_key_dictionary[rfml_test.rfml_id]
    file_dir = File.dirname(rfml_test.file_name)

    matches.each do |match|
      file_name = File.expand_path(File.join(file_dir, match[1]))

      if File.exist?(file_name)
        puts "test id is #{test_id}"
        resp = @http_client.post(
          "/tests/#{test_id}/files",
          {
            mime_type: MimeMagic.by_path(file_name),
            size: File.new(file_name).size,
            name: file_name.gsub(/[^\w\d,\.\+\/=]/, ''),
          }
        )
        puts resp
      else
        logger.warn "\t\tNo such file exists: #{file_name}"
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
