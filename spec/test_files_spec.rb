# frozen_string_literal: true
describe RainforestCli::TestFiles do
  describe '#initialize' do
    context 'when test folder name is not supplied' do
      let(:test_folder) { described_class::DEFAULT_TEST_FOLDER }
      it 'uses the default file folder' do
        # NOTE: if default folder changes, make sure to change the folder to be
        # removed at the end. Do not delete any folders actually in the project!
        expect(test_folder).to eq('./spec/rainforest')

        expect(Dir.exist?(test_folder)).to eq(false)
        described_class.new
        expect(Dir.exist?(test_folder)).to eq(true)
        FileUtils.remove_dir('./spec/rainforest')
        expect(Dir.exist?(test_folder)).to eq(false)
      end
    end

    context 'when test folder name is supplied' do
      let(:test_folder) { './foo' }

      it 'creates the expected folder' do
        expect(Dir.exist?(test_folder)).to eq(false)
        described_class.new(test_folder)
        expect(Dir.exist?(test_folder)).to eq(true)
        FileUtils.remove_dir(test_folder)
        expect(Dir.exist?(test_folder)).to eq(false)
      end
    end

    context 'creates multiple levels of folders' do
      let(:test_folder) { './foo/bar/baz' }

      it 'creates all folders' do
        expect(Dir.exist?(test_folder)).to eq(false)
        described_class.new(test_folder)
        expect(Dir.exist?(test_folder)).to eq(true)
        FileUtils.remove_dir('./foo')
        expect(Dir.exist?(test_folder)).to eq(false)
      end
    end
  end

  describe '#test_data' do
    let(:test_folder) { File.dirname(__FILE__) + '/rainforest-example' }
    subject { described_class.new(test_folder) }

    let(:rfml_test) { subject.test_data.first }
    let(:text_file) { File.read(test_folder + '/example_test.rfml') }

    it 'parses all available tests on initialization' do
      expect(rfml_test.title).to eq(text_file.match(/^# title: (.+)$/)[1])
      expect(rfml_test.rfml_id).to eq(text_file.match(/^#! (.+?)($| .+?$)/)[1])
    end
  end

  describe '#test_dictionary' do
    Test = Struct.new(:rfml_id, :id)

    subject { described_class.new }
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
