# frozen_string_literal: true
describe RainforestCli::Resources do
  let(:options) { instance_double('RainforestCli::Options', token: 'fake_token') }
  subject { described_class.new(options) }

  shared_examples 'a properly formatted resource' do |tested_method|
    before do
      allow_any_instance_of(RainforestCli::HttpClient).to receive(:get).and_return(api_response)
    end

    it 'calls the print method with the correct information' do
      expect(subject).to receive(:print_table) do |_headers, rows|
        expect(rows).to eq(expected_rows)
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
          { 'id' => 123, 'name' => 'The Foo Site', 'category' => 'cat A' },
          { 'id' => 456, 'name'=> 'The Bar Site', 'category' => 'cat B' },
          { 'id' => 789, 'name' => 'The Baz Site', 'category' => 'cat C' },
        ]
      end
      let(:expected_rows) do
        [
          [123, 'The Foo Site', 'cat A'],
          [456, 'The Bar Site', 'cat B'],
          [789, 'The Baz Site', 'cat C'],
        ]
      end

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
          { 'id' => 789, 'title' => 'The Baz Folder' },
        ]
      end
      let(:expected_rows) do
        [
          [123, 'The Foo Folder'],
          [456, 'The Bar Folder'],
          [789, 'The Baz Folder'],
        ]
      end

      it_should_behave_like 'a properly formatted resource', :folders

      it 'should make a get request for 100 pages' do
        allow(subject).to receive(:print_table)

        expect_any_instance_of(RainforestCli::HttpClient).to receive(:get) do |_obj, message|
          expect(message).to include('page_size=100')
        end.and_return(api_response)
        subject.folders
      end
    end
  end

  describe '#browsers' do
    let(:api_resources) do
      [
        { 'name' => 'chrome', 'description' => 'Chrome' },
        { 'name' => 'safari', 'description' => 'Safari' },
        { 'name' => 'firefox', 'description' => 'Firefox' },
      ]
    end
    let(:api_response) { { 'available_browsers' => api_resources } }
    let(:expected_rows) do
      [
        ['chrome', 'Chrome'],
        ['safari', 'Safari'],
        ['firefox', 'Firefox'],
      ]
    end

    it_should_behave_like 'a properly formatted resource', :browsers
  end
end
