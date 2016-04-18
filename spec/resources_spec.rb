# frozen_string_literal: true
describe RainforestCli::Resources do
  let(:options) { instance_double('RainforestCli::Options', token: 'fake_token') }
  subject { described_class.new(options) }

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
      let(:sites) do
        [
          {
            'id' => 123,
            'name' => 'The Foo Site'
          },
          {
            'id' => 456,
            'name'=> 'The Bar Site'
          },
          {
            'id' => 789,
            'name' => 'The Baz Site'
          }
        ]
      end

      before do
        allow_any_instance_of(RainforestCli::HttpClient).to receive(:get).and_return(sites)
      end

      it 'calls the print method' do
        expect(subject).to receive(:print_table).with('Site', sites).and_yield(sites.first)
        subject.sites
      end

      it 'correctly formats the site information in the given block' do
        expect(subject).to receive(:print_table) do |_resource_name, _resource, &blk|
          site = sites.first
          expect(blk.call(site)).to include(id: site['id'], name: site['name'])
        end
        subject.sites
      end
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
      let(:folders) do
        [
          {
            'id' => 123,
            'title' => 'The Foo Folder'
          },
          {
            'id' => 456,
            'title'=> 'The Bar Folder'
          },
          {
            'id' => 789,
            'title' => 'The Baz Folder'
          }
        ]
      end

      before do
        allow_any_instance_of(RainforestCli::HttpClient).to receive(:get).and_return(folders)
      end

      it 'calls the print method' do
        expect(subject).to receive(:print_table).with('Folder', folders).and_yield(folders.first)
        subject.folders
      end

      it 'correctly formats the site information in the given block' do
        expect(subject).to receive(:print_table) do |_resource_name, _resource, &blk|
          folder = folders.first
          expect(blk.call(folder)).to include(id: folder['id'], name: folder['title'])
        end
        subject.folders
      end
    end
  end

  describe '#print_table' do
    let(:resource_id) { 123456 }
    let(:resource_name) { 'resource name' }
    let(:resources) { [ { id: resource_id, name: resource_name } ] }

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

      subject.print_table('Resource', resources) do
        { id: resource_id, name: resource_name }
      end
    end
  end
end
