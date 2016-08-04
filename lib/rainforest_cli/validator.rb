# frozen_string_literal: true
class RainforestCli::Validator
  API_TOKEN_ERROR = 'Please supply API token and try again'
  VALIDATIONS_PASSED = '[VALID]'
  VALIDATIONS_FAILED = '[INVALID] - Please see log to correct errors.'

  attr_reader :local_tests, :remote_tests

  def initialize(options, local_tests = nil, remote_tests = nil)
    @local_tests = local_tests || RainforestCli::TestFiles.new(options)
    @remote_tests = remote_tests || RainforestCli::RemoteTests.new(options.token)
  end

  def validate
    check_test_directory_for_tests!
    exit 1 if invalid?
  end

  def validate_with_exception!
    check_test_directory_for_tests!

    unless remote_tests.api_token_set?
      logger.error API_TOKEN_ERROR
      exit 2
    end

    exit 1 if invalid?
  end

  def invalid?
    # Assign result to variables to ensure both methods are called
    # (no short-circuiting with ||)
    parsing_errors = has_parsing_errors?
    dependency_errors = has_test_dependency_errors?
    duplicate_rfml_id_errors = has_duplicate_rfml_id_errors?

    is_invalid = parsing_errors || dependency_errors || duplicate_rfml_id_errors

    logger.info ''
    logger.info(is_invalid ? VALIDATIONS_FAILED : VALIDATIONS_PASSED)
    is_invalid
  end

  private

  def check_test_directory_for_tests!
    unless local_tests.count > 0
      logger.error "No tests found in directory: #{local_tests.test_folder}"
      exit 3
    end
  end

  def has_duplicate_rfml_id_errors?
    duped_rfml_ids_and_counts = collect_duplicate_rfml_ids
    return false unless duped_rfml_ids_and_counts.size > 0
    duplicate_rfml_ids_notification(duped_rfml_ids_and_counts)
    true
  end

  def collect_duplicate_rfml_ids
    rfml_ids_and_counts = Hash.new(0)
    local_tests.test_data.each do |test|
      rfml_ids_and_counts[test.rfml_id] += 1
    end
    rfml_ids_and_counts.select {|_, count| count > 1}
  end

  def has_parsing_errors?
    logger.info 'Validating parsing errors...'
    has_parsing_errors = rfml_tests.select { |t| t.errors.any? }

    return false unless has_parsing_errors.any?

    parsing_error_notification(has_parsing_errors)
    true
  end

  def has_test_dependency_errors?
    logger.info 'Validating embedded test IDs...'

    # Assign result to variables to ensure both methods are called
    nonexisting_tests = has_nonexisting_tests?
    circular_dependencies = has_circular_dependencies?
    nonexisting_tests || circular_dependencies
  end

  def has_nonexisting_tests?
    contains_nonexistent_ids = rfml_tests.select { |t| (t.embedded_ids - all_rfml_ids).any? }

    return false unless contains_nonexistent_ids.any?

    nonexisting_embedded_id_notification(contains_nonexistent_ids)
    true
  end

  def has_circular_dependencies?
    # TODO: Check embedded tests for remote tests as well. The Rainforest Ruby client
    # doesn't appear to support elements yet.
    has_circular_dependencies = false
    rfml_tests.each do |rfml_test|
      has_circular_dependencies ||= check_for_nested_embed(rfml_test, rfml_test.rfml_id, rfml_test.file_name)
    end
    has_circular_dependencies
  end

  def check_for_nested_embed(rfml_test, root_id, root_file)
    rfml_test.embedded_ids.each do |embed_id|
      descendant = test_dictionary[embed_id]

      # existence for embedded tests is covered in #has_nonexisting_tests?
      next unless descendant

      if descendant.embedded_ids.include?(root_id)
        circular_dependencies_notification(root_file, descendant.file_name) if descendant.embedded_ids.include?(root_id)
        return true
      end

      check_for_nested_embed(descendant, root_id, root_file)
    end
    false
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

  def parsing_error_notification(rfml_tests)
    logger.error 'Parsing errors:'
    logger.error ''
    rfml_tests.each do |rfml_test|
      logger.error "\t#{rfml_test.file_name}"
      rfml_test.errors.each do |_line, error|
        logger.error "\t#{error}"
      end
    end
    logger.error ''
  end

  def nonexisting_embedded_id_notification(rfml_tests)
    logger.error 'The following files contain unknown embedded test IDs:'
    logger.error ''
    rfml_tests.each do |rfml_test|
      logger.error "\t#{rfml_test.file_name}"
    end
    logger.error ''
  end

  def circular_dependencies_notification(file_a, file_b)
    logger.error 'The following files are embedding one another:'
    logger.error ''
    logger.error "\t#{file_a}"
    logger.error "\t#{file_b}"
    logger.error ''
  end

  def duplicate_rfml_ids_notification(duplicate_rfml_ids_and_counts)
    logger.error "The test ids are not unique!"
    logger.error ''
    duplicate_rfml_ids_and_counts.each do |rfml_id, count|
      count_str = count == 1 ? 'is 1 file' : "are #{count} files"
      logger.error "\tThere #{count_str} with an id of #{rfml_id}"
    end
    logger.error ''
  end

  def logger
    RainforestCli.logger
  end
end
