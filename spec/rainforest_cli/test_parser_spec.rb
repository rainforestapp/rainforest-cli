# frozen_string_literal: true
describe RainforestCli::TestParser do
  describe RainforestCli::TestParser::Step do
    describe '#to_element' do
      subject { described_class.new('action', 'response', 'redirect').to_element }
      its([:redirection]) { is_expected.to eq('redirect') }

      context 'with no redirect' do
        subject { described_class.new('action', 'response', nil).to_element }
        its([:redirection]) { is_expected.to eq('true') }
      end
    end
  end

  describe RainforestCli::TestParser::EmbeddedTest do
    describe '#to_element' do
      let(:primary_key_id) { 123 }
      subject { described_class.new.to_element(primary_key_id) }

      context 'with no redirect' do
        its([:redirection]) { is_expected.to eq('true') }
      end
    end
  end

  describe RainforestCli::TestParser::Parser do
    subject { described_class.new(file_name) }

    describe '#initialize' do
      let(:file_name) { './spec/rainforest-example/example_test.rfml' }

      it 'expands the file name path' do
        test = subject.instance_variable_get(:@test)
        expect(test.file_name).to eq(File.expand_path(file_name))
      end
    end

    describe '#process' do
      context 'redirection' do
        context 'step' do
          context 'no redirection specified' do
            let(:file_name) { File.dirname(__FILE__) + '/../rainforest-example/example_test.rfml' }
            it { has_redirection_value(subject.process, 'true') }
          end

          context 'redirection specified as true' do
            let(:file_name) { File.dirname(__FILE__) + '/../redirection-examples/redirect.rfml' }
            it { has_redirection_value(subject.process, 'true') }
          end

          context 'redirection specified as not true or false' do
            let(:file_name) { File.dirname(__FILE__) + '/../redirection-examples/wrong_redirect.rfml' }
            let(:redirect_line_no) { 7 }
            it { has_parsing_error(subject.process, redirect_line_no) }
          end
        end

        context 'embedded test' do
          context 'redirection specified as false' do
            let(:file_name) { File.dirname(__FILE__) + '/../redirection-examples/no_redirect_embedded.rfml' }
            it { has_redirection_value(subject.process, 'false') }
          end

          context 'redirection specified as true' do
            let(:file_name) { File.dirname(__FILE__) + '/../redirection-examples/redirect_embedded.rfml' }
            it { has_redirection_value(subject.process, 'true') }
          end

          context 'redirection specified as not true or false' do
            let(:file_name) { File.dirname(__FILE__) + '/../redirection-examples/wrong_redirect_embedded.rfml' }
            let(:redirect_line_no) { 8 }
            it { has_parsing_error(subject.process, redirect_line_no) }
          end
        end

        context 'poor syntax' do
          context 'extra empty lines' do
            let(:file_name) { File.dirname(__FILE__) + '/../redirection-examples/wrong_redirect_spacing.rfml' }
            let(:empty_line_no) { 8 }
            it { has_parsing_error(subject.process, empty_line_no) }
          end
        end

        def has_redirection_value(test, value)
          step = test.steps.last
          expect(step.redirection).to eq(value)
        end

        def has_parsing_error(test, line_no)
          expect(test.errors[line_no]).to_not be_nil
        end
      end
    end
  end
end
