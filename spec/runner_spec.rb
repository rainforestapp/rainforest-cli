# frozen_string_literal: true
describe RainforestCli::Runner do
  let(:args) { %w(run all) }
  let(:options) { RainforestCli::OptionParser.new(args) }
  subject { described_class.new(options) }

  describe '#get_environment_id' do
    context 'with an invalid URL' do
      it 'errors out and exits' do
        expect do
          subject.get_environment_id('some=weird')
        end.to raise_error(SystemExit) { |error|
          expect(error.status).to eq 2
        }
      end
    end

    context 'on API error' do
      before do
        allow(subject.client).to receive(:post).and_return({'error'=>'Some API error'})
      end

      it 'errors out and exits' do
        expect_any_instance_of(Logger).to receive(:fatal).with('Error creating the ad-hoc URL: Some API error')
        expect do
          subject.get_environment_id('http://example.com')
        end.to raise_error(SystemExit) { |error|
          expect(error.status).to eq 1
        }
      end
    end
  end

  describe '#url_valid?' do
    subject { super().url_valid?(url) }
    [
      'http://example.org',
      'https://example.org',
      'http://example.org/',
      'http://example.org?foo=bar',
    ].each do |valid_url|
      context "#{valid_url}" do
        let(:url) { valid_url }
        it { should be(true) }
      end
    end

    [
      'ftp://example.org',
      'example.org',
      '',
    ].each do |valid_url|
      context "#{valid_url}" do
        let(:url) { valid_url }
        it { should be(false) }
      end
    end
  end

end
