# frozen_string_literal: true
describe RainforestCli::Resources do
  let(:options) { instance_double('RainforestCli::Options', token: 'fake_token') }
  subject { described_class.new(options) }

  shared_examples 'a properly formatted resource' do |tested_method|
    before do
      allow_any_instance_of(RainforestCli::HttpClient).to receive(:get).and_return(api_response)
    end

    it 'calls the print method with the correct information' do
      expect(subject).to receive(:print_table) do |_name, resources|
        resource = resources.first
        expect(resource.identifier).to eq(expected_id)
        expect(resource.name).to eq(expected_name)
      end

      subject.send(tested_method)
    end
  end

  describe '#sites' do
    context 'no sites configured' do
      before do
        allow_any_instance_of(RainforestCli::HttpClient).to receive(:get).and_return([])
      end

      it 'directs you to the site settings' do
        # first line of text
        expect(RainforestCli.logger).to receive(:info).with(an_instance_of(String))
        # second line of text
        expect(RainforestCli.logger).to receive(:info) do |message|
          expect(message).to include('/settings/sites')
        end

        subject.sites
      end
    end

    context 'with sites in account' do
      let(:api_response) do
        [
          { 'id' => 123, 'name' => 'The Foo Site' },
          { 'id' => 456, 'name'=> 'The Bar Site' },
          { 'id' => 789, 'name' => 'The Baz Site' }
        ]
      end
      let(:expected_id) { api_response.first['id'] }
      let(:expected_name) { api_response.first['name'] }

      it_should_behave_like 'a properly formatted resource', :sites
    end
  end

  describe '#folders' do
    context 'no folders in account' do
      before do
        allow_any_instance_of(RainforestCli::HttpClient).to receive(:get).and_return([])
      end

      it 'directs you to the folders page' do
        # first line of text
        expect(RainforestCli.logger).to receive(:info).with(an_instance_of(String))
        # second line of text
        expect(RainforestCli.logger).to receive(:info) do |message|
          expect(message).to include('/folders')
        end

        subject.folders
      end
    end

    context 'with folders in account' do
      let(:api_response) do
        [
          { 'id' => 123, 'title' => 'The Foo Folder' },
          { 'id' => 456, 'title'=> 'The Bar Folder' },
          { 'id' => 789, 'title' => 'The Baz Folder' }
        ]
      end
      let(:expected_id) { api_response.first['id'] }
      let(:expected_name) { api_response.first['title'] }

      it_should_behave_like 'a properly formatted resource', :folders
    end
  end

  describe '#browsers' do
    let(:api_resources) do
      [
        { 'name' => 'chrome', 'description' => 'Chrome' },
        { 'name' => 'safari', 'description' => 'Safari' },
        { 'name' => 'firefox', 'description' => 'Firefox' }
      ]
    end
    let(:api_response) { { 'available_browsers' => api_resources } }
    let(:expected_id) { api_resources.first['name'] }
    let(:expected_name) { api_resources.first['description'] }

    it_should_behave_like 'a properly formatted resource', :browsers
  end

  describe '#print_table' do
    let(:resource_id) { 123456 }
    let(:resource_name) { 'resource name' }
    let(:resources) { [ RainforestCli::Resources::Resource.new(resource_id, resource_name) ] }

    it 'prints out the resources' do
      expect(subject).to receive(:puts) do |message|
        expect(message).to include('Resource ID')
        expect(message).to include('Resource Name')
      end

      # Stub dashed the line dividing table header and body
      expect(subject).to receive(:puts)

      expect(subject).to receive(:puts) do |message|
        expect(message).to include(resource_id.to_s)
        expect(message).to include(resource_name)
      end

      subject.print_table('Resource', resources)
    end
  end
end
