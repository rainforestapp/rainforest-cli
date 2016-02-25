# frozen_string_literal: true
describe RainforestCli::Uploader do
  let(:rfml_id_regex) { /^#! (.+?)($| .+?$)/ }

  describe '#upload' do
    let(:options) { instance_double('RainforestCli::Options', token: 'foo', test_folder: test_directory, debug: true) }
    subject { described_class.new(options) }

    before do
      allow(subject).to receive(:rfml_id_to_primary_key_map).and_return({})
    end

    context 'with a incorrect embedded RFML ID' do
      let(:test_directory) { File.join(File.dirname(__FILE__), 'embedded-examples/missing_embeds') }

      it 'raises a TestNotFound error for tests embedding a missing test' do
        expect { subject.upload }.to raise_error do |error|
          expect(error).to be_a(described_class::TestsNotFound)
          expect(error.message).to include('parent_test.rfml')
          expect(error.message).to include('other_parent_test.rfml')
        end
      end

      it 'does not raise the error for tests embedding the offending tests' do
        expect { subject.upload }.to raise_error do |error|
          expect(error).to be_a(described_class::TestsNotFound)
          expect(error.message).to_not include('correct_test.rfml')
        end
      end
    end

    context 'with circular embeds' do
      let(:test_directory) { File.join(File.dirname(__FILE__), 'embedded-examples/circular_embeds') }

      it 'raises a CircularEmbeds error for both tests' do
        expect { subject.upload }.to raise_error do |error|
          expect(error).to be_a(described_class::CircularEmbeds)
          # raises in the first circular embed it finds
          match1 = error.message.match(/test1\.rfml/)
          match2 = error.message.match(/test2\.rfml/)
          expect(match1 || match2).to be_truthy
        end
      end
    end

    context 'with new tests' do
      let(:test_directory) { File.join(File.dirname(__FILE__), 'embedded-examples/correct_embeds') }
      let(:embedded_id) { File.read("#{test_directory}/embedded_test.rfml").match(rfml_id_regex)[1] }
      let(:parent_id) { File.read("#{test_directory}/test_with_embedded.rfml").match(rfml_id_regex)[1] }
      let(:progress_bar_double) { double('ProgressBar') }

      before do
        allow(ProgressBar).to receive(:create).and_return(progress_bar_double)
        allow(progress_bar_double).to receive(:increment)
        allow(subject).to receive(:upload_test)
      end

      it 'creates uploads the new tests with no steps' do
        expect(subject).to receive(:create_test) do |rfml_test|
          expect([embedded_id, parent_id]).to include(rfml_test.rfml_id)
        end.twice

        subject.upload
      end
    end
  end
end
