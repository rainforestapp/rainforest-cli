# frozen_string_literal: true
describe RainforestCli::Commands do
  subject { described_class.new }

  describe '#initialize' do
    it 'works without giving a block' do
      expect { described_class.new }.to_not raise_error
    end

    it 'adds commands to list' do
      commands = described_class.new do |c|
        c.add('foo', 'bar') { 'baz' }
        c.add('blah', 'blah') { 'blah' }
      end

      expect(commands.commands.length).to eq(2)
    end
  end

  describe '#add' do
    let(:first_command) { subject.commands.first }

    it 'properly adds a command to the list of commands' do
      subject.add('foo', 'bar') { 'baz' }
      expect(first_command.name).to eq('foo')
      expect(first_command.description).to eq('bar')
      expect(first_command.block.call).to eq('baz')
    end
  end

  describe '#call' do
    let(:service_object) { double }

    before do
      subject.add('my_command', 'desc') { service_object.my_method }
    end

    it "calls the corresponding command's block" do
      expect(service_object).to receive(:my_method)
      subject.call('my_command')
    end

    it 'raises a fatal error if the command does not exist' do
      expect { subject.call('my_other_command') }.to raise_error(SystemExit)
    end
  end
end
