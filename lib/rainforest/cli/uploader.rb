# frozen_string_literal: true
class RainforestCli::Uploader
  attr_reader :test_files

  def initialize(options)
    ::Rainforest.api_key = options.token
    @test_files = RainforestCli::TestFiles.new(options.test_folder)
  end

  # NOTE: Embedded tests must be successfully uploaded before parent tests.
  def upload
    validate_embedded_tests!

    # Prioritize tests so that no parent tests are not uploaded before their
    # children exist.
    priority_groups = prioritize_tests

    RainforestCli.logger.info 'Uploading tests...'

    # Upload tests in parallel if there is only one upload group (no priority).
    if priority_groups.count == 1
      upload_group_in_parallel(priority_groups.first)
    else
      upload_groups_sequentially(priority_groups)
    end
  end

  private

  def validate_embedded_tests!
    contains_nonexistent_ids = rfml_tests.select { |t| (t.embedded_ids - all_rfml_ids).any? }

    if contains_nonexistent_ids.any?
      raise TestsNotFound.new(contains_nonexistent_ids.map(&:file_name))
    end
  end

  # Prioritize embedded tests before other tests
  def prioritize_tests
    priority_groups = []
    remaining_tests = rfml_tests.dup

    until remaining_tests.empty?
      prioritized_ids = priority_groups.flatten.map(&:rfml_id)
      priority_group, remaining_tests = make_priority_group(remaining_tests, prioritized_ids)

      priority_groups << priority_group
    end

    priority_groups
  end

  def make_priority_group(rfml_tests, prioritized_ids)
    # priotizable == tests whose embedded tests have already been prioritized, if any
    prioritizable = []
    # unprioritizable == tests whose embedded tests have not been prioritized yet
    unprioritizable = []

    rfml_tests.each do |rfml_test|
      group = priotizable?(rfml_test, prioritized_ids) ? prioritizable : unprioritizable
      group << rfml_test
    end

    [prioritizable, unprioritizable]
  end

  def prioritizable?(rfml_test, prioritized_ids)
    (rfml_test.embedded_ids - prioritized_ids).empty?
  end

  def upload_groups_sequentially(upload_groups)
    progress_bar = ProgressBar.create(title: 'Rows', total: test_files.count, format: '%a %B %p%% %t')
    upload_groups.each_with_index do |rfml_tests, idx|
      if idx == (rfml_tests.length - 1)
        upload_group_in_parallel(rfml_tests, progress_bar)
      else
        rfml_tests.each { |rfml_test| upload_test(rfml_test) }
        progress_bar.increment
      end
    end
  end

  def upload_group_in_parallel(rfml_tests, progress_bar = nil)
    progress_bar ||= ProgressBar.create(title: 'Rows', total: rfml_tests.count, format: '%a %B %p%% %t')
    Parallel.each(
      rfml_tests,
      in_threads: threads,
      finish: lambda { |_item, _i, _result| progress_bar.increment }
    ) do |rfml_test|
      upload_test(rfml_test)
    end
  end

  def upload_test(rfml_test)
    return unless rfml_test.steps.count > 0

    test_obj = create_test_obj(rfml_test)
    # Upload the test
    begin
      if primary_key_ids[rfml_test.rfml_id]
        Rainforest::Test.update(primary_key_ids[rfml_test.rfml_id], test_obj)
      else
        t = Rainforest::Test.create(test_obj)
        primary_key_ids[rfml_test.rfml_id] = t.id
      end
    rescue => e
      logger.fatal "Error: #{rfml_test.rfml_id}: #{e}"
      exit 2
    end
  end

  def primary_key_ids
    if @primary_key_ids.nil?
      @primary_key_ids = {}.tap do |primary_key_ids|
        Rainforest::Test.all(page_size: 1000, rfml_ids: test_files.rfml_ids).each do |rf_test|
          rfml_id = rf_test.rfml_id
          next if rfml_id.nil?

          primary_key_ids[rfml_id] = rf_test.id
        end
      end
    end
    @primary_key_ids
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
        step.to_element(primary_key_ids[step.rfml_id])
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

  def rfml_tests
    @rfml_tests ||= test_files.test_data
  end

  def all_rfml_ids
    @rfml_ids ||= test_files.rfml_ids
  end

  def threads
    RainforestCli::THREADS
  end

  class TestsNotFound < RuntimeError
    def initialize(file_names)
      super("The following tests contain embedded tests not found in test directory:\n\t#{file_names.join("\n\t")}\n\n")
    end
  end
end
