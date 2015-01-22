describe Rainforest::Cli::GitTrigger do
  subject { described_class }

  describe ".last_commit_message" do
    xit "returns a string" do
      default_dir = Dir.pwd

      Dir.chdir(File.join([default_dir, 'spec', 'test-repo']))
      expect(described_class.last_commit_message).to eq "Initial commit"
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
      expect(described_class.extract_hashtags('#hello,#world')).to eq ['hello', 'world']
      expect(described_class.extract_hashtags('#dashes-work, #underscores_work #007')).to eq ['dashes-work', 'underscores_work', '007']
    end
  end
end
