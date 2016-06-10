# frozen_string_literal: true
describe RainforestCli::Exporter do
  let(:options) do
    instance_double(
      'RainforestCli::Options',
      token: nil,
      test_folder: nil,
      command: nil,
      debug: nil,
      embed_tests: nil,
      tests: []
    )
  end
  subject { described_class.new(options) }

  describe '#export' do
    # Collect everything printed to file in an array-like file double object
    class FileDouble < Array
      alias_method :puts, :push

      def to_s
        join("\n")
      end
    end

    let(:file) { FileDouble.new }
    let(:file_str) { file.to_s }
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
    let(:test_elements) do
      [
        {
          type: 'test',
          redirection: true,
          element: embedded_test
        },
        {
          type: 'step',
          redirection: false,
          element: {
            action: 'Step Action',
            response: 'Step Response'
          }
        },
        {
          type: 'test',
          redirection: true,
          element: embedded_test
        },
        {
          type: 'step',
          redirection: false,
          element: {
            action: 'Last step',
            response: 'Last step?'
          }
        }
      ]
    end
    let(:single_test) do
      Rainforest::Test.new(
        {
          id: 123,
          title: 'Test title',
          start_uri: '/uri',
          tags: ['foo', 'bar'],
          browsers: [
            { name: 'chrome', state: 'enabled' },
            { name: 'safari', state: 'enabled' },
            { name: 'firefox', state: 'disabled' }
          ],
          elements: test_elements
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

    context 'action and/or question contain newlines' do
      let(:action) { "Step Action\nwith newlines\n" }
      let(:expected_action) { 'Step Action with newlines' }
      let(:response) { "Step Response\nwith\nnewlines\n" }
      let(:expected_response) { 'Step Response with newlines' }
      let(:test_elements) do
        [
          {
            type: 'step',
            element: {
              action: action,
              response: response
            }
          }
        ]
      end

      it 'removes the newlines' do
        expect(file).to include(expected_action)
        expect(file).to include(expected_response)
      end
    end

    context 'with embed-tests flag' do
      let(:options) do
        instance_double(
          'RainforestCli::Options',
          token: nil, test_folder: nil, command: nil,
          debug: nil, embed_tests: true, tests: [],
        )
      end

      it 'prints an embedded test rfml id rather than the steps' do
        expect(file_str.scan(/# redirect: true\n- #{embedded_rfml_id}/).count).to eq(2)
        expect(file_str.scan(/# redirect: false\nStep Action/).count).to eq(1)

        # The last step exists but no redirect with it
        expect(file_str.match(/Last step/)).to_not be_nil
        expect(file_str.match(/# redirect: false\nLast step/)).to be_nil

        expect(file_str).to_not include('Embedded Action')
        expect(file_str).to_not include('Embedded Response')
      end
    end

    context 'with specific tests' do
      let(:tests) { (123..127).to_a }
      let(:options) do
        instance_double(
          'RainforestCli::Options',
          token: nil, test_folder: nil, command: nil,
          debug: nil, embed_tests: nil, tests: tests
        )
      end

      it 'gets specific tests instead of all' do
        expect(Rainforest::Test).to receive(:retrieve).exactly(tests.length).times
        expect(Rainforest::Test).to receive(:all).exactly(0).times
        subject.export
      end

      it 'opens correct number of files' do
        expect(File).to receive(:open).exactly(tests.length).times
        expect_any_instance_of(RainforestCli::TestFiles).to receive(:create_file).exactly(tests.length).times
        subject.export
      end
    end
  end
end
