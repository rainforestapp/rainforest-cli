# frozen_string_literal: true
describe RainforestCli::TestParser::Step do
  subject { described_class.new }
  let(:screenshot_string) do
    'Picture 1: {{ file.screenshot(./foo) }}. Picture 2: {{file.screenshot(bar/baz)   }}'
  end

  let(:download_string) do
    'Download 1: {{ file.download(./foo) }}. Download 2: {{file.download(bar/baz)   }}'
  end

  let(:parameterized_screenshot_string) do
    'Picture 1: {{ file.screenshot(./foo) }}. Picture 2: {{file.screenshot(123, signat)   }}'
  end

  let(:parameterized_download_string) do
    'Download 1: {{ file.download(./foo) }}. Download 2: {{file.download(123, signat, baz.txt)   }}'
  end

  shared_examples 'a method that detects step variables' do |att|
    it 'correctly detects a screenshot step variable' do
      subject[att] = screenshot_string
      expect(subject.send(:"uploadable_in_#{att}"))
        .to eq([['screenshot', './foo'], ['screenshot', 'bar/baz']])
    end

    it 'correctly detects a download step variable' do
      subject[att] = download_string
      expect(subject.send(:"uploadable_in_#{att}"))
        .to eq([['download', './foo'], ['download', 'bar/baz']])
    end

    it 'does not detect a parameterized screenshot step variable' do
      subject[att] = parameterized_screenshot_string
      expect(subject.send(:"uploadable_in_#{att}"))
        .to eq([['screenshot', './foo']])
    end

    it 'does not detect a parameterized download step variable' do
      subject[att] = parameterized_download_string
      expect(subject.send(:"uploadable_in_#{att}"))
        .to eq([['download', './foo']])
    end
  end

  describe '#uploadable_in_action' do
    it_behaves_like 'a method that detects step variables', :action
  end

  describe '#uploadable_in_response' do
    it_behaves_like 'a method that detects step variables', :response
  end

  describe '#has_uploadable_files?' do
    let(:action) { 'Regular action' }
    let(:response) { 'Regular response' }
    subject { described_class.new(action, response).has_uploadable_files? }

    context 'with no uploadables' do
      it { is_expected.to be(false) }
    end

    context 'uploadable in action' do
      let(:action) { 'Action with uploadable {{ file.download(/foo) }}' }
      it { is_expected.to be(true) }
    end

    context 'uploadable in response' do
      let(:response) { 'Response with uploadable {{ file.download(/foo) }}' }
      it { is_expected.to be(true) }
    end
  end
end
