# frozen_string_literal: true
class RainforestCli::Validator
  attr_reader :local_tests, :remote_tests

  def initialize(local_tests, remote_tests)
    @local_tests = local_tests
    @remote_tests = remote_tests
  end

  def validate_all!
    validate_parsing_errors!
    validate_embedded_tests!
  end

  def validate_parsing_errors!
    logger.info 'Detecting parsing errors...'
    has_parsing_errors = rfml_tests.select { |t| t.errors.any? }
    parsing_error_notification!(has_parsing_errors) if has_parsing_errors.any?
  end

  def validate_embedded_tests!
    logger.info 'Validating embedded test IDs...'
    validate_embedded_test_existence!
    validate_circular_dependencies!
  end

  def validate_embedded_test_existence!
    contains_nonexistent_ids = rfml_tests.select { |t| (t.embedded_ids - all_rfml_ids).any? }
    nonexisting_embedded_id_notification!(contains_nonexistent_ids) if contains_nonexistent_ids.any?
  end

  def validate_circular_dependencies!
    # TODO: Add validation for circular dependencies in server tests as well
    rfml_tests.each do |rfml_test|
      check_for_nested_embed(rfml_test, rfml_test.rfml_id, rfml_test.file_name)
    end
  end

  private

  def check_for_nested_embed(rfml_test, root_id, root_file)
    rfml_test.embedded_ids.each do |embed_id|
      descendant = test_dictionary[embed_id]
      circular_dependencies_notification!(root_file, descendant.file_name) if descendant.embedded_ids.include?(root_id)
      check_for_nested_embed(descendant, root_id, root_file)
    end
  end

  def rfml_tests
    @rfml_tests ||= local_tests.test_data
  end

  def all_rfml_ids
    local_rfml_ids + remote_rfml_ids
  end

  def local_rfml_ids
    @local_rfml_ids ||= local_tests.rfml_ids
  end

  def remote_rfml_ids
    @remote_rfml_ids ||= remote_tests.rfml_ids
  end

  def test_dictionary
    @test_dictionary ||= local_tests.test_dictionary
  end

  def parsing_error_notification!(rfml_tests)
    logger.error 'Parsing errors:'
    logger.error ''
    rfml_tests.each do |rfml_test|
      logger.error "\t#{rfml_test.file_name}"
      rfml_test.errors.each do |_line, error|
        logger.error "\t#{error}"
      end
    end

    exit 1
  end

  def nonexisting_embedded_id_notification!(rfml_tests)
    logger.error 'The following files contain unknown embedded test IDs:'
    logger.error ''
    rfml_tests.each do |rfml_test|
      logger.error "\t#{rfml_test.file_name}"
    end

    exit 2
  end

  def circular_dependencies_notification!(file_a, file_b)
    logger.error 'The following files are embedding one another:'
    logger.error ''
    logger.error "\t#{file_a}"
    logger.error "\t#{file_b}"

    exit 3
  end

  def logger
    RainforestCli.logger
  end
end
