# frozen_string_literal: true
describe RainforestCli::TestImporter do
  describe '#upload' do
    let(:options) { instance_double('RainforestCli::Options', token: 'foo', test_folder: test_directory) }
    subject { described_class.new(options) }

    context 'with a missing embedded test' do
      let(:test_directory) { File.join(File.dirname(__FILE__), 'embedded_examples/missing_embeds') }

      it 'raises a TestNotFound error' do
        # TODO: Do these specs!
        # expect { subject.upload }.to raise_error(described_class::TestNotFound)
      end
    end
  end
end
