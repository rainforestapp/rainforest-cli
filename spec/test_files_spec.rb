# frozen_string_literal: true
describe RainforestCli::TestFiles do
  describe '#test_data' do
    let(:test_directory) { File.dirname(__FILE__) + '/rainforest-example' }
    let(:options) { instance_double('RainforestCli::Options', test_folder: test_directory) }
    subject { described_class.new(options) }

    let(:rfml_test) { subject.test_data.first }
    let(:text_file) { File.read(test_directory + '/example_test.rfml') }

    it 'parses all available tests on initialization' do
      expect(rfml_test.title).to eq(text_file.match(/^# title: (.+)$/)[1])
      expect(rfml_test.rfml_id).to eq(text_file.match(/^#! (.+?)($| .+?$)/)[1])
    end
  end

  describe '#test_dictionary' do
    Test = Struct.new(:rfml_id, :id)

    let(:options) { instance_double('RainforestCli::Options', test_folder: nil) }
    subject { described_class.new(options) }
    let(:tests) { [Test.new('foo', 123), Test.new('bar', 456), Test.new('baz', 789)] }

    before do
      allow(FileUtils).to receive(:mkdir_p)
      allow_any_instance_of(described_class).to receive(:test_data)
        .and_return(tests)
    end

    it "correctly formats the dictionary's keys and values" do
      expect(subject.test_dictionary).to include({})
    end
  end
end
