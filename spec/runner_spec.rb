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

  describe '#upload_app' do

    context 'with valid URL' do
      it 'returns the given app_source_url' do
        expect(subject.upload_app('http://my.app.url')).to eq 'http://my.app.url'
      end
    end

    context 'with invalid URL' do

      it 'errors out and exits if the file does not exists' do
        expect_any_instance_of(Logger).to receive(:fatal).with('App source file: fobar not found')
        expect do
          subject.upload_app('fobar')
        end.to raise_error(SystemExit) { |error|
          expect(error.status).to eq 1
        }
      end

      it 'errors out and exits if not an .ipa file' do
        File.should_receive(:exist?).with('fobar.txt') { true }
        expect_any_instance_of(Logger).to receive(:fatal).with('Invalid app source file: fobar.txt')
        expect do
          subject.upload_app('fobar.txt')
        end.to raise_error(SystemExit) { |error|
          expect(error.status).to eq 1
        }
      end

      it 'errors out and exits if valid file but invalid token' do
        File.should_receive(:exist?).with('fobar.ipa') { true }
        subject.client.should_receive(:get).with('/uploads', {}, retries_on_failures: true) { nil }
        expect_any_instance_of(Logger).to receive(:fatal).with(
          'Failed to upload file fobar.ipa. Please, check your API token.')
        expect do
          subject.upload_app('fobar.ipa')
        end.to raise_error(SystemExit) { |error|
          expect(error.status).to eq 1
        }
      end

      it 'errors out and exits if was not possible to upload the file' do
        File.should_receive(:exist?).with('fobar.ipa') { true }
        File.should_receive(:read).with('fobar.ipa') { 'File data' }
        url = {'host' => 'host', 'port' => 'port', 'uri' => 'uri', 'path' => 'path'}
        subject.client.should_receive(:get).with('/uploads', {}, retries_on_failures: true) { url }
        subject.should_receive(:upload_file).with(url, 'File data') { '500' }
        expect_any_instance_of(Logger).to receive(:fatal).with('Failed to upload file fobar.ipa')
        expect do
          subject.upload_app('fobar.ipa')
        end.to raise_error(SystemExit) { |error|
          expect(error.status).to eq 1
        }
      end

      it 'returns the new app_source_url in case of success' do
        File.should_receive(:exist?).with('fobar.ipa') { true }
        File.should_receive(:read).with('fobar.ipa') { 'File data' }
        url = {'host' => 'host', 'port' => 'port', 'uri' => 'uri', 'path' => 'path'}
        subject.client.should_receive(:get).with('/uploads', {}, retries_on_failures: true) { url }
        subject.should_receive(:upload_file).with(url, 'File data') { '200' }

        expect(subject.upload_app('fobar.ipa')).to eq 'path'
      end

    end
  end

  describe '#generate_junit_ouptut' do
    context 'on /runs/{run_id}.json API error' do
      before do
        allow(subject.client).to receive(:get).with('/runs/some=id.json').and_return({'error'=>'Some API error'})
      end

      it 'errors out and exits' do
        expect_any_instance_of(Logger).to receive(:fatal).with('Error retrieving results for your run: Some API error')
        expect do
          subject.generate_junit_output('some=id', 'some=filename')
        end.to raise_error(SystemExit) { |error|
          expect(error.status).to eq 1
        }
      end

    end

    context 'on /runs/{run_id}/tests.json API error' do
      before do
        allow(subject.client).to receive(:get).and_return('fake get')
        allow(subject.client).to receive(:get).with('/runs/some=id.json').and_return({'total_tests'=>'1'})
        allow(subject.client).to receive(:get).with('/runs/some=id/tests.json?page_size=1').and_return({'error'=>'Some API error'})
      end

      it 'errors and exits' do
        expect_any_instance_of(Logger).to receive(:fatal).with('Error retrieving test details for your run: Some API error')
        expect do
          subject.generate_junit_output('some=id', 'some=filename')
        end.to raise_error(SystemExit) { |error|
          expect(error.status).to eq 1
        }
      end

    end

  end

end
