# frozen_string_literal: true
describe RainforestCli::TestImporter do
  let(:rfml_id_regex) { /^#! (.+?)($| .+?$)/ }

  describe '#upload' do
    let(:options) { instance_double('RainforestCli::Options', token: 'foo', test_folder: test_directory) }
    subject { described_class.new(options) }

    context 'with a incorrect embedded RFML ID' do
      let(:test_directory) { File.join(File.dirname(__FILE__), 'embedded-examples/missing_embeds') }

      it 'raises a TestNotFound error' do
        expect { subject.upload }.to raise_error(described_class::TestNotFound)
      end
    end

    context 'with a correctly embedded test' do
      let(:test_directory) { File.join(File.dirname(__FILE__), 'embedded-examples/correct_embeds') }
      let(:embedded_id) { File.read("#{test_directory}/embedded_test.rfml").match(rfml_id_regex)[1] }
      let(:parent_id) { File.read("#{test_directory}/test_with_embedded.rfml").match(rfml_id_regex)[1] }

      it 'prioritizes the embedded tests before the parent tests' do
        expect_any_instance_of(described_class).to receive(:upload_groups_sequentially) do |_, test_groups|
          first_test = test_groups.dig(0, 0)
          expect(first_test.rfml_id).to eq(embedded_id)

          second_test = test_groups.dig(1, 0)
          expect(second_test.rfml_id).to eq(parent_id)
        end

        subject.upload
      end
    end
  end
end
