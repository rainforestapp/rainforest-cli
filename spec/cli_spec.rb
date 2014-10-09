describe Rainforest::Cli do
  before do
    Kernel.stub(:sleep)
    stub_const("Rainforest::Cli::API_URL", 'http://app.rainforest.dev/api/1')
  end

  describe ".last_commit_message" do
    it "returns a string" do
      default_dir = Dir.pwd
      Dir.chdir(File.join(default_dir, %w{spec test_git_repo}))
      expect(described_class.last_commit_message).to eq "Add a test repo"
      Dir.chdir(default_dir)
    end
  end

  describe ".git_trigger_should_run?" do
    it "returns true when @rainforest is in the string" do
      expect(described_class.git_trigger_should_run?('hello, world')).to eq false
      expect(described_class.git_trigger_should_run?('hello @rainforest')).to eq true
    end
  end

  describe ".extract_hashtags" do
    it "returns a list of hashtags" do
      expect(described_class.extract_hashtags('hello, world')).to eq []
      expect(described_class.extract_hashtags('#hello, #world')).to eq ['hello', 'world']
      expect(described_class.extract_hashtags('#dashes-work, #underscores_work #007')).to eq ['dashes-work', 'underscores_work', '007']
    end
  end

  describe ".start" do
    let(:valid_args) { %w(--token foo run --fg all) }
    let(:ok_progress) { {"state" => "in_progress", "current_progress" => {"percent" => "1"} } }

    context "with bad parameters" do
      context "no token" do
        let(:params) { %w(--custom-url http://ad-hoc.example.com) }
        it 'errors out' do
          expect_any_instance_of(Logger).to receive(:fatal).with('You must pass your API token using: --token TOKEN')
          begin
            described_class.start(params)
          rescue SystemExit => e
            # That's fine, this is expected but tested in a differnet assertion
          end
        end

        it 'exits with exit code 2' do
          expect {
            described_class.start(params)
          }.to raise_error { |error|
            expect(error).to be_a(SystemExit)
            expect(error.status).to eq 2
          }
        end
      end

      context "with custom-url with no site-id" do
        let(:params) { %w(--token x --custom-url http://ad-hoc.example.com) }

        it 'errors out' do
          expect_any_instance_of(Logger).to receive(:fatal).with('The site-id and custom-url options are both required.')
          begin
            described_class.start(params)
          rescue SystemExit => e
            # That's fine, this is expected but tested in a differnet assertion
          end
        end

        it 'exits with exit code 2' do
          expect {
            described_class.start(params)
          }.to raise_error { |error|
            expect(error).to be_a(SystemExit)
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
        described_class.stub(:last_commit_message) { commit_message }
      end
      
      describe "with tags parameter passed" do
        let(:params) { %w(--token x --tag x --git-trigger) }

        it "warns about the parameter being ignored" do
          expect_any_instance_of(Logger).to receive(:warn).with("Specified tags are ignored when using the git_trigger option")
          
          start_with_params(params, 0)
        end
      end
      
      describe "with tags parameter passed" do
        let(:params) { %w(all --token x --git-trigger) }

        it "warns about the parameter being ignored" do
          expect_any_instance_of(Logger).to receive(:warn).with("Specified tests are ignored when using the git_trigger option")
          
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
          expect(described_class).to receive(:post).with(
            "http://app.rainforest.dev/api/1/runs",
            { :tags=>['run-me'], :gem_version=>Rainforest::Cli::VERSION }
          ).and_return( {} )

          start_with_params(params, 0)
        end
      end
    end

    context "with site-id and custom-url" do
      let(:params) { %w(--token x --site 3 --custom-url http://ad-hoc.example.com) }
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
          described_class.start(params)
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
        described_class.start(params)
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
        expect_any_instance_of(Logger).to receive(:fatal).with("The custom URL is invalid")
        expect {
          described_class.get_environment_id('http://some=weird')
        }.to raise_error { |error|
          expect(error).to be_a(SystemExit)
          expect(error.status).to eq 2
        }
      end
    end

    context 'on API error' do
      before do
        allow(described_class).to receive(:post).and_return( {"error"=>"Some API error"} )
      end

      it 'errors out and exits' do
        expect_any_instance_of(Logger).to receive(:fatal).with("Error creating the ad-hoc URL: Some API error")
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
