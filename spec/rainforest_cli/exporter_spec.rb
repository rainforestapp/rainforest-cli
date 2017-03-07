# frozen_string_literal: true
describe RainforestCli::Exporter do
  let(:options) do
    instance_double(
      'RainforestCli::Options',
      token: 'token',
      test_folder: nil,
      command: nil,
      debug: nil,
      embed_tests: nil,
      tests: [],
      tags: [],
      folder: nil,
      site_id: nil
    )
  end
  let(:http_client) { instance_double(RainforestCli::HttpClient, api_token_set?: true) }
  subject { described_class.new(options) }

  before do
    allow(RainforestCli).to receive(:http_client).and_return(http_client)
  end

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
    let(:tests) { [{ 'id' => 123, 'rfml_id' => 'rfml_id_123' }] }
    let(:embedded_rfml_id) { 'embedded_test_rfml_id' }
    let(:embedded_test) do
      {
        'rfml_id' => embedded_rfml_id,
        'elements' => [
          {
            'type' => 'step',
            'element' => {
              'action' => 'Embedded Action',
              'response' => 'Embedded Response',
            },
          },
        ],
      }
    end
    let(:test_elements) do
      [
        {
          'type' => 'test',
          'redirection' => true,
          'element' => embedded_test,
        },
        {
          'type' => 'step',
          'redirection' => false,
          'element' => {
            'action' => 'Step Action',
            'response' => 'Step Response',
          },
        },
        {
          'type' => 'test',
          'redirection' => true,
          'element' => embedded_test,
        },
        {
          'type' => 'step',
          'redirection' => false,
          'element' => {
            'action' => 'Last step',
            'response' => 'Last step?',
          },
        },
      ]
    end
    let(:single_test_id) { 123 }
    let(:single_test_site_id) { 987 }
    let(:single_test) do
      {
        'id' => single_test_id,
        'title' => 'Test title',
        'start_uri' => '/uri',
        'site_id' => single_test_site_id,
        'tags' => ['foo', 'bar'],
        'browsers' => [
          { 'name' => 'chrome', 'state' => 'enabled' },
          { 'name' => 'safari', 'state' => 'enabled' },
          { 'name' => 'firefox', 'state' => 'disabled' },
        ],
        'elements' => test_elements,
      }
    end

    before do
      allow(File).to receive(:open) do |_file_name, _, &blk|
        blk.call(file)
      end

      allow_any_instance_of(RainforestCli::TestFiles).to receive(:create_file).and_return('file_name')
      allow(File).to receive(:truncate)

      allow(http_client).to receive(:get).with('/tests/rfml_ids', an_instance_of(Hash))
        .and_return(tests)
      allow(http_client).to receive(:get).with("/tests/#{single_test_id}").and_return(single_test)
    end

    it 'prints an action and response for a step' do
      subject.export
      expect(file).to include('Step Action')
      expect(file).to include('Step Response')
    end

    it 'prints embedded steps' do
      subject.export
      expect(file).to include('Embedded Action')
      expect(file).to include('Embedded Response')
      expect(file).to_not include("- #{embedded_rfml_id}")
    end

    it 'prints enabled browsers only' do
      subject.export
      meta_data = file[0]
      expect(meta_data).to include('chrome')
      expect(meta_data).to include('safari')
      expect(meta_data).to_not include('firefox')
    end

    it 'prints site id' do
      subject.export
      meta_data = file[0]
      expect(meta_data).to include("# site_id: #{single_test_site_id}")
    end

    context 'action and/or question contain newlines' do
      let(:action) { "Step Action\nwith newlines\n" }
      let(:expected_action) { 'Step Action with newlines' }
      let(:response) { "Step Response\nwith\nnewlines\n" }
      let(:expected_response) { 'Step Response with newlines' }
      let(:test_elements) do
        [
          {
            'type' => 'step',
            'element' => {
              'action' => action,
              'response' => response,
            },
          },
        ]
      end

      it 'removes the newlines' do
        subject.export
        expect(file).to include(expected_action)
        expect(file).to include(expected_response)
      end
    end

    context 'with embed-tests flag' do
      let(:options) do
        instance_double(
          'RainforestCli::Options',
          token: 'token', test_folder: nil, command: nil, debug: nil, embed_tests: true,
          tests: [], tags: [], folder: nil, site_id: nil
        )
      end

      it 'prints an embedded test rfml id' do
        subject.export
        expect(file).to include("- #{embedded_rfml_id}")
        expect(file_str).to_not include('Embedded Action')
        expect(file_str).to_not include('Embedded Response')
      end

      it 'prints the redirects in the correct location' do
        subject.export
        # the first embedded test should not have a redirect before it
        expect(file_str.scan(/# redirect: true\n- #{embedded_rfml_id}/).count).to eq(1)

        # First real step should have a redirect
        expect(file_str).to include("# redirect: false\nStep Action")

        # The last step exists but no redirect with it
        expect(file_str).to include('Last step')
        expect(file_str).to_not include("# redirect: false\nLast step")
      end
    end

    context 'with specific tests' do
      let(:test_ids) { (123..127).to_a }
      let(:options) do
        instance_double(
          'RainforestCli::Options',
          token: nil, test_folder: nil, command: nil,
          debug: nil, embed_tests: nil, tests: test_ids
        )
      end

      before do
        allow(http_client).to receive(:get).with(/\/tests\/\d+/).and_return(single_test)
      end

      it 'gets specific tests instead of all' do
        expect(http_client).to receive(:get).with(/\/tests\/\d+/).exactly(test_ids.length)
        expect_any_instance_of(RainforestCli::RemoteTests).to_not receive(:primary_ids)
        subject.export
      end

      it 'opens correct number of files' do
        expect(File).to receive(:open).exactly(test_ids.length).times
        expect_any_instance_of(RainforestCli::TestFiles).to receive(:create_file).exactly(test_ids.length).times
        subject.export
      end
    end
  end
end
