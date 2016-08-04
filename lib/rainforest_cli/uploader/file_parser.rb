# frozen_string_literal: true
require 'mimemagic'
require 'rainforest_cli/uploader/multi_form_post_request'

class RainforestCli::Uploader::FileParser
  def initialize(rfml_test, test_id, uploaded_files)
    @rfml_test = rfml_test
    @test_id = test_id
    @uploaded_files = uploaded_files
  end

  def parse_files!
    @rfml_test.steps.each do |step|
      next if step.type == :test
      parse_action_files(step) if step.uploadable_in_action
      parse_response_files(step) if step.uploadable_in_response
    end
  end

  def parse_action_files(step)
    step.uploadable_in_action.each do |match|
      step.action = replace_paths_in_text(step.action, match)
    end
  end

  def parse_response_files(step)
    step.uploadable_in_response.each do |match|
      step.response = replace_paths_in_text(step.response, match)
    end
  end

  def replace_paths_in_text(text, match)
    step_var, relative_file_path = match
    file_path = File.expand_path(File.join(test_directory, relative_file_path))

    unless File.exist?(file_path)
      logger.warn "\tError for test: #{@rfml_test.file_name}:"
      logger.warn "\t\tNo such file exists: #{File.basename(file_path)}"
      return text
    end

    file = File.open(file_path, 'rb')
    if file_already_uploaded?(file)
      aws_info = get_uploaded_data(file)
    else
      aws_info = upload_to_rainforest(file)
      upload_to_aws(file, aws_info)
    end

    sig = aws_info['file_signature'][0...6]
    if step_var == 'screenshot'
      text.gsub(relative_file_path, "#{aws_info['file_id']}, #{sig}")
    elsif step_var == 'download'
      text.gsub(relative_file_path, "#{aws_info['file_id']}, #{sig}, #{File.basename(file_path)}")
    end
  end

  def upload_to_rainforest(file)
    logger.info "\tUploading file metadata..."

    resp = http_client.post(
      "/tests/#{@test_id}/files",
      mime_type: MimeMagic.by_path(file).to_s,
      size: file.size,
      name: File.basename(file.path),
      digest: file_digest(file)
    )

    if resp['aws_url'].nil?
      logger.fatal "\tThere was a problem with uploading your file: #{file_path}."
      logger.fatal "\t\t#{resp.to_json}"
      exit 2
    end

    resp
  end

  def upload_to_aws(file, aws_info)
    logger.info "\tUploading file data..."

    resp = RainforestCli::Uploader::MultiFormPostRequest.request(
      aws_info['aws_url'],
      'key' => aws_info['aws_key'],
      'AWSAccessKeyId' => aws_info['aws_access_id'],
      'acl' => aws_info['aws_acl'],
      'policy' => aws_info['aws_policy'],
      'signature' => aws_info['aws_signature'],
      'Content-Type' => MimeMagic.by_path(file),
      'file' => file,
    )

    unless resp.code.between?(200, 299)
      logger.fatal "\tThere was a problem with uploading your file: #{file.path}."
      logger.fatal "\t\t#{resp.to_json}"
      exit 3
    end
  end

  def file_already_uploaded?(file)
    @uploaded_files.any? { |f| f['digest'] == file_digest(file) }
  end

  def get_uploaded_data(file)
    file_data = @uploaded_files.find { |f| f['digest'] == file_digest(file) }
    {
      'file_signature' => file_data['signature'],
      'file_id' => file_data['id'],
    }
  end

  private

  def file_digest(file)
    Digest::MD5.file(file).hexdigest
  end

  def test_directory
    @test_directory ||= File.dirname(@rfml_test.file_name)
  end

  def http_client
    RainforestCli.http_client
  end

  def logger
    RainforestCli.logger
  end
end
