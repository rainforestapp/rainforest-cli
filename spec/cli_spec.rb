describe Rainforest::Cli do
  before do
    Kernel.stub(:sleep)
  end

  describe ".start" do
    let(:valid_args) { %w(--token foo run --fg all) }
    let(:ok_progress) { {"state" => "in_progress", "current_progress" => {"percent" => "1"} } }

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
