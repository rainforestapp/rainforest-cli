# frozen_string_literal: true
describe RainforestCli::TestParser do
  describe RainforestCli::TestParser::Parser do
    subject { described_class.new(file_name) }

    describe '#initialize' do
      let(:file_name) { './spec/rainforest-example/example_test.rfml' }

      it 'expands the file name path' do
        test = subject.instance_variable_get(:@test)
        expect(test.file_name).to eq(File.expand_path(file_name))
      end
    end
  end
end
