# frozen_string_literal: true
require 'rspec/its'
require 'rainforest_cli'

RainforestCli.logger = Logger.new(StringIO.new)

RSpec.configure do |config|
  config.run_all_when_everything_filtered = true
  config.filter_run :focus

  # Run specs in random order to surface order dependencies. If you find an
  # order dependency and want to debug it, you can fix the order by providing
  # the seed, which is printed after each run.
  #     --seed 1234
  config.order = 'random'

  config.before do
    progressbar_mock = double('ProgressBar')
    allow(ProgressBar).to receive(:create).and_return(progressbar_mock)
    allow(progressbar_mock).to receive(:increment)
    ENV['RAINFOREST_API_TOKEN'] = nil
  end
end

RSpec::Matchers.define :test_with_file_name do |expected_name|
  match do |actual|
    actual.file_name == expected_name
  end
end

RSpec::Matchers.define :array_excluding do |expected_exclusion|
  match do |actual|
    !actual.include?(expected_exclusion)
  end
end
