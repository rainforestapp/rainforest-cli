describe Rainforest::Cli::OptionParser do
  subject { Rainforest::Cli::OptionParser.new(args) }

  describe "#initialize" do
    context "importing csv file" do
      let(:args) { ["--import-variable-csv-file", "some_file.csv"] }
      its(:import_file_name) { should == "some_file.csv" }
    end

    context "importing name" do
      let(:args) { ["--import-variable-name", "some_name"] }
      its(:import_name) { should == "some_name" }
    end

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

    context "only run in specific browsers" do
      let(:args) { ["run", "--browsers", "ie8"] }
      its(:browsers) { should == ["ie8"]}
    end

    context "raises errors with invalid browsers" do
      let(:args) { ["run", "--browsers", "lulbrower"] }

      it "raises an exception" do
        expect{subject}.to raise_error(Rainforest::Cli::BrowserException)
      end
    end

    context "accepts multiple browsers" do
      let(:args) { ["run", "--browsers", "ie8,chrome"] }
      its(:browsers) { should == ["ie8", "chrome"]}
    end

    context "it parses the --git-trigger flag" do
      let(:args) { ["run", "--git-trigger", "all"] }
      its(:tests) { should == ["all"]}
      its(:git_trigger?) { should be_true }
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

    context "it parses the fail-fast flag" do
      let(:args) { ["run", "--fail-fast"] }
      its(:failfast?) { should be_true }
    end

    context "it parses the site-id flag" do
      let(:args) { ["run", "--site-id", '3'] }
      its(:site_id) { should eq 3 }
    end

    context "it parses the custom-url flag" do
      let(:args) { ["run", "--custom-url", 'http://ad-hoc.example.com'] }
      its(:custom_url) { should eq 'http://ad-hoc.example.com' }
    end
  end

  describe "#validate!" do
    def does_not_raise_a_validation_exception
      expect do
        subject.validate!
      end.to_not raise_error
    end

    def raises_a_validation_exception
      expect do
        subject.validate!
      end.to raise_error(described_class::ValidationError)
    end

    context "with valid arguments" do
      let(:args) { %w(--token foo) }
      it { does_not_raise_a_validation_exception }
    end

    context "with missing token" do
      let(:args) { %w() }
      it { raises_a_validation_exception }
    end

    context "with a custom url but no site id" do
      let(:args) { %w(--token foo --custom-url http://www.example.com) }
      it { raises_a_validation_exception }
    end

    context "with a import_file_name but no import name" do
      let(:args) { %w(--token foo --import-variable-csv-file foo.csv) }
      it { raises_a_validation_exception }
    end

    context "with a import_file_name and a import_name" do
      context "for an existing file" do
        let(:args) { %W(--token foo --import-variable-csv-file #{__FILE__} --import-variable-name my-var) }
        it { does_not_raise_a_validation_exception }
      end

      context "for a non existing file" do
        let(:args) { %W(--token foo --import-variable-csv-file does_not_exist --import-variable-name my-var) }
        it { raises_a_validation_exception }
      end
    end
  end
end
