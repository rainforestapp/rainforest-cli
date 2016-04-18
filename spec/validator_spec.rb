# frozen_string_literal: true

describe RainforestCli::Validator do
  let(:file_path) { File.join(test_directory, correct_file_name) }

  def notifies_with_correct_file_name
    expect(subject).to receive(notification_method)
      .with(array_including(test_with_file_name(file_path)))
      .and_call_original

    if raises_error
      expect { subject.public_send(tested_method) }.to raise_error(SystemExit)
    else
      expect { subject.public_send(tested_method) }.to_not raise_error
    end
  end

  def does_not_notify_for_wrong_file_names
    expect(subject).to receive(notification_method)
      .with(array_excluding(test_with_file_name(file_path)))
      .and_call_original

    if raises_error
      expect { subject.public_send(tested_method) }.to raise_error(SystemExit)
    else
      expect { subject.public_send(tested_method) }.to_not raise_error
    end
  end

  shared_examples 'it detects all the correct errors' do
    let(:tested_method) { :validate }
    let(:raises_error) { false }
    let(:options) { instance_double('RainforestCli::Options', test_folder: test_directory, token: 'api_token') }
    subject { described_class.new(options) }

    before do
      allow(Rainforest::Test).to receive(:all).and_return([])
    end

    context 'with parsing errors' do
      let(:notification_method) { :parsing_error_notification }
      let(:test_directory) { File.expand_path(File.join(__FILE__, '../validation-examples/parse_errors')) }

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
      let(:notification_method) { :nonexisting_embedded_id_notification }
      let(:test_directory) { File.expand_path(File.join(__FILE__, '../validation-examples/missing_embeds')) }

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
      let(:test_directory) { File.expand_path(File.join(__FILE__, '../validation-examples/circular_embeds')) }
      let(:file_name_a) { File.join(test_directory, 'test1.rfml') }
      let(:file_name_b) { File.join(test_directory, 'test2.rfml') }

      it 'raises a CircularEmbeds error for both tests' do
        expect(subject).to receive(:circular_dependencies_notification) do |a, b|
          expect([a, b] - [file_name_a, file_name_b]).to be_empty
        end.and_call_original

        expect { subject.validate_with_exception! }.to raise_error(SystemExit)
      end
    end
  end

  describe '#validate' do
    it_behaves_like 'it detects all the correct errors'

    context 'when multiple tests have the same rfml_ids' do
      let(:test_directory) { File.expand_path(File.join(__FILE__, '../validation-examples/duplicate_rfml_ids')) }
      let(:options) { instance_double('RainforestCli::Options', test_folder: test_directory, token: 'api_token') }

      subject { described_class.new(options) }

      before do
        allow(subject).to receive(:remote_rfml_ids) { [] }
      end

      it 'logs the errors' do
        expect(subject).to receive(:duplicate_rfml_ids_notification).with({'a-test' => 2}).and_call_original

        subject.validate
      end
    end
  end

  describe '#validate_with_exception!' do
    let(:tested_method) { :validate_with_exception! }
    let(:raises_error) { true }

    it_behaves_like 'it detects all the correct errors'

    context 'without a token option' do
      let(:test_directory) { File.expand_path(File.join(__FILE__, '../validation-examples')) }
      let(:options) { instance_double('RainforestCli::Options', test_folder: test_directory, token: nil) }
      subject { described_class.new(options) }

      it 'validates locally and tells the user to include a token to valid with server tests as well' do
        expect_any_instance_of(Logger).to receive(:error).with(described_class::API_TOKEN_ERROR)

        expect { subject.validate_with_exception! }.to raise_error(SystemExit)
      end
    end
  end
end
