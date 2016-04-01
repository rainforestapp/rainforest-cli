# frozen_string_literal: true
describe RainforestCli::Exporter do
  let(:options) do
    instance_double('RainforestCli::Options', token: nil, test_folder: nil, debug: nil, embed_tests: nil)
  end
  subject { described_class.new(options) }

  describe '#export' do
    # Collect everything printed to file in an array-like file double object
    class FileDouble < Array
      alias_method :puts, :push
    end

    let(:file) { FileDouble.new }
    let(:tests) { [Rainforest::Test.new(id: 123, title: 'Test title')] }
    let(:embedded_rfml_id) { 'embedded_test_rfml_id' }
    let(:embedded_test) do
      {
        rfml_id: embedded_rfml_id,
        elements: [
          {
            type: 'step',
            element: {
              action: 'Embedded Action',
              response: 'Embedded Response'
            }
          }
        ]
      }
    end
    let(:single_test) do
      Rainforest::Test.new(
        {
          id: 123,
          title: 'Test title',
          start_uri: '/uri',
          tags: ['foo', 'bar'],
          browsers: [
            {
              name: 'chrome',
              state: 'enabled'
            },
            {
              name: 'safari',
              state: 'enabled'
            },
            {
              name: 'firefox',
              state: 'disabled'
            }
          ],
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
              element: embedded_test
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

      subject.export
    end

    it 'prints an action and response for a step' do
      expect(file).to include('Step Action')
      expect(file).to include('Step Response')
    end

    it 'prints embedded steps' do
      expect(file).to include('Embedded Action')
      expect(file).to include('Embedded Response')
      expect(file).to_not include("- #{embedded_rfml_id}")
    end

    it 'print enabled browsers only' do
      comments = file[0]
      expect(comments).to include('chrome')
      expect(comments).to include('safari')
      expect(comments).to_not include('firefox')
    end

    context 'with embed-tests flag' do
      let(:options) do
        instance_double('RainforestCli::Options', token: nil, test_folder: nil, debug: nil, embed_tests: true)
      end

      it 'prints an embedded test rfml id rather than the steps' do
        expect(file).to include("- #{embedded_rfml_id}")
        expect(file).to_not include('Embedded Action')
        expect(file).to_not include('Embedded Response')
      end
    end
  end
end
