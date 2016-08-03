# frozen_string_literal: true
require 'mimemagic'
require 'rainforest/cli/uploader/multi_form_post_request'

class RainforestCli::Uploader::FileParser
  def initialize(rfml_test, test_id)
    @rfml_test = rfml_test
    @test_id = test_id
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
      logger.warn "\t\tNo such file exists: #{File.basename(file_path)}"
      return
    end

    file = File.open(file_path, 'rb')

    rf_response = upload_to_rainforest(file)
    upload_to_aws(file, rf_response)

    sig = rf_response['file_signature'][0...6]
    if step_var == 'screenshot'
      text.gsub(relative_file_path, "#{rf_response['file_id']}, #{sig}")
    elsif step_var == 'download'
      text.gsub(relative_file_path, "#{rf_response['file_id']}, #{sig}, #{File.basename(file_path)}")
    end
  end

  def upload_to_rainforest(file)
    logger.info "\t\t\tUploading metadata..."

    resp = http_client.post(
      "/tests/#{@test_id}/files",
      mime_type: MimeMagic.by_path(file).to_s,
      size: file.size,
      name: File.basename(file.path),
      digest: Digest::MD5.file(file).hexdigest
    )

    if resp['aws_url'].nil?
      logger.fatal "There was a problem with uploading your file: #{file_path}."
      logger.fatal resp.to_json
      exit 2
    end

    resp
  end

  def upload_to_aws(file, aws_info)
    logger.info "\t\t\tUploading data..."

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
      logger.fatal "There was a problem with uploading your file: #{file.path}."
      logger.fatal resp.to_json
      exit 3
    end
  end

  private

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