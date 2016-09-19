#!/usr/bin/env ruby
# frozen_string_literal: true
require 'csv'
require 'parallel'
require 'ruby-progressbar'

module RainforestCli
  class CSVImporter
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
        {name: column.downcase.strip.gsub(/\s/, '_')}
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

      print 'Creating new tabular variable'
      response = http_client.post '/generators', { name: @generator_name, description: @generator_name, columns: columns }
      raise "Error creating tabular variable: #{response['error']}" if response['error']
      puts "\t[OK]"

      @columns = response['columns']
      @generator_id = response['id']

      puts 'Uploading data...'
      p = ProgressBar.create(title: 'Rows', total: rows.count, format: '%a %B %p%% %t')

      # Insert the data
      Parallel.each(rows, in_threads: RainforestCli::THREADS, finish: lambda { |_item, _i, _result| p.increment }) do |row|
        response = http_client.post("/generators/#{@generator_id}/rows", {data: row_data(@columns, row)})
        raise response['error'] if response['error']
      end
    end

    private

    def http_client
      RainforestCli.http_client
    end
  end
end
