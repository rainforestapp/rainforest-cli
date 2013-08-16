describe Rainforest::Cli::OptionParser do
  subject { Rainforest::Cli::OptionParser.new(args) }

  context "run all tests" do
    let(:args) { ["run", "all"] }
    its(:tests) { should == ["all"]}
    its(:conflict) { should == nil}
  end

  context "run from tags" do
    let(:args) { ["run", "--tag", "run-me"] }
    its(:tests) { should == []}
    its(:tags) { should == ["run-me"]}
  end

  context "it parses the --fg flag" do
    let(:args) { ["run", "--fg", "all"] }
    its(:tests) { should == ["all"]}
    its(:foreground?) { should be_true }
  end

  context "it parses the api token" do
    let(:args) { ["run", "--token", "abc",  "all"] }
    its(:token) { should == "abc"}
  end

  context "it parses the conflict flag" do
    let(:args) { ["run", "--conflict", "abort",  "all"] }
    its(:conflict) { should == "abort"}
  end
end
