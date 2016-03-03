# frozen_string_literal: true
describe RainforestCli::Validator do
  let(:rfml_id_regex) { /^#! (.+?)($| .+?$)/ }

  describe '#validate_all!' do
    RSpec::Matchers.define :test_with_file_name do |expected_name|
      match do |actual|
        actual.file_name == expected_name
      end
    end

    def notifies_with_correct_file_name
      expect(subject).to receive(notification_method)
        .with(array_including(test_with_file_name(file_path)))
        .and_call_original
      expect { subject.validate_all! }.to raise_error(SystemExit)
    end

    def does_not_notify_for_wrong_file_names
      expect(subject).to receive(notification_method)
        .with(array_excluding(test_with_file_name(file_path)))
        .and_call_original
      expect { subject.validate_all! }.to raise_error(SystemExit)
    end

    let(:test_files) { RainforestCli::TestFiles.new(test_directory) }
    let(:remote_tests) { RainforestCli::RemoteTests.new('api_token') }
    subject { described_class.new(test_files, remote_tests) }

    let(:file_path) { File.join(test_directory, correct_file_name) }

    before do
      allow(Rainforest::Test).to receive(:all).and_return([])
    end

    context 'with parsing errors' do
      RSpec::Matchers.define :array_excluding do |expected_exclusion|
        match do |actual|
          !actual.include?(expected_exclusion)
        end
      end

      let(:notification_method) { :parsing_error_notification! }
      let(:test_directory) { File.expand_path(File.join(__FILE__, '../embedded-examples/parse_errors')) }

      context 'no rfml id' do
        let(:correct_file_name) { 'no_rfml_id.rfml' }
        it { notifies_with_correct_file_name }
      end

      context 'no question' do
        let(:correct_file_name) { 'no_rfml_id.rfml' }
        it { notifies_with_correct_file_name }
      end

      context 'no question mark' do
        let(:correct_file_name) { 'no_question_mark.rfml' }
        it { notifies_with_correct_file_name }
      end

      context 'no parse errors' do
        let(:correct_file_name) { 'no_question_mark.rfml' }
        it { does_not_notify_for_wrong_file_names }
      end
    end

    context 'with a incorrect embedded RFML ID' do
      let(:notification_method) { :nonexisting_embedded_id_notification! }
      let(:test_directory) { File.expand_path(File.join(__FILE__, '../embedded-examples/missing_embeds')) }

      context 'the file containing in the incorrect id' do
        let(:correct_file_name) { 'incorrect_test.rfml' }
        it { notifies_with_correct_file_name }
      end

      context 'the file with the correct id' do
        let(:correct_file_name) { 'correct_test.rfml' }
        it { does_not_notify_for_wrong_file_names }
      end
    end

    context 'with circular embeds' do
      let(:test_directory) { File.expand_path(File.join(__FILE__, '../embedded-examples/circular_embeds')) }
      let(:file_name_a) { File.join(test_directory, 'test1.rfml') }
      let(:file_name_b) { File.join(test_directory, 'test2.rfml') }

      it 'raises a CircularEmbeds error for both tests' do
        expect(subject).to receive(:circular_dependencies_notification!) do |a, b|
          expect([a, b] - [file_name_a, file_name_b]).to be_empty
        end.and_call_original

        expect { subject.validate_all! }.to raise_error(SystemExit)
      end
    end
  end
end
