# frozen_string_literal: true
describe RainforestCli::CSVImporter do
  let(:csv_file) { "#{File.dirname(__FILE__)}/../fixtures/variables.txt" }
  let(:http_client) { instance_double('RainforestCli::HttpClient') }

  before do
    allow(RainforestCli).to receive(:http_client).and_return(http_client)
  end

  describe '.import' do
    subject { described_class.new(options) }

    let(:options) do
      instance_double(
      'RainforestCli::Options',
      {
        import_name: 'variables',
        import_file_name: csv_file,
        overwrite_variable: overwrite_variable,
        single_use_tabular_variable: single_use_tabular_variable,
      })
    end

    let(:overwrite_variable) { nil }
    let(:single_use_tabular_variable) { nil }

    let(:columns) { %w(email pass) }
    let(:generator_id) { 12345 }
    let(:existing_generators) { [] }

    let(:success_response) do
      {
        'id' => generator_id,
        'columns' => columns.each_with_index.map { |col, i| { 'id' => i, 'name' => col } },
      }
    end

    before do
      allow(http_client).to receive(:get).with('/generators').and_return(existing_generators)
    end

    shared_examples 'it properly uploads variables' do
      it 'makes the proper API interactions' do
        expect(http_client).to receive(:post)
                                .with('/generators', {
                                        name: 'variables',
                                        description: 'variables',
                                        columns: columns,
                                        single_use: (single_use_tabular_variable || false),
                                      }, retries_on_failures: true)
                                .and_return success_response

        expect(http_client).to receive(:post)
                                .with("/generators/#{generator_id}/rows/batch", {
                                        data: [
                                          {
                                            0 => 'russ@rainforestqa.com',
                                            1 => 'abc123',
                                          },
                                          {
                                            0 => 'bob@example.com',
                                            1 => 'hunter2',
                                          },
                                        ],
                                      }, retries_on_failures: true).and_return({})
        subject.import
      end
    end

    it_behaves_like 'it properly uploads variables'

    context 'tabular variable with given name already exists' do
      let(:existing_generators) do
        [
          {
            'id' => 98765,
            'name' => 'existing',
          },
          {
            'id' => generator_id,
            'name' => 'variables',
          },
        ]
      end

      before do
        expect(http_client).to_not receive(:delete)
      end

      it_behaves_like 'it properly uploads variables'
    end

    context 'with variable overwriting' do
      let(:overwrite_variable) { true }

      it_behaves_like 'it properly uploads variables'

      context 'tabular variable with given name already exists' do
        let(:existing_generators) do
          [
            {
              'id' => 98765,
              'name' => 'existing',
            },
            {
              'id' => generator_id,
              'name' => 'variables',
            },
          ]
        end

        before do
          expect(http_client).to receive(:delete).with("/generators/#{generator_id}").and_return({})
        end

        it_behaves_like 'it properly uploads variables'
      end
    end

    context 'with single-use flag' do
      let(:single_use_tabular_variable) { true }

      it_behaves_like 'it properly uploads variables'
    end
  end
end
