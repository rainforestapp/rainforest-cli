#!/usr/bin/env ruby
require 'csv'
require 'httparty'
require 'parallel'
require 'ruby-progressbar'

module Rainforest
  module Cli
    class CSVImporter
      def initialize name, file, token
        @generator_name = name
        @file = file
        @token = token
      end

      def post url, body
        HTTParty.post(Rainforest::Cli::API_URL + url, body: body.to_json, headers: {'Content-Type' => 'application/json', "CLIENT_TOKEN" => @token})
      end

      def row_data columns, values
        Hash[(columns.map {|c| c['id']}).zip(values)]
      end

      def import
        rows = []
        puts "Loading data"
        CSV.foreach(@file, encoding: 'windows-1251:utf-8') do |row|
          rows << row
        end

        # Create the generator
        puts "Parsing rows"
        columns = rows.shift.map do |column|
          {name: column.downcase.strip.gsub(/\s/, '_')}
        end
        raise 'Invalid generator schema' if not columns

        puts "Creating generator"
        response = post "/generators", { name: @generator_name, description: @generator_name, columns: columns }
        raise "Error creating generator: #{response.code}, #{response.body}" unless response.code == 201

        @columns = response['columns']
        @generator_id = response['id']

        puts "Uploading data"
        p = ProgressBar.create(title: 'Rows', total: rows.count, format: '%a %B %p%% %t')

        # Insert the data
        Parallel.each(rows, in_processes: 16, finish: lambda { |item, i, result| p.increment }) do |row|
          response = post("/generators/#{@generator_id}/rows", {data: row_data(@columns, row)})
          raise response.to_json unless response.code == 201
        end
      end
    end
  end
end