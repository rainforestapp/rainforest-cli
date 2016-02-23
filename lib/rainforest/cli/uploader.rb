class RainforestCli::Uploader
  attr_reader :test_files

  def initialize(options)
    ::Rainforest.api_key = options.token
    @test_files = RainforestCli::TestFiles.new(options.test_folder)
  end

  def upload
    upload_groups = make_test_priority_groups

    RainforestCli.logger.info 'Uploading tests...'

    # Upload in parallel if order doesn't matter
    if upload_groups.count > 1
      upload_groups_sequentially(upload_groups)
    else
      upload_group_in_parallel(upload_groups.first)
    end
  end

  private

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
      in_threads: THREADS,
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
      if rfml_id_mappings[rfml_test.rfml_id]
        t = Rainforest::Test.update(rfml_id_mappings[rfml_test.rfml_id], test_obj)
      else
        t = Rainforest::Test.create(test_obj)
        rfml_id_mappings[rfml_test.rfml_id] = t.id
      end
    rescue => e
      logger.fatal "Error: #{rfml_test.rfml_id}: #{e}"
      exit 2
    end
  end

  def make_test_priority_groups
    # Prioritize embedded tests before other tests
    upload_groups = []
    unordered_tests = []
    queued_tests = test_files.test_data.dup

    until queued_tests.empty?
      new_ordered_group = []
      ordered_ids = upload_groups.flatten.map(&:rfml_id)

      queued_tests.each do |rfml_test|
        if (rfml_test.embedded_ids - ordered_ids).empty?
          new_ordered_group << rfml_test
        else
          unordered_tests << rfml_test
        end
      end

      # If all the queued tests make it to the unordered tests group, then
      # they contain non-existent RFML ids.
      if queued_tests.length == unordered_tests.length
        misconfigured_tests = filter_misconfigured_tests(queued_tests)
        raise TestNotFound.new(misconfigured_tests.map(&:file_name))
      end

      upload_groups << new_ordered_group
      queued_tests = unordered_tests
      unordered_tests = []
    end

    upload_groups
  end

  # Filter out tests that depend on the actual misconfigured tests
  def filter_misconfigured_tests(unfiltered_tests)
    all_ids = unfiltered_tests.map(&:rfml_id)
    unfiltered_tests.reject { |test| (test.embedded_ids - all_ids).empty? }
  end

  def rfml_id_mappings
    if @id_mappings.nil?
      @id_mappings = {}.tap do |id_mappings|
        Rainforest::Test.all(page_size: 1000, rfml_ids: test_files.rfml_ids).each do |rf_test|
          rfml_id = rf_test.rfml_id
          next if rfml_id.nil?

          id_mappings[rfml_id] = rf_test.id
        end
      end
    end
    @id_mappings
  end

  class TestNotFound < RuntimeError
    def initialize(file_names)
      super("The following tests contain embedded tests not found in test directory:\n\t#{file_names.join("\n\t")}\n\n")
    end
  end
end
