# frozen_string_literal: true
describe RainforestCli::CSVImporter do
  let(:csv_file) { "#{File.dirname(__FILE__)}/../fixtures/variables.txt" }
  let(:http_client) { instance_double('RainforestCli::HttpClient') }

  before do
    allow(RainforestCli).to receive(:http_client).and_return(http_client)
  end

  describe '.import' do
    let(:options) { instance_double('RainforestCli::Options', import_name: 'variables', import_file_name: csv_file) }
    subject { described_class.new(options) }
    let(:columns) { %w(email pass) }

    let(:success_response) do
      {
        'id' => 12345,
        'columns' => columns.each_with_index.map { |col, i| { 'id' => i, 'name' => col } },
      }
    end

    it 'should post the schema to the generators API' do
      expect(http_client).to receive(:post)
                              .with('/generators', {
                                      name: 'variables',
                                      description: 'variables',
                                      columns: columns.map {|col| { name: col } },
                                    })
                              .and_return success_response

      expect(http_client).to receive(:post)
                              .with('/generators/12345/rows', {
                                      data: {
                                        0 => 'russ@rainforestqa.com',
                                        1 => 'abc123',
                                      },
                                    }).and_return({})

      expect(http_client).to receive(:post)
                              .with('/generators/12345/rows', {
                                      data: {
                                        0 => 'bob@example.com',
                                        1 => 'hunter2',
                                      },
                                    }).and_return({})

      subject.import
    end
  end
end
