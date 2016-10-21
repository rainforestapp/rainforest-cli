# frozen_string_literal: true
describe RainforestCli::HttpClient do
  let(:path) { '/my/path' }
  let(:success_response) { instance_double('HTTParty::Response', code: 200, body: {'success'=>'true'}.to_json) }

  [:get, :post, :delete].each do |m|
    describe "##{m}" do
      subject { described_class.new({ token: 'foo' }) }

      it 'makes the correct type of request' do
        expect(HTTParty).to receive(m).and_return(success_response)
        subject.public_send(m, path)
      end
    end
  end

  describe '#request' do
    let(:method) { :get }
    let(:body) { { foo: :bar } }
    let(:options) { {} }
    subject { described_class.new({ token: 'foo' }).request(method, path, body, options) }

    describe 'maximum tolerated exceptions' do
      before do
        allow(HTTParty).to receive(method).and_raise(SocketError)
      end

      context 'retries_on_failures omitted' do
        it 'raises the error on the first exception' do
          expect(HTTParty).to receive(method).once
          expect { subject }.to raise_error(Http::Exceptions::HttpException)
        end
      end

      context 'retries_on_failures == false' do
        let(:options) { { retries_on_failures: false } }

        it 'raises the error on the first exception' do
          expect(HTTParty).to receive(method).once
          expect { subject }.to raise_error(Http::Exceptions::HttpException)
        end
      end

      context 'retries_on_failures == true' do
        let(:options) { { retries_on_failures: true } }
        let(:response) { instance_double('HTTParty::Response', code: 200, body: {foo: :bar}.to_json) }
        let(:delay_interval) { described_class::RETRY_INTERVAL }

        it 'it sleeps after failures before a retry' do
          expect(HTTParty).to receive(method).and_raise(SocketError).once.ordered
          expect(HTTParty).to receive(method).and_return(response).ordered
          expect(Kernel).to receive(:sleep)
          expect { subject }.to_not raise_error
        end

        it 'sleeps for longer periods with repeated exceptions' do
          expect(HTTParty).to receive(method).and_raise(SocketError).exactly(3).times.ordered
          expect(HTTParty).to receive(method).and_return(response).ordered
          expect(Kernel).to receive(:sleep).exactly(3).times
          expect { subject }.to_not raise_error
        end

        it 'returns the response upon success' do
          expect(HTTParty).to receive(method).and_raise(SocketError).once.ordered
          expect(HTTParty).to receive(method).and_return(response).ordered
          expect(Kernel).to receive(:sleep).once
          expect(subject).to eq(JSON.parse(response.body))
        end
      end
    end

    describe 'unsuccessful status codes' do
      let(:url) { 'http://some.url.com' }
      let(:bad_response) { instance_double('HTTParty::Response', code: 404, body: {'error'=>'some error', 'type'=>'some type'}.to_json) }

      before do
        allow(HTTParty).to receive(method).and_raise(SocketError)
      end

      it 'gets an error 404 and prints the error and exits' do
        expect(HTTParty).to receive(method).and_return(bad_response)
        expect_any_instance_of(Logger).to receive(:fatal).with(a_string_including('Non 200 code received'))
        expect_any_instance_of(Logger).to receive(:fatal).with(a_string_including(bad_response.body.to_s))
        expect { subject }.to raise_error(SystemExit)
      end
    end
  end
end
