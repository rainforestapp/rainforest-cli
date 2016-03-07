# frozen_string_literal: true
class RainforestCli::RemoteTests
  def initialize(api_token = nil)
    Rainforest.api_key = api_token
  end

  def api_token_set?
    !Rainforest.api_key.nil?
  end

  def rfml_ids
    @rfml_ids ||= tests.map(&:rfml_id)
  end

  def tests
    @tests ||= fetch_tests
  end

  def fetch_tests
    if api_token_set?
      Rainforest::Test.all(page_size: 1000)
    else
      RainforestCli.logger.warn 'Please supply an API Token in order to sync with tests on the Rainforest server'
      []
    end
  end

  def primary_key_dictionary
    @primary_key_dictionary ||= {}.tap do |primary_key_dictionary|
      tests.each do |rf_test|
        primary_key_dictionary[rf_test.rfml_id] = rf_test.id
      end
    end
  end
end
