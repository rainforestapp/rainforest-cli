# frozen_string_literal: true
describe RainforestCli::RemoteTests do
  subject { described_class.new('api_token') }

  describe '#primary_key_dictionary' do
    Test = Struct.new(:rfml_id, :id)

    let(:test1) { { 'rfml_id' => 'foo', 'id' => 123 } }
    let(:test2) { { 'rfml_id' => 'bar', 'id' => 456 } }
    let(:test3) { { 'rfml_id' => 'baz', 'id' => 789 } }
    let(:tests) { [test1, test2, test3] }

    before do
      allow_any_instance_of(RainforestCli::HttpClient).to receive(:get)
        .with('/tests/rfml_ids').and_return(tests)
    end

    it "correctly formats the dictionary's keys and values" do
      expect(subject.primary_key_dictionary)
        .to include({
                      test1['rfml_id'] => test1['id'],
                      test2['rfml_id'] => test2['id'],
                      test3['rfml_id'] => test3['id']
                    })
    end

    context 'no api token set' do
      subject { described_class.new }

      it 'does not make an API call but returns an empty dictionary' do
        expect_any_instance_of(RainforestCli::HttpClient).to_not receive(:get)
        dictionary = subject.primary_key_dictionary
        expect(dictionary).to be_a(Hash)
        expect(dictionary).to eq({})
      end
    end
  end
end
