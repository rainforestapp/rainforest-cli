#!/usr/bin/env ruby
# frozen_string_literal: true
require 'csv'

module RainforestCli
  class CSVImporter
    BATCH_SIZE = 50

    def initialize(options)
      @generator_name = options.import_name
      @file = options.import_file_name
      @overwrite_variable = options.overwrite_variable
    end

    def row_data columns, values
      Hash[(columns.map {|c| c['id']}).zip(values)]
    end

    def import
      rows = []
      CSV.foreach(@file, encoding: 'windows-1251:utf-8') do |row|
        rows << row
      end

      # Create the generator
      columns = rows.shift.map do |column|
        column.downcase.strip.gsub(/\s/, '_')
      end
      raise 'Invalid schema in CSV. You must include headers in first row.' if !columns

      if @overwrite_variable
        puts 'Checking for existing tabular variables.'
        generators = http_client.get('/generators')
        generator = generators.find { |g| g['name'] == @generator_name }

        if generator
          puts 'Existing tabular variable found. Deleting old data.'
          response = http_client.delete("/generators/#{generator['id']}")
          raise "Error deleting old tabular variable: #{response['error']}" if response['error']
        end
      end

      puts 'Creating new tabular variable'
      response = http_client.post(
        '/generators',
        { name: @generator_name, description: @generator_name, columns: columns },
        { retries_on_failures: true },
      )
      raise "Error creating tabular variable: #{response['error']}" if response['error']
      puts "\t[OK]"

      columns = response['columns']
      generator_id = response['id']
      data = rows.map { |row| row_data(columns, row) }

      puts 'Uploading data...'
      row_count = (1.0 * data.count / BATCH_SIZE).ceil
      p = ProgressBar.create(title: 'Rows', total: row_count, format: '%a %B %p%% %t')

      data.each_slice(BATCH_SIZE) do |data_slice|
        response = http_client.post(
          "/generators/#{generator_id}/rows/batch",
          { data: data_slice },
          { retries_on_failures: true },
        )
        # NOTE: Response for this endpoint will usually be an array representing all the rows created
        raise response['error'] if response.is_a?(Hash) && response['error']
        p.increment
      end
      puts 'Upload complete.'
    end

    private

    def http_client
      RainforestCli.http_client
    end
  end
end
