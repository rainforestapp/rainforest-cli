# frozen_string_literal: true
require 'rainforest'
require 'parallel'
require 'ruby-progressbar'
require 'rainforest/cli/notifications'

class RainforestCli::Uploader
  include RainforestCli::Notifications

  attr_reader :test_files

  def initialize(options)
    ::Rainforest.api_key = options.token
    @test_files = RainforestCli::TestFiles.new(options.test_folder)
  end

  def upload
    report_parsing_errors!
    validate_embedded_tests!

    # Create new tests first to ensure that they can be embedded
    if new_tests.any?
      logger.info 'Syncing new tests...'
      each_in_parallel(new_tests) { |rfml_test| create_test(rfml_test) }
    end

    # Update all tests
    logger.info 'Uploading tests...'
    each_in_parallel(rfml_tests) { |rfml_test| upload_test(rfml_test) }
  end

  def report_parsing_errors!
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
      descendant = rfml_id_to_test_map[embed_id]
      circular_dependencies_notification!(root_file, descendant.file_name) if descendant.embedded_ids.include?(root_id)
      check_for_nested_embed(descendant, root_id, root_file)
    end
  end

  def each_in_parallel(tests, &blk)
    progress_bar = ProgressBar.create(title: 'Rows', total: new_tests.count, format: '%a %B %p%% %t')
    Parallel.each(tests, in_threads: threads, finish: lambda { |_item, _i, _result| progress_bar.increment }) do |rfml_test|
      blk.call(rfml_test)
    end
  end

  def rfml_tests
    @rfml_tests ||= test_files.test_data
  end

  def all_rfml_ids
    local_rfml_ids + remote_rfml_ids
  end

  def local_rfml_ids
    @rfml_ids ||= test_files.rfml_ids
  end

  def remote_rfml_ids
    @remote_rfml_ids ||= rfml_id_to_primary_key_map.keys
  end

  def rfml_id_to_primary_key_map
    if @rfml_id_to_primary_key_map.nil?
      logger.info 'Syncing with server...'

      @rfml_id_to_primary_key_map = {}.tap do |rfml_id_to_primary_key_map|
        Rainforest::Test.all(page_size: 1000, rfml_ids: test_files.rfml_ids).each do |rf_test|
          rfml_id = rf_test.rfml_id
          next if rfml_id.nil?

          rfml_id_to_primary_key_map[rfml_id] = rf_test.id
        end
      end
    end
    @rfml_id_to_primary_key_map
  end

  def rfml_id_to_test_map
    @rfml_id_to_test_map ||= {}.tap do |rfml_id_to_test_map|
      rfml_tests.each { |rfml_test| rfml_id_to_test_map[rfml_test.rfml_id] = rfml_test }
    end
  end

  def new_tests
    @new_tests ||= rfml_tests.select { |t| rfml_id_to_primary_key_map[t.rfml_id].nil? }
  end

  def create_test(rfml_test)
    test_obj = {
      title: rfml_test.title,
      start_uri: rfml_test.start_uri,
      rfml_id: rfml_test.rfml_id
    }
    rf_test = Rainforest::Test.create(test_obj)

    rfml_id_to_primary_key_map[rf_test.rfml_id] = rf_test.id
  end

  def upload_test(rfml_test)
    return unless rfml_test.steps.count > 0

    test_obj = create_test_obj(rfml_test)
    # Upload the test
    begin
      if rfml_id_to_primary_key_map[rfml_test.rfml_id]
        Rainforest::Test.update(rfml_id_to_primary_key_map[rfml_test.rfml_id], test_obj)
      else
        t = Rainforest::Test.create(test_obj)
        rfml_id_to_primary_key_map[rfml_test.rfml_id] = t.id
      end
    rescue => e
      logger.fatal "Error: #{rfml_test.rfml_id}: #{e}"
      exit 2
    end
  end

  def threads
    RainforestCli::THREADS
  end

  def logger
    RainforestCli.logger
  end

  def create_test_obj(rfml_test)
    test_obj = {
      start_uri: rfml_test.start_uri || '/',
      title: rfml_test.title,
      description: rfml_test.description,
      tags: (['ro'] + rfml_test.tags).uniq,
      rfml_id: rfml_test.rfml_id
    }

    test_obj[:elements] = rfml_test.steps.map do |step|
      if step.respond_to?(:rfml_id)
        step.to_element(rfml_id_to_primary_key_map[step.rfml_id])
      else
        step.to_element
      end
    end

    unless rfml_test.browsers.empty?
      test_obj[:browsers] = rfml_test.browsers.map do|b|
        {'state' => 'enabled', 'name' => b}
      end
    end

    test_obj
  end
end
