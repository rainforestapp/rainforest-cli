# frozen_string_literal: true
describe RainforestCli::TestParser do
  describe RainforestCli::TestParser::Parser do
    subject { described_class.new(file_name) }
    let(:file_name) { './spec/rainforest-example/example_test.rfml' }

    describe '#initialize' do
      it 'expands the file name path' do
        test = subject.instance_variable_get(:@test)
        expect(test.file_name).to eq(File.expand_path(file_name))
      end
    end

    describe '#process' do
      let(:parsed_test) { described_class.new(file_name).process }

      describe 'parsing comments' do
        it 'properly identifies comments' do
          expect(parsed_test.description).to include('This is a comment', 'this is also a comment')
        end

        it 'properly identifies metadata fields' do
          expect(parsed_test.title).to eq('Example Test')
          expect(parsed_test.tags).to eq(['foo', 'bar', 'baz'])
          expect(parsed_test.site_id).to eq('456')
          expect(parsed_test.start_uri).to eq('/start_uri')
        end

        it 'properly parses step metadata' do
          parsed_step = parsed_test.steps.first
          expect(parsed_step.redirect).to eq('false')
        end
      end

      describe 'errors' do
        context 'no title' do
          let(:file_name) { './spec/validation-examples/parse_errors/no_title.rfml' }

          specify do
            expect(parsed_test.errors[:title]).to_not be_nil
          end
        end
      end
    end
  end
end
