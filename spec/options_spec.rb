# frozen_string_literal: true
describe RainforestCli::OptionParser do
  subject { RainforestCli::OptionParser.new(args) }

  describe '#initialize' do
    context 'importing csv file' do
      let(:args) { ['--import-variable-csv-file', 'some_file.csv'] }

      its(:import_file_name) { should == 'some_file.csv' }
    end

    context 'test folder (when passed)' do
      let(:args) { ['--test-folder', '/path/to/folder'] }
      its(:test_folder) { should == '/path/to/folder' }
    end

    context 'importing name' do
      let(:args) { ['--import-variable-name', 'some_name'] }
      its(:import_name) { should == 'some_name' }
    end

    context 'crowd' do
      let(:args) { ['--crowd', 'some_name'] }
      its(:crowd) { should == 'some_name' }
    end

    context 'new' do
      let(:args) { ['new', 'foo.rfml']}
      its(:command) { should == 'new' }
      its(:file_name) { should == 'foo.rfml' }
    end

    context 'rm' do
      let(:args) { ['rm', 'foo.rfml']}
      its(:command) { should == 'rm' }
      its(:file_name) { should == File.expand_path('foo.rfml') }
    end

    context 'app_source_url' do
      let(:args) { ['--app-source-url', 'some_app'] }
      its(:app_source_url) { should == 'some_app' }
    end

    context 'run all tests' do
      let(:args) { ['run', 'all'] }
      its(:tests) { should == ['all']}
      its(:conflict) { should == nil}
    end

    context 'run from tags' do
      let(:args) { ['run', '--tag', 'run-me'] }
      its(:tests) { should == []}
      its(:tags) { should == ['run-me']}
    end

    context 'run from folder' do
      let(:args) { ['run', '--folder', '12'] }
      its(:folder) { should == '12'}
    end

    context 'only run in specific browsers' do
      let(:args) { ['run', '--browsers', 'ie8'] }
      its(:browsers) { should == ['ie8']}
    end

    context 'accepts multiple browsers' do
      let(:args) { ['run', '--browsers', 'ie8,chrome'] }
      its(:browsers) { should == ['ie8', 'chrome']}
    end

    context 'it parses the --git-trigger flag' do
      let(:args) { ['run', '--git-trigger', 'all'] }
      its(:tests) { should == ['all']}
      its(:git_trigger?) { is_expected.to eq(true) }
    end

    context 'it parses the --fg flag' do
      let(:args) { ['run', '--fg', 'all'] }
      its(:tests) { should == ['all']}
      its(:foreground?) { is_expected.to eq(true) }
    end

    context 'it parses the --wait flag' do
      let(:args) { ['run', '--wait', '12345'] }
      its(:wait?) { is_expected.to be true }
      its(:run_id) { is_expected.to eq 12345 }
    end

    context 'it parses the api token' do
      let(:args) { ['run', '--token', 'abc', 'all'] }
      its(:token) { should == 'abc'}
    end

    context 'it parses the conflict flag' do
      context 'when abort' do
        let(:args) { ['run', '--conflict', 'abort', 'all'] }
        its(:conflict) { should == 'abort'}
      end

      context 'when abort-all' do
        let(:args) { ['run', '--conflict', 'abort-all', 'all']}
        its(:conflict) { should == 'abort-all' }
      end
    end

    context 'it parses the fail-fast flag' do
      let(:args) { ['run', '--fail-fast'] }
      its(:failfast?) { is_expected.to eq(true) }
    end

    context 'it parses the site-id flag' do
      let(:args) { ['run', '--site-id', '3'] }
      its(:site_id) { should eq 3 }
    end

    context 'it parses the environment-id flag' do
      let(:args) { ['run', '--environment-id', '3'] }
      its(:environment_id) { should eq 3 }
    end

    context 'it parses the custom-url flag' do
      let(:args) { ['run', '--custom-url', 'http://ad-hoc.example.com'] }
      its(:custom_url) { should eq 'http://ad-hoc.example.com' }
    end

    context 'it add a run description' do
      let(:args) { ['run', '--description', 'test description'] }
      its(:description) { should eq 'test description' }
    end
  end

  describe '#validate!' do
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

    context 'with valid arguments' do
      let(:args) { %w(--token foo) }
      it { does_not_raise_a_validation_exception }
    end

    context 'with missing token' do
      let(:args) { %w() }
      it { raises_a_validation_exception }
    end

    context 'with missing filename' do
      let(:args) { %w(--token foo rm) }
      it { raises_a_validation_exception }
    end

    context 'with a custom url but no site id' do
      let(:args) { %w(--token foo --custom-url http://www.example.com) }
      it { raises_a_validation_exception }
    end

    context 'with a import_file_name but no import name' do
      let(:args) { %w(--token foo --import-variable-csv-file foo.csv) }
      it { raises_a_validation_exception }
    end

    context 'with a import_file_name and a import_name' do
      context 'for an existing file' do
        let(:args) { %W(--token foo --import-variable-csv-file #{__FILE__} --import-variable-name my-var) }
        it { does_not_raise_a_validation_exception }
      end

      context 'for a non existing file' do
        let(:args) { %W(--token foo --import-variable-csv-file does_not_exist --import-variable-name my-var) }
        it { raises_a_validation_exception }
      end
    end
  end
end
