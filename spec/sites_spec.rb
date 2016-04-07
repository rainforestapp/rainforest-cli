# frozen_string_literal: true
describe RainforestCli::Sites do
  let(:options) { instance_double('RainforestCli::Options', token: 'fake_token') }
  subject { described_class.new(options) }

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

  describe '#list_sites' do
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

        subject.list_sites
      end
    end

    context 'with sites configured' do
      before do
        allow_any_instance_of(RainforestCli::HttpClient).to receive(:get).and_return(sites)
      end

      it 'prints the site table' do
        expect(subject).to receive(:print_site_table).with(sites)
        subject.list_sites
      end

    end
  end

  describe '#print_site_table' do
    it 'prints out the sites' do
      expect(subject).to receive(:puts) do |message|
        expect(message).to include('Site ID')
        expect(message).to include('Site Name')
      end

      # The line dividing table header and body
      expect(subject).to receive(:puts)

      sites.each do |site|
        expect(subject).to receive(:puts) do |message|
          expect(message).to include(site['id'].to_s)
          expect(message).to include(site['name'])
        end
      end

      subject.print_site_table(sites)
    end
  end
end
