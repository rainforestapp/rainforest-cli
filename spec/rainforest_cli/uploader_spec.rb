# frozen_string_literal: true
describe RainforestCli::Uploader do
  let(:options) { instance_double('RainforestCli::Options', token: 'foo', test_folder: test_directory, command: '') }
  subject { described_class.new(options) }

  before do
    allow_any_instance_of(RainforestCli::Validator).to receive(:validate_with_exception!)
    allow_any_instance_of(RainforestCli::RemoteTests).to receive(:primary_key_dictionary)
      .and_return({})
  end

  describe '#upload' do
    context 'with new tests' do
      let(:test_directory) { File.expand_path(File.dirname(__FILE__) + '/../validation-examples/correct_embeds') }
      let(:rf_test_double) { instance_double('Rainforest::Test', id: 123) }

      before do
        allow(subject).to receive(:upload_test)
      end

      it 'creates uploads the new tests with no steps' do
        expect(Rainforest::Test).to receive(:create).with(hash_including(:title, :start_uri, :rfml_id, :source))
          .and_return(rf_test_double).twice
        subject.upload
      end
    end
  end

  describe 'uploaded test object' do
    let(:test_directory) { File.expand_path(File.dirname(__FILE__) + '/../rainforest-example') }
    let(:test_double) { instance_double('Rainforest::Test', id: 123) }

    before do
      allow(Rainforest::Test).to receive(:create).and_return(test_double)
    end

    it 'contains the correct attributes' do
      expect(Rainforest::Test).to receive(:update) do |_id, test_attrs|
        expect(test_attrs).to include({
                                        start_uri: '/start_uri',
                                        title: 'Example Test',
                                        source: 'rainforest-cli',
                                        tags: ['foo', 'bar', 'baz'],
                                        rfml_id: 'example_test',
                                        site_id: '456'
                                      })
      end

      subject.upload
    end
  end
end
