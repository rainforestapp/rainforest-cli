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
        HTTParty.post(Rainforest::Cli::API_URL + url, body: body.to_json, headers: request_headers)
      end

      def get url
        HTTParty.get(Rainforest::Cli::API_URL + url, headers: request_headers)
      end

      def delete url
        HTTParty.delete(Rainforest::Cli::API_URL + url, headers: request_headers)
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
        response = create_variable(name: @generator_name, description: @generator_name, columns: columns)
        raise "Error creating custom step variable: #{response.code}, #{response.body}" unless response.code == 201
        puts "\t[OK]"

        @columns = response['columns']
        @generator_id = response['id']

        puts "Uploading data..."
        p = ProgressBar.create(title: 'Rows', total: rows.count, format: '%a %B %p%% %t')

        # Insert the data
        Parallel.each(rows, in_processes: 16, finish: lambda { |item, i, result| p.increment }) do |row|
          response = post("/generators/#{@generator_id}/rows", {data: row_data(@columns, row)})
          raise response.to_json unless response.code == 201
        end
      end

      def create_variable(name:, description:, columns:)

        response = post "/generators", {name: name, description: description, columns: columns}

        if response.code == 400 and response["error"] == "Name is already in use"
          generators = get("/generators")
          id = generators.detect { |g| g["name"] == name }["id"]
          print "\nFound existing generator with name \"#{name}\". Replacing it."
          delete("/generators/#{id}")
          response = post "/generators", {name: name, description: description, columns: columns}
        end

        response
      end

      def request_headers
        {'Content-Type' => 'application/json', "CLIENT_TOKEN" => @token}
      end
    end
  end
end
