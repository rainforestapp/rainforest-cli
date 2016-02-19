# frozen_string_literal: true
describe RainforestCli::CSVImporter do
  let(:csv_file) { "#{File.dirname(__FILE__)}/fixtures/variables.txt" }

  describe '.import' do
    subject { described_class.new('variables', csv_file, 'abc123') }
    let(:progressbar_mock) { double('ProgressBar') }

    before do
      # suppress output in terminal
      allow_any_instance_of(described_class).to receive(:print)
      allow_any_instance_of(described_class).to receive(:puts)
      allow(ProgressBar).to receive(:create).and_return(progressbar_mock)
      allow(progressbar_mock).to receive(:increment)
    end

    let(:http_client) { double }
    let(:columns) { %w(email pass) }

    let(:success_response) do
      {
        'id' => 12345,
        'columns' => columns.each_with_index.map { |col, i| { 'id' => i, 'name' => col } }
      }
    end

    before do
      RainforestCli::HttpClient.stub(:new).and_return(http_client)
    end

    it 'should post the schema to the generators API' do
      expect(http_client).to receive(:post)
                              .with('/generators', {
                                      name: 'variables',
                                      description: 'variables',
                                      columns: columns.map {|col| { name: col } }
                                    })
                              .and_return success_response

      expect(http_client).to receive(:post)
                              .with('/generators/12345/rows', {
                                      data: {
                                        0 => 'russ@rainforestqa.com',
                                        1 => 'abc123'
                                      }
                                    }).and_return({})

      expect(http_client).to receive(:post)
                              .with('/generators/12345/rows', {
                                      data: {
                                        0 => 'bob@example.com',
                                        1 => 'hunter2'
                                      }
                                    }).and_return({})

      subject.import
    end
  end
end
