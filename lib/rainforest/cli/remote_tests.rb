# frozen_string_literal: true
class RainforestCli::RemoteTests
  def initialize(api_token = nil)
    Rainforest.api_key = api_token
  end

  def rfml_ids
    @rfml_ids ||= tests.map(&:rfml_id)
  end

  def tests
    if @tests.nil?
      logger.info 'Syncing with server...'
      @tests = Rainforest::Test.all(page_size: 1000)
    end
    @tests
  end

  def primary_key_dictionary
    @primary_key_dictionary ||= {}.tap do |primary_key_dictionary|
      tests.each do |rf_test|
        primary_key_dictionary[test.rfml_id] = rf_test.id
      end
    end
  end
end
