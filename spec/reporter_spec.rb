

describe RainforestCli::Reporter do
  let(:args) { %w(report --token abc123 --run_id 12345 --junit somefile.xml) }
  let(:options) { RainforestCli::OptionParser.new(args) }

  subject { described_class.new(options) }

  describe '#report' do
    context 'on /runs/{run_id}.json API error' do
      before do
        allow(subject.client).to receive(:get).with('/runs/12345.json').and_return({'error'=>'Some API error'})
      end

      it 'errors out and exits' do
        expect_any_instance_of(Logger).to receive(:fatal).with('Error retrieving results for your run: Some API error')
        expect do
          subject.report
        end.to raise_error(SystemExit) { |error|
          expect(error.status).to eq 1
        }
      end

    end

    context 'on /runs/{run_id}/tests.json API error' do
      before do
        allow(subject.client).to receive(:get).with('/runs/12345.json').and_return({'total_tests'=>'1'})
        allow(subject.client).to receive(:get).with('/runs/12345/tests.json?page_size=1').and_return({'error'=>'Some API error'})
      end

      it 'errors and exits' do
        expect_any_instance_of(Logger).to receive(:fatal).with('Error retrieving test details for your run: Some API error')
        expect do
          subject.report
        end.to raise_error(SystemExit) { |error|
          expect(error.status).to eq 1
        }
      end

    end

    context 'with working API calls creates JunitOutputter' do
      before do
        allow(subject.client).to receive(:get).with('/runs/12345.json').and_return()
        allow(subject.client).to receive(:get).with('/runs/12345/tests.json?page_size=1').and_return()
      end
    end

  end
end
