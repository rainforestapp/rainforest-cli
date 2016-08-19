# frozen_string_literal: true
describe RainforestCli::RemoteTests do
  let(:options) { instance_double('RainforestCli::OptionsParser', tags: [], folder: nil, site_id: nil) }
  let(:tests) { [] }
  let(:http_client) { instance_double('RainforestCli::HttpClient', api_token_set?: true) }
  subject { described_class.new(options) }

  before do
    allow(RainforestCli).to receive(:http_client).and_return(http_client)
    allow(http_client).to receive(:get).with('/tests/rfml_ids', an_instance_of(Hash)).and_return(tests)
  end

  describe '#primary_key_dictionary' do
    Test = Struct.new(:rfml_id, :id)

    let(:test1) { { 'rfml_id' => 'foo', 'id' => 123 } }
    let(:test2) { { 'rfml_id' => 'bar', 'id' => 456 } }
    let(:test3) { { 'rfml_id' => 'baz', 'id' => 789 } }
    let(:tests) { [test1, test2, test3] }

    it "correctly formats the dictionary's keys and values" do
      expect(subject.primary_key_dictionary)
        .to include({
                      test1['rfml_id'] => test1['id'],
                      test2['rfml_id'] => test2['id'],
                      test3['rfml_id'] => test3['id'],
                    })
    end

    context 'no api token set' do
      let(:http_client) { instance_double('RainforestCli::HttpClient', api_token_set?: false) }

      it 'does not make an API call but returns an empty dictionary' do
        expect(http_client).to_not receive(:get)
        dictionary = subject.primary_key_dictionary
        expect(dictionary).to be_a(Hash)
        expect(dictionary).to eq({})
      end
    end
  end

  describe '#fetch_tests' do
    let(:url) { '/tests/rfml_ids' }
    let(:tags) { ['foo', 'bar'] }
    let(:smart_folder_id) { 123 }
    let(:site_id) { 456 }

    context 'with tags option set' do
      let(:options) { instance_double('RainforestCli::OptionsParser', tags: tags, folder: nil, site_id: nil) }

      it 'uses tags as a request parameter' do
        expect(http_client).to receive(:get) do |url, params|
          expect(url).to eq(url)
          expect(params[:tags]).to eq(tags)
        end
        subject.fetch_tests
      end
    end

    context 'with smart folder option set' do
      let(:options) { instance_double('RainforestCli::OptionsParser', tags: [], folder: smart_folder_id, site_id: nil) }

      it 'uses smart_folder_id as a request parameter' do
        expect(http_client).to receive(:get) do |url, params|
          expect(url).to eq(url)
          expect(params[:smart_folder_id]).to eq(smart_folder_id)
        end
        subject.fetch_tests
      end
    end

    context 'with site id set' do
      let(:options) { instance_double('RainforestCli::OptionsParser', tags: [], folder: nil, site_id: site_id) }

      it 'uses site_id as a request parameter' do
        expect(http_client).to receive(:get) do |url, params|
          expect(url).to eq(url)
          expect(params[:site_id]).to eq(site_id)
        end
        subject.fetch_tests
      end
    end

    context 'with all parameters set' do
      let(:options) { instance_double('RainforestCli::OptionsParser', tags: tags, folder: smart_folder_id, site_id: site_id) }

      it 'uses all request parameters' do
        expect(http_client).to receive(:get) do |url, params|
          expect(url).to eq(url)
          expect(params[:tags]).to eq(tags)
          expect(params[:smart_folder_id]).to eq(smart_folder_id)
          expect(params[:site_id]).to eq(site_id)
        end
        subject.fetch_tests
      end
    end
  end
end
