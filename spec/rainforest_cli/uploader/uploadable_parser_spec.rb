# frozen_string_literal: true
describe RainforestCli::Uploader::UploadableParser do
  let(:rfml_test) { instance_double('RainforestCli::TestParser::Test', file_name: 'rfml_test_filename.rfml') }
  let(:test_id) { 12345 }
  let(:uploaded_files) { [] }
  subject { described_class.new(rfml_test, test_id, uploaded_files) }

  describe '#replace_paths_in_text' do
    let(:file) { instance_double('File') }
    let(:file_path) { 'path/to/my/file.ext' }
    let(:file_id) { 9876 }
    let(:file_sig) { 'abcdef' }
    let(:aws_info) do
      { 'file_id' => file_id, 'file_signature' => "#{file_sig}xyz123" }
    end

    before do
      allow(subject).to receive(:test_directory).and_return('my/test/directory')
      allow(File).to receive(:exist?).and_return(true)
      allow(File).to receive(:open).and_return(file)
      allow(subject).to receive(:get_aws_upload_info).with(file).and_return(aws_info)
    end

    context 'screenshot' do
      let(:original_text) { "My screenshot is {{ file.screenshot(#{file_path}) }}" }
      let(:match) { ['screenshot', file_path] }
      let(:expected_text) { "My screenshot is {{ file.screenshot(#{file_id}, #{file_sig}) }}" }

      it 'replaces the screenshot variable arguments with the correct values' do
        expect(subject.replace_paths_in_text(original_text, match)).to eq(expected_text)
      end
    end

    context 'download' do
      let(:original_text) { "My download is {{ file.download(#{file_path}) }}" }
      let(:match) { ['download', file_path] }
      let(:file_name) { File.basename(file_path) }
      let(:expected_text) { "My download is {{ file.download(#{file_id}, #{file_sig}, #{file_name}) }}" }

      it 'replaces the download variable arguments with the correct values' do
        expect(subject.replace_paths_in_text(original_text, match)).to eq(expected_text)
      end
    end
  end

  describe '#upload_to_rainforest' do
    let(:file_path) { 'test_file.txt' }
    let(:mime_type) { 'text/plain' }
    let(:size) { 4545 }
    let(:file) { instance_double('File', path: file_path, size: size) }
    let(:digest) { 'myDigest' }

    before do
      allow(subject).to receive(:file_digest).with(file).and_return(digest)
      allow(MimeMagic).to receive(:by_path).with(file).and_return(mime_type)
    end

    it 'uploads the correct metadata to Rainforest' do
      expect(RainforestCli.http_client).to receive(:post) do |url, params|
        expect(url).to include(test_id.to_s)
        expect(params).to include(mime_type: mime_type, size: size, name: file_path, digest: digest)
      end.and_return({ 'aws_url' => 'http://aws.amazon.com/lol' })

      subject.upload_to_rainforest(file)
    end
  end
end
