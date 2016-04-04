# frozen_string_literal: true
describe RainforestCli::TestParser do
  describe RainforestCli::TestParser::Step do
    subject { described_class.new('action', 'response', 'redirect') }

    describe '#to_element' do
      it 'includes redirection key/value pair' do
        expect(subject.to_element).to include(redirection: 'redirect')
      end
    end
  end

  describe RainforestCli::TestParser::Parser do
    let(:file_name) { File.dirname(__FILE__) + '/rainforest-example/example_test.rfml' }
    subject { described_class.new(file_name) }

    describe '#process' do
      context 'redirection' do
        it 'sets redirection to true by default' do
          test = subject.process
          step = test.steps.first
          expect(step.redirection).to eq('true')
        end

        context 'redirection specified as false' do
          let(:file_name) { File.dirname(__FILE__) + '/redirection-examples/no_redirect.rfml' }

          it 'sets redirection to false' do
            test = subject.process
            step = test.steps.first
            expect(step.redirection).to eq('false')
          end
        end

        context 'redirection specified as true' do
          let(:file_name) { File.dirname(__FILE__) + '/redirection-examples/redirect.rfml' }

          it 'sets redirection to false' do
            test = subject.process
            step = test.steps.first
            expect(step.redirection).to eq('true')
          end
        end

        context 'redirection specified as not true or false' do
          let(:file_name) { File.dirname(__FILE__) + '/redirection-examples/wrong_redirect.rfml' }
          let(:redirect_line_no) { 5 }

          it 'creates an error with the correct line number' do
            test = subject.process
            expect(test.errors[redirect_line_no]).to_not be_nil
          end
        end
      end
    end
  end
end
