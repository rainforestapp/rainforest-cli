# frozen_string_literal: true

require 'json'
require 'mimemagic'

class RainforestCli::FileUploader
  require 'rainforest/cli/file_uploader/multi_form_post_request'

  def initialize(options)
    @http_client = RainforestCli::HttpClient.new(token: options.token)
    @test_files = RainforestCli::TestFiles.new(options)
    @remote_tests = RainforestCli::RemoteTests.new(options.token)
  end

  def upload
    unless upload_enabled?
      logger.fatal 'File uploads are not enabled for this environment.'
      logger.fatal 'For further information on how to enable uploads, please visit https://github.com/rainforestapp/rainforest-cli.'
      exit 1
    end

    if tests_with_uploadables.empty?
      logger.info 'Nothing to upload'
    else
      primary_key_dictionary # fetch Test IDs
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
    test_id = get_test_id(rfml_test)
    test_path = rfml_test.file_name
    test_dir = File.dirname(test_path)

    matches.each do |match|
      step_var, relative_file_path = match
      file_path = File.expand_path(File.join(test_dir, relative_file_path))

      if File.exist?(file_path)
        file_name = File.basename(file_path)
        file = File.open(file_path, 'rb')
        mime_type = MimeMagic.by_path(file_path)

        logger.info "\t\tUploading file:"
        logger.info "\t\t\t#{file_path}"

        resp = upload_to_rainforest(test_id, mime_type, file.size, file_name)
        upload_to_aws(resp, file, mime_type)

        logger.info "\t\tSuccessfully uploaded file."

        sig = resp['file_signature'][0...6]

        if step_var == 'screenshot'
          content = File.read(test_path).gsub(relative_file_path, "#{resp['file_id']}, #{sig}")
        elsif step_var == 'download'
          content = File.read(test_path).gsub(relative_file_path, "#{resp['file_id']}, #{sig}, #{file_name}")
        end

        File.open(test_path, 'w') { |f| f.puts content }
        logger.info "\t\tRFML test updated with new variable values:"
        logger.info "\t\t\t#{test_path}"
      else
        logger.warn "\t\tNo such file exists: #{file_name}"
      end
    end
  end

  def upload_to_rainforest(test_id, mime_type, file_size, file_name)
    logger.info "\t\t\tUploading metadata..."

    resp = @http_client.post(
      "/tests/#{test_id}/files",
      mime_type: mime_type,
      size: file_size,
      name: file_name
    )

    if resp['aws_url'].nil?
      logger.fatal "There was a problem with uploading your file: #{file_path}."
      logger.fatal resp.to_json
      exit 2
    end

    resp
  end

  def upload_to_aws(aws_info, file, mime_type)
    logger.info "\t\t\tUploading data..."

    resp = MultiFormPostRequest.request(
      aws_info['aws_url'],
      'key' => aws_info['aws_key'],
      'AWSAccessKeyId' => aws_info['aws_access_id'],
      'acl' => aws_info['aws_acl'],
      'policy' => aws_info['aws_policy'],
      'signature' => aws_info['aws_signature'],
      'Content-Type' => mime_type,
      'file' => file,
    )

    unless resp.code.between?(200, 299)
      logger.fatal "There was a problem with uploading your file: #{file.path}."
      logger.fatal resp.to_json
      exit 3
    end
  end

  def get_test_id(rfml_test)
    if primary_key_dictionary[rfml_test.rfml_id].nil?
      upload_empty_test(rfml_test)
    end

    primary_key_dictionary[rfml_test.rfml_id]
  end

  def upload_empty_test(rfml_test)
    test_obj = {
      title: rfml_test.title,
      start_uri: rfml_test.start_uri,
      rfml_id: rfml_test.rfml_id,
      source: 'rainforest-cli',
    }
    rf_test = Rainforest::Test.create(test_obj)

    primary_key_dictionary[rfml_test.rfml_id] = rf_test.id
  end

  private

  def upload_enabled?
    ENV.fetch('RAINFOREST_ENABLE_FILE_UPLOAD', false)
  end

  def primary_key_dictionary
    @remote_tests.primary_key_dictionary
  end

  def tests_with_uploadables
    @files_with_uploadables ||= @test_files.with_uploadables
  end

  def logger
    RainforestCli.logger
  end
end
