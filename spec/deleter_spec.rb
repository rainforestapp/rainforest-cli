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
    File.stub(:delete)
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
      let(:remote_tests) { OpenStruct.new(primary_key_dictionary: {'embedded_test' => 25})}

      context 'with remote rfml test' do
        before do
          allow(
            RainforestCli::RemoteTests
          ).to receive(:new).and_return(remote_tests)
        end

        it 'deletes the remote rfml test' do
          expect(Rainforest::Test).to receive(:delete).with(25)
          subject.delete
        end

        it 'deletes the local file' do
          allow(Rainforest::Test).to receive(:delete).with(25).and_return({})
          expect(File).to receive(:delete).with('./spec/validation-examples/correct_embeds/embedded_test.rfml')
          subject.delete
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
