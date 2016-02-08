describe RainforestCli::TestFiles do
  describe '#initialize' do
    before do
      allow(Dir).to receive(:mkdir)
    end

    it 'uses the default file folder' do
      expect(described_class.new.test_folder).to eq(described_class::DEFAULT_TEST_FOLDER)
    end

    context 'when test folder name is supplied' do
      let(:folder_name) { './foo' }

      it 'creates the supplied folder if file does not exist' do
        expect(Dir).to receive(:mkdir).with(folder_name)
        described_class.new(folder_name)
      end
    end
  end
end
