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
      begin
        logger.info 'Syncing tests...'
        Rainforest::Test.all(page_size: 1000)
      rescue Rainforest::ApiError => e
        logger.error "Encountered API Error: #{e.message}"
        exit 4
      end
    else
      logger.info ''
      logger.info 'No API Token set. Using local tests only...'
      logger.info ''
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

  def logger
    RainforestCli.logger
  end
end
