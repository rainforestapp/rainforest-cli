describe Rainforest::Cli do
  before do
    Kernel.stub(:sleep)
    stub_const("Rainforest::Cli::API_URL", 'http://app.rainforest.dev/api/1')
  end

  describe ".start" do
    let(:valid_args) { %w(--token foo run --fg all) }
    let(:ok_progress) { {"state" => "in_progress", "current_progress" => {"percent" => "1"} } }

    context "with bad parameters" do
      context "with custom-url with no site-id" do
        it 'errors out' do
          expect(STDOUT).to receive(:puts).with('The site-id and custom-url options work together, you need both of them.')
          begin
            described_class.start(%w(--custom-url http://ad-hoc.example.com))
          rescue SystemExit => e
            # That's fine, this is expected but tested in a differnet assertion
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

    context "with site-id and custom-url" do
      it "creates a new environment" do
        allow(described_class).to receive(:post).and_return { exit }
        expect(described_class).to receive(:post).with(
            "http://app.rainforest.dev/api/1/environments",
            {
              :name => "temporary-env-for-custom-url-via-CLI",
              :url=>"http://ad-hoc.example.com"
            }
          ).and_return(
            { 'id' => 333 }
          )

        # This is a hack because when expecting a function to be called with
        # parameters, the last call is compared but I want to compare the first
        # call, not the call to create a run, so I exit, but rescue from it here
        # so that the spec doesn't fail. It's horrible, sorry!
        begin
          described_class.start(%w(--site 3 --custom-url http://ad-hoc.example.com))
        rescue SystemExit => e
          # That's fine, this is expected but tested in a differnet assertion
        end
      end

      it "starts the run with site_id and environment_id" do
        allow(described_class).to receive(:get_environment_id).and_return(333)

        expect(described_class).to receive(:post).with(
          "http://app.rainforest.dev/api/1/runs",
          { :tests=>[], :site_id=>3, :gem_version=>Rainforest::Cli::VERSION, :environment_id=>333 }
        ).and_return( {} )
        described_class.start(%w(--site 3 --custom-url http://ad-hoc.example.com))
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

  describe ".get_environment_id" do
    context "with an invalid URL" do
      it 'errors out and exits' do
        expect(STDOUT).to receive(:puts).with("The custom URL is invalid")
        expect {
          described_class.get_environment_id('http://some=weird')
        }.to raise_error { |error|
          expect(error).to be_a(SystemExit)
          expect(error.status).to eq 1
        }
      end
    end

    context 'on API error' do
      before do
        allow(described_class).to receive(:post).and_return( {"error"=>"Some API error"} )
      end

      it 'errors out and exits' do
        expect(STDOUT).to receive(:puts).with("Error creating the ad-hoc URL: Some API error")
        expect {
          described_class.get_environment_id('http://example.com')
        }.to raise_error { |error|
          expect(error).to be_a(SystemExit)
          expect(error.status).to eq 1
        }
      end
    end
  end
end
