# frozen_string_literal: true
class RainforestCli::RemoteTests
  def initialize(options)
    @options = options
  end

  def tests
    @tests ||= fetch_tests
  end

  def rfml_ids
    @rfml_ids ||= tests.map { |t| t['rfml_id'] }
  end

  def primary_ids
    @primary_ids ||= tests.map { |t| t['id'] }
  end

  def primary_key_dictionary
    @primary_key_dictionary ||= make_test_dictionary
  end

  def make_test_dictionary
    primary_key_dictionary = {}
    tests.each do |rf_test|
      primary_key_dictionary[rf_test['rfml_id']] = rf_test['id']
    end
    primary_key_dictionary
  end

  def fetch_tests
    if http_client.api_token_set?
      logger.info 'Fetching test data from server...'
      test_list = http_client.get('/tests/rfml_ids', filters)
      logger.info 'Fetch complete.'
      test_list
    else
      logger.info 'No API Token set. Using local tests only...'
      []
    end
  end

  private

  def filters
    {}.tap do |f|
      f[:tags] = @options.tags if @options.tags.any?
      f[:smart_folder_id] = @options.folder if @options.folder
      f[:site_id] = @options.site_id if @options.site_id
    end
  end

  def logger
    RainforestCli.logger
  end

  def http_client
    RainforestCli.http_client
  end
end
