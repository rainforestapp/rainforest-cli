# frozen_string_literal: true
describe RainforestCli::TestFiles do
  let(:options) { instance_double('RainforestCli::Options', test_folder: nil) }
  subject { described_class.new(options) }

  describe '#test_data' do
    let(:test_directory) { File.dirname(__FILE__) + '/rainforest-example' }
    let(:options) { instance_double('RainforestCli::Options', test_folder: test_directory) }
    let(:rfml_test) { subject.test_data.first }
    let(:text_file) { File.read(test_directory + '/example_test.rfml') }

    it 'parses all available tests on initialization' do
      expect(rfml_test.title).to eq(text_file.match(/^# title: (.+)$/)[1])
      expect(rfml_test.rfml_id).to eq(text_file.match(/^#! (.+?)($| .+?$)/)[1])
    end
  end

  describe '#test_dictionary' do
    Test = Struct.new(:rfml_id, :id)
    let(:tests) { [Test.new('foo', 123), Test.new('bar', 456), Test.new('baz', 789)] }

    before do
      allow(FileUtils).to receive(:mkdir_p)
      allow_any_instance_of(described_class).to receive(:test_data)
        .and_return(tests)
    end

    it "correctly formats the dictionary's keys and values" do
      expect(subject.test_dictionary).to include({})
    end
  end

  describe '#create_file' do
    let(:file) { instance_double('File') }
    let(:test_title) { 'My Test Title' }

    it 'sets the file name as the given test title' do
      allow(file).to receive(:write)

      expect(File).to receive(:open).with(a_string_including("#{test_title}.rfml"), 'w').and_yield(file)
      subject.create_file(test_title)
    end

    it 'sets the title of the test as the given test title' do
      allow(File).to receive(:open).and_yield(file)

      expect(file).to receive(:write).with(a_string_including("# title: #{test_title}"))
      subject.create_file(test_title)
    end

    context 'when there is an existing file with the same title' do
      before do
        allow(file).to receive(:write)
        expect(File).to receive(:exist?).twice.and_return(true).ordered
        expect(File).to receive(:exist?).and_return(false).ordered
      end

      it 'sets the file name as given test title with a number for uniqueness' do
        expect(File).to receive(:open).with(a_string_including("#{test_title} (2).rfml"), 'w').and_yield(file)
        subject.create_file(test_title)
      end
    end
  end
end
