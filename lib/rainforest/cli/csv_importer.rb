#!/usr/bin/env ruby
require 'csv'
require 'httparty'
require 'parallel'
require 'ruby-progressbar'

module RainforestCli
  class CSVImporter
    attr_reader :client

    def initialize name, file, token
      @generator_name = name
      @file = file
      @client = HttpClient.new token: token
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
      raise 'Invalid schema in CSV. You must include headers in first row.' if not columns

      print "Creating custom step variable"
      response = client.post "/generators", { name: @generator_name, description: @generator_name, columns: columns }
      raise "Error creating custom step variable: #{response['error']}" if response['error']
      puts "\t[OK]"

      @columns = response['columns']
      @generator_id = response['id']

      puts "Uploading data..."
      p = ProgressBar.create(title: 'Rows', total: rows.count, format: '%a %B %p%% %t')

      # Insert the data
      Parallel.each(rows, in_threads: 16, finish: lambda { |item, i, result| p.increment }) do |row|
        response = client.post("/generators/#{@generator_id}/rows", {data: row_data(@columns, row)})
        raise response['error'] if response['error']
      end
    end
  end
end
