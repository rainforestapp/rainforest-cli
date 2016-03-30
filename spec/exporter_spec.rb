# frozen_string_literal: true
describe RainforestCli::Exporter do
  let(:options) { instance_double('RainforestCli::Options', token: nil, test_folder: nil, debug: nil) }
  subject { described_class.new(options) }

  describe '#export' do
    # Collect everything printed to file in an array-like file double object
    class FileDouble < Array
      alias_method :puts, :push
    end

    let(:file) { FileDouble.new }
    let(:tests) { [Rainforest::Test.new(id: 123, title: 'Test title')] }
    let(:rfml_id) { 'embedded_test_rfml_id' }
    let(:single_test) do
      Rainforest::Test.new(
        {
          id: 123,
          title: 'Test title',
          description: '! primary_rfml_id',
          elements: [
            {
              type: 'step',
              element: {
                action: 'Step Action',
                response: 'Step Response'
              }
            },
            {
              type: 'test',
              element: { rfml_id: rfml_id }
            }
          ]
        }
      )
    end

    before do
      allow(File).to receive(:open) do |_file_name, _, &blk|
        blk.call(file)
      end

      allow_any_instance_of(RainforestCli::TestFiles).to receive(:create_file).and_return('file_name')
      allow(File).to receive(:truncate)

      allow(Rainforest::Test).to receive(:all).and_return(tests)
      allow(Rainforest::Test).to receive(:retrieve).and_return(single_test)
    end

    it 'prints an action and response for a step' do
      subject.export
      expect(file).to include('Step Action')
      expect(file).to include('Step Response')
    end

    it 'prints an embedded test rfml id rather than the steps' do
      subject.export
      expect(file).to include("- #{rfml_id}")
    end
  end
end
