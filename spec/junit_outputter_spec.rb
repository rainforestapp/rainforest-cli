require 'json'
require 'stringio'
#frozen_string_literal: true
describe RainforestCli::JunitOutputter do
  let(:runs_json_results) { JSON.parse(File.read("#{File.dirname(__FILE__)}/fixtures/runs_response.json")) }
  let(:tests_json_results) { JSON.parse(File.read("#{File.dirname(__FILE__)}/fixtures/tests_response.json")) }
  let(:failed_test_json) { JSON.parse(File.read("#{File.dirname(__FILE__)}/fixtures/failed_test_response.json")) }
  let(:test_io) { StringIO.new }

  describe '.build_test_suite' do

    subject { described_class.new('abc123', runs_json_results, tests_json_results) }

    before do
      allow(subject.client).to receive(:get).and_return(failed_test_json)
    end

    context 'With a valid response' do
      it 'Parses the response' do
        subject.parse
        subject.output(test_io)

        expect(test_io.string).to include('name="Test Description"')
        expect(test_io.string).to include('failures="2"')
        expect(test_io.string).to include("This feedback should appear")
        expect(test_io.string).not_to include("This feedback shouldn't &quot;appear")
      end
    end

  end # end .build_test_suite

end
