describe Rainforest::Cli do
  let(:http_client) { Rainforest::Cli::HttpClient.any_instance }

  before do
    Kernel.stub(:sleep)
  end

  describe ".start" do
    let(:valid_args) { %w(--token foo run --fg all) }
    let(:ok_progress) do
      {
        "state" => "in_progress",
        "current_progress" => {"percent" => "1"},
        "state_details" => { "is_final_state" => false },
        "result" => "no_result",
      }
    end

    let(:complete_response) do
      {
        "state" => "complete",
        "current_progress" => {"percent" => "100"},
        "state_details" => { "is_final_state" => true},
        "result" => "passed",
      }
    end

    context "with bad parameters" do
      context "no token" do
        let(:params) { %w(--custom-url http://ad-hoc.example.com) }
        it 'errors out' do
          expect_any_instance_of(Logger).to receive(:fatal).with('You must pass your API token using: --token TOKEN')
          expect {
            described_class.start(params)
          }.to raise_error(SystemExit) { |error|
            expect(error.status).to eq 2
          }
        end
      end
    end

    context "git-trigger" do
      let(:params) { %w(--token x --git-trigger) }
      let(:commit_message) { 'a test commit message' }

      def start_with_params(params, expected_exit_code = 2)
        begin
          described_class.start(params)
        rescue SystemExit => error
          expect(error.status).to eq expected_exit_code
        end
      end

      before do
        Rainforest::Cli::GitTrigger.stub(:last_commit_message) { commit_message }
      end

      describe "with tags parameter passed" do
        let(:params) { %w(--token x --tag x --git-trigger) }

        it "warns about the parameter being ignored" do
          expect_any_instance_of(Logger).to receive(:warn).with("Specified tags are ignored when using --git-trigger")

          start_with_params(params, 0)
        end
      end

      describe "without tags parameter passed" do
        let(:params) { %w(all --token x --git-trigger) }

        it "warns about the parameter being ignored" do
          expect_any_instance_of(Logger).to receive(:warn).with("Specified tests are ignored when using --git-trigger")

          start_with_params(params, 0)
        end
      end

      describe "with no @rainforest in the commit message" do
        it "exit 0's and logs the reason" do
          expect_any_instance_of(Logger).to receive(:info).with("Not triggering as @rainforest was not mentioned in last commit message.")
          start_with_params(params, 0)
        end
      end

      describe "with @rainforest in the commit message, but no tags" do
        let(:commit_message) { 'a test commit message @rainforest' }

        it "exit 2's and logs the reason" do
          expect_any_instance_of(Logger).to receive(:error).with("Triggered via git, but no hashtags detected. Please use commit message format:")
          expect_any_instance_of(Logger).to receive(:error).with("\t'some message. @rainforest #tag1 #tag2")

          start_with_params(params, 2)
        end
      end

      describe "with @rainforest in the commit message + hashtags" do
        let(:commit_message) { 'a test commit message @rainforest #run-me' }

        it "starts the run with the specified tags" do
          http_client.should_receive(:post) do |url, options|
            expect(url).to eq("/runs")
            expect(options[:tags]).to eq(['run-me'])
            {}
          end

          start_with_params(params, 0)
        end
      end
    end

    context "with site-id and custom-url" do
      let(:params) { %w(--token x --site 3 --custom-url http://ad-hoc.example.com) }
      it "creates a new environment" do
        http_client.should_receive(:post).with("/environments",
            {
              :name => "temporary-env-for-custom-url-via-CLI",
              :url=>"http://ad-hoc.example.com"
            }
          ).and_return(
            { 'id' => 333 }
          )

        http_client.should_receive(:post).with("/runs", anything).and_return( { "id" => 1 } )

        # This is a hack because when expecting a function to be called with
        # parameters, the last call is compared but I want to compare the first
        # call, not the call to create a run, so I exit, but rescue from it here
        # so that the spec doesn't fail. It's horrible, sorry!
        begin
          described_class.start(params)
        rescue SystemExit => e
          # That's fine, this is expected but tested in a differnet assertion
        end
      end

      it "starts the run with site_id and environment_id" do
        Rainforest::Cli::Runner.any_instance.stub(get_environment_id: 333)

        http_client.should_receive(:post).with(
          "/runs",
          { :tests=>[], :site_id=>3, :environment_id=>333 }
        ).and_return( {} )
        described_class.start(params)
      end
    end

    context "with environment-id" do
      let(:params) { %w(--token x --environment 123) }

      it "starts the run with environment_id" do
        Rainforest::Cli::Runner.any_instance.stub(get_environment_id: 333)

        http_client.should_receive(:post).with(
          "/runs",
          { :tests=>[], :environment_id=>123 }
        ).and_return( {} )
        described_class.start(params)
      end
    end

    context "with smart_folder_id" do
      let(:params) { %w(--token x --folder 123) }

      it "starts the run with smart folder" do
        http_client.should_receive(:post).with(
          "/runs",
          { :smart_folder_id=>123 }
        ).and_return( {} )
        described_class.start(params)
      end
    end

    context "a simple run" do
      before do
        http_client.stub(:post) { {"id" => 1} }
        3.times do
          http_client.should_receive(:get) { ok_progress }
        end
        http_client.should_receive(:get) { complete_response }
      end

      it "should return true" do
        described_class.start(valid_args).should be_true
      end
    end

    context "a run where the server 500s after a while" do
      before do
        http_client.stub(:post) { {"id" => 1} }
        2.times do
          http_client.should_receive(:get) { ok_progress }
        end

        http_client.should_receive(:get) { nil }

        2.times do
          http_client.should_receive(:get) { ok_progress }
        end

        http_client.should_receive(:get) { complete_response }
      end

      it "should return true" do
        described_class.start(valid_args).should be_true
      end
    end
  end
end
