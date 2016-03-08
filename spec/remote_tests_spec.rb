# frozen_string_literal: true
describe RainforestCli::RemoteTests do
  subject { described_class.new('api_token') }

  describe '#primary_key_dictionary' do
    Test = Struct.new(:rfml_id, :id)

    let(:test1) { Test.new('foo', 123) }
    let(:test2) { Test.new('bar', 456) }
    let(:test3) { Test.new('baz', 789) }
    let(:tests) { [test1, test2, test3] }

    before do
      allow(subject).to receive(:tests).and_return(tests)
    end

    it "correctly formats the dictionary's keys and values" do
      expect(subject.primary_key_dictionary)
        .to include({'foo' => test1.id, 'bar' => test2.id, 'baz' => test3.id})
    end
  end
end
