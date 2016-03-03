# frozen_string_literal: true
describe RainforestCli::Uploader do
  let(:test_directory) { File.expand_path(File.join(__FILE__, '../embedded-examples/correct_embeds')) }
  let(:options) { instance_double('RainforestCli::Options', token: 'foo', test_folder: test_directory) }
  subject { described_class.new(options) }

  before do
    allow_any_instance_of(RainforestCli::Validator).to receive(:validate_all!)
    allow_any_instance_of(RainforestCli::RemoteTests).to receive(:primary_key_dictionary)
      .and_return({})
  end

  describe '#upload' do
    context 'with new tests' do
      let(:progress_bar_double) { double('ProgressBar') }
      let(:rf_test_double) { instance_double('Rainforest::Test', id: 123) }

      before do
        allow(ProgressBar).to receive(:create).and_return(progress_bar_double)
        allow(progress_bar_double).to receive(:increment)
        allow(subject).to receive(:upload_test)
      end

      it 'creates uploads the new tests with no steps' do
        expect(Rainforest::Test).to receive(:create).with(hash_including(:title, :start_uri, :rfml_id))
          .and_return(rf_test_double).twice
        subject.upload
      end
    end
  end
end
