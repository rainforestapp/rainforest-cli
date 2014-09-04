describe Rainforest::Cli do
  before do
    Kernel.stub(:sleep)
  end

  describe ".start" do
    let(:valid_args) { %w(--token foo run --fg all) }
    let(:ok_progress) { {"state" => "in_progress", "current_progress" => {"percent" => "1"} } }

    context "with bad parameters" do
      context "with only site-id" do
        it 'errors out' do
          expect(STDOUT).to receive(:puts).with('The site-id and custom-url options work together, you need both of them.')
          begin
            described_class.start(%w(--site 3))
          rescue SystemExit => e
            # That's fine, this is expected but estede in a differnet assertion
          end
        end

        it 'exits with exit code 1' do
          expect {
            described_class.start(%w(--site 3))
          }.to raise_error { |error|
            expect(error).to be_a(SystemExit)
            expect(error.status).to eq 1
          }
        end
      end

      context "with only custom-url" do
        it 'errors out' do
          expect(STDOUT).to receive(:puts).with('The site-id and custom-url options work together, you need both of them.')
          begin
            described_class.start(%w(--custom-url http://ad-hoc.example.com))
          rescue SystemExit => e
            # That's fine, this is expected but estede in a differnet assertion
          end
        end

        it 'exits with exit code 1' do
          expect {
            described_class.start(%w(--custom-url http://ad-hoc.example.com))
          }.to raise_error { |error|
            expect(error).to be_a(SystemExit)
            expect(error.status).to eq 1
          }
        end
      end
    end

    context "a simple run" do
      before do
        described_class.stub(:post) { {"id" => 1} }
        3.times do 
          described_class.should_receive(:get) { ok_progress }
        end
        described_class.should_receive(:get) { {"state" => "complete", "result" => "passed" } }
      end

      it "should return true" do
        described_class.start(valid_args).should be_true
      end
    end

    context "a run where the server 500s after a while" do
      before do
        described_class.stub(:post) { {"id" => 1} }
        2.times do 
          described_class.should_receive(:get) { ok_progress }
        end

        described_class.should_receive(:get) { nil }

        2.times do 
          described_class.should_receive(:get) { ok_progress }
        end

        described_class.should_receive(:get) { {"state" => "complete", "result" => "passed" } }
      end

      it "should return true" do
        described_class.start(valid_args).should be_true
      end
    end
  end
end
