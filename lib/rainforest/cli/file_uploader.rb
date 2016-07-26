# frozen_string_literal: true

require 'httmultiparty'
require 'mimemagic'
require 'json'

class RainforestCli::FileUploader
  def initialize(options)
    @http_client = RainforestCli::HttpClient.new(token: options.token)
    @test_files = RainforestCli::TestFiles.new(options)
    @remote_tests = RainforestCli::RemoteTests.new(options.token)
  end

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
  end

  def upload_from_match_data(matches, rfml_test)
    test_id = @remote_tests.primary_key_dictionary[rfml_test.rfml_id]
    test_path = rfml_test.file_name
    test_dir = File.dirname(test_path)

    matches.each do |match|
      step_var, relative_file_path = match
      file_path = File.expand_path(File.join(test_dir, relative_file_path))

      if File.exist?(file_path)
        file_name = File.split(file_path).last.gsub(/[^\w\d,\.\+\/=]/, '')
        file = File.new(file_path)
        mime_type = MimeMagic.by_magic(file)

        resp = @http_client.post(
          "/tests/#{test_id}/files",
          mime_type: mime_type,
          size: file.size,
          name: file_name
        )

        if resp['aws_url']
          upload_to_aws(resp, file, mime_type)
        else
          logger.error "There was a problem with uploading your file: #{file_path}."
          logger.error resp.to_json
          exit 1
        end

        sig = resp['file_signature'][0...6]

        if step_var == 'screenshot'
          content = File.read(test_path).gsub(relative_file_path, "#{resp['file_id']}, #{sig}")
        elsif step_var == 'download'
          content = File.read(test_path).gsub(relative_file_path, "#{resp['file_id']}, #{sig}, #{file_name}")
        end

        File.open(test_path, 'w') { |f| f.puts content }
      else
        logger.warn "\t\tNo such file exists: #{file_name}"
      end
    end
  end

  def upload_to_aws(aws_info, file, mime_type)
    resp = HTTMultiParty.post(
      aws_info['aws_url'],
      query: {
        'key' => aws_info['aws_key'],
        'AWSAccessKeyId' => aws_info['aws_access_id'],
        'acl' => aws_info['aws_acl'],
        'policy' => aws_info['aws_policy'],
        'signature' => aws_info['aws_signature'],
        'Content-Type' => mime_type,
        'file' => file,
      }
    )

    unless resp.code.between?(200, 299)
      logger.fatal "There was a problem with uploading your file: #{file.path}."
      logger.fatal resp.to_json
      exit 2
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
