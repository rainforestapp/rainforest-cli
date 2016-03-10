# frozen_string_literal: true
describe RainforestCli::HttpClient do
  subject { described_class.new({ token: 'foo' }) }

  describe '#get' do
    describe 'maximum tolerated exceptions' do
      let(:url) { 'http://some.url.com' }

      before do
        allow(HTTParty).to receive(:get).and_raise(SocketError)
      end

      context 'max exceptions == 0' do
        it 'raises the error on the first exception' do
          expect(HTTParty).to receive(:get).once
          expect { subject.get(url, {}, max_exceptions: 0) }.to raise_error(Http::Exceptions::HttpException)
        end
      end

      context 'max exceptions > 0' do
        it 'raises an error after (max_exceptions + 1) tries' do
          expect(HTTParty).to receive(:get).exactly(4).times
          expect { subject.get(url, {}, max_exceptions: 3) }.to raise_error(Http::Exceptions::HttpException)
        end
      end

      context 'get a result before reaching max exceptions' do
        let(:response) { instance_double('HTTParty::Response', code: 200, body: {foo: :bar}.to_json) }

        before do
          expect(HTTParty).to receive(:get).exactly(3).times.and_raise(SocketError).ordered
          expect(HTTParty).to receive(:get).once.and_return(response).ordered
        end

        it 'does not raise an error but returns a value instead' do
          expect(subject.get(url, {}, max_exceptions: 3)).to_not be_nil
        end
      end
    end
  end
end
