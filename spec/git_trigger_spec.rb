describe RainforestCli::GitTrigger do
  subject { described_class }

  let(:default_dir) { __dir__ }
  let(:test_repo_dir) { File.join(Dir.tmpdir, 'raiforest-cli', 'test-repo') }

  describe ".last_commit_message" do
    around(:each) do |spec|
      default_dir = Dir.pwd
      FileUtils.mkdir_p test_repo_dir
      Dir.chdir(test_repo_dir)
      setup_test_repo
      begin
        spec.call
      ensure
        Dir.chdir(default_dir)
      end
    end

    # Commented because it kills CI
    xit "returns a string" do
      expect(described_class.last_commit_message).to eq "Initial commit"
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
      expect(described_class.extract_hashtags('#hello,#world')).to eq ['hello', 'world']
      expect(described_class.extract_hashtags('#dashes-work, #underscores_work #007')).to eq ['dashes-work', 'underscores_work', '007']
    end
  end

  def setup_test_repo
    FileUtils.rm_rf File.join(test_repo_dir, "*")
    [
      "git init",
      "git commit --allow-empty -m 'Initial commit'",
    ].each do |cmd|
        system cmd
      end
  end
end
