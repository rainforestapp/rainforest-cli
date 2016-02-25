# frozen_string_literal: true
describe RainforestCli::Uploader do
  let(:rfml_id_regex) { /^#! (.+?)($| .+?$)/ }

  describe '#upload' do
    RSpec::Matchers.define :test_with_file_name do |expected_name|
      match do |actual|
        actual.file_name == expected_name
      end
    end

    def notifies_with_correct_file_name
      expect(subject).to receive(notification_method)
        .with(array_including(test_with_file_name(file_path)))
        .and_call_original
      expect { subject.upload }.to raise_error(SystemExit)
    end

    def does_not_notify_for_wrong_file_names
      expect(subject).to receive(notification_method)
        .with(array_excluding(test_with_file_name(file_path)))
        .and_call_original
      expect { subject.upload }.to raise_error(SystemExit)
    end

    let(:options) { instance_double('RainforestCli::Options', token: 'foo', test_folder: test_directory, debug: true) }
    subject { described_class.new(options) }

    let(:file_path) { File.join(test_directory, correct_file_name) }

    before do
      allow(Rainforest::Test).to receive(:all).and_return([])
      allow(subject).to receive(:rfml_id_to_primary_key_map).and_return({})
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

        expect { subject.upload }.to raise_error(SystemExit)
      end
    end

    context 'with new tests' do
      let(:test_directory) { File.expand_path(File.join(__FILE__, '../embedded-examples/correct_embeds')) }
      let(:embedded_id) { File.read("#{test_directory}/embedded_test.rfml").match(rfml_id_regex)[1] }
      let(:parent_id) { File.read("#{test_directory}/test_with_embedded.rfml").match(rfml_id_regex)[1] }
      let(:progress_bar_double) { double('ProgressBar') }

      before do
        allow(ProgressBar).to receive(:create).and_return(progress_bar_double)
        allow(progress_bar_double).to receive(:increment)
        allow(subject).to receive(:upload_test)
      end

      it 'creates uploads the new tests with no steps' do
        expect(subject).to receive(:create_test) do |rfml_test|
          expect([embedded_id, parent_id]).to include(rfml_test.rfml_id)
        end.twice

        subject.upload
      end
    end
  end
end
