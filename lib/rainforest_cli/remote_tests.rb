# frozen_string_literal: true
class RainforestCli::RemoteTests
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

  private

  def fetch_tests
    if http_client.api_token_set?
      logger.info 'Fetching test data from server...'
      test_list = http_client.get('/tests/rfml_ids')
      logger.info 'Fetch complete.'
      test_list
    else
      logger.info 'No API Token set. Using local tests only...'
      []
    end
  end

  def logger
    RainforestCli.logger
  end

  def http_client
    RainforestCli.http_client
  end
end
