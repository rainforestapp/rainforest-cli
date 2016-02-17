# frozen_string_literal: true
describe RainforestCli::TestFiles do
  let(:test_folder) { File.dirname(__FILE__) + '/rainforest-example' }
  subject { described_class.new(test_folder) }

  let(:rfml_test) { subject.test_data.first }
  let(:text_file) { File.read(test_folder + '/example_test.rfml') }

  describe '#initialize' do
    context 'when test folder name is not supplied' do
      before do
        allow(Dir).to receive(:mkdir)
      end

      it 'uses the default file folder' do
        expect(described_class.new.test_folder).to eq(described_class::DEFAULT_TEST_FOLDER)
      end
    end

    context 'when test folder name is supplied' do
      def creates_the_expected_folder
        expect(Dir).to receive(:mkdir).with(expected_folder_arg).once
        described_class.new(folder_name)
      end

      context 'when folder name begins with ./' do
        let(:folder_name) { './foo' }
        let(:expected_folder_arg) { 'foo' }

        it { creates_the_expected_folder }
      end

      context 'when folder name has to prefix' do
        let(:folder_name) { 'foo' }
        let(:expected_folder_arg) { 'foo' }

        it { creates_the_expected_folder }
      end

      context 'when absolute path is used' do
        let(:folder_name) { '/foo' }
        let(:expected_folder_arg) { '/foo' }

        it { creates_the_expected_folder }
      end

      context 'when existing folder name already given' do
        let(:folder_name) { 'spec/rainforest-example' }

        it 'does not throw an error' do
          expect(described_class.new(folder_name)).to_not raise_error
        end
      end
    end

    context 'creates multiple levels of folders' do
      let(:folder_name) { './foo/bar/baz' }

      it 'creates all folders' do
        expect(Dir.exist?(folder_name)).to be_false
        described_class.new(folder_name)
        expect(Dir.exist?(folder_name)).to be_true
        Dir.rmdir(folder_name)
        Dir.rmdir('foo/bar')
        Dir.rmdir('foo')
      end
    end
  end

  describe '#test_data' do
    it 'parses all available tests on initialization' do
      expect(rfml_test.title).to eq(text_file.match(/^# title: (.+)$/)[1])
      expect(rfml_test.rfml_id).to eq(text_file.match(/^#! (.+?)($| .+?$)/)[1])
    end
  end
end
