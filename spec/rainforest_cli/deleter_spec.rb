# frozen_string_literal: true
describe RainforestCli::Deleter do
  let(:options) do
    instance_double(
      'RainforestCli::Options',
      token: 'foo',
      test_folder: test_directory,
      command: 'rm',
      file_name: file_name
    )
  end

  let(:test_directory) do
    File.expand_path(File.join(__FILE__, '../validation-examples/correct-embeds'))
  end

  before(:each) do
    allow(File).to receive(:delete)
  end

  subject { described_class.new(options) }

  describe '#delete' do
    context 'with incorrect file extension' do
      let(:file_name) { 'embedded_test' }

      it 'exits with the correct error' do
        begin
          expect_any_instance_of(
            Logger
          ).to receive(:fatal).with('Error: file extension must be .rfml')

          subject.delete
        rescue SystemExit => exception
          expect(exception.status).to eq(2)
        end
      end
    end

    context 'with correct file extension' do
      let(:file_name) { 'embedded_test.rfml' }
      let(:rfml_id) { 'embedded_test' }
      let(:test_id) { 25 }
      let(:primary_key_dictionary) { Hash[rfml_id, test_id] }
      let(:test_data) do
        [instance_double(RainforestCli::TestParser::Test, file_name: file_name, rfml_id: rfml_id)]
      end

      context 'with remote rfml test' do
        before do
          allow_any_instance_of(RainforestCli::RemoteTests)
            .to receive(:primary_key_dictionary).and_return(primary_key_dictionary)
          allow_any_instance_of(RainforestCli::TestFiles)
            .to receive(:test_data).and_return(test_data)
        end

        it 'deletes the remote rfml test' do
          expect(Rainforest::Test).to receive(:delete).with(test_id)
          # make sure that we don't reach SystemExit lines of file
          expect { subject.delete }.to_not raise_error
        end

        it 'deletes the local file' do
          allow(Rainforest::Test).to receive(:delete).with(test_id).and_return({})
          expect(File).to receive(:delete).with(file_name)
          # make sure that we don't reach SystemExit lines of file
          expect { subject.delete }.to_not raise_error
        end
      end

      context 'without remote rfml test' do
        let(:test_data) { [OpenStruct.new(file_name: 'foobar.rfml')] }
        let(:test_files) { OpenStruct.new(test_data: test_data) }

        before do
          allow(RainforestCli::TestFiles).to receive(:new).and_return(test_files)
        end

        it 'exits with the correct error' do
          begin
            expect_any_instance_of(
              Logger
            ).to receive(:fatal).with('Unable to delete remote rfml test')

            subject.delete
          rescue SystemExit => exception
            expect(exception.status).to eq(2)
          end
        end
      end
    end
  end
end
