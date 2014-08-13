require "rainforest/cli/version"
require "rainforest/cli/options"
require "rainforest/cli/csv_importer"
require "httparty"
require "json"

module Rainforest
  module Cli 
    API_URL = 'https://app.rainforestqa.com/api/1'.freeze
    
    def self.start(args)
      @options = OptionParser.new(args)
      
      if @options.import_file_name && @options.import_name
        unless File.exists?(@options.import_file_name)
          puts "Input file: #{@options.import_file_name} not found"
          exit
        end
        
        delete_generator(@options.import_name)
        CSVImporter.new(@options.import_name, @options.import_file_name, @options.token).import
      elsif @options.import_file_name || @options.import_name
        puts "You must pass both --import-variable-csv-file and --import-variable-name"
        exit
      end

      post_opts = {}
      if !@options.tags.empty?
        post_opts[:tags] = @options.tags
      else
        post_opts[:tests] = @options.tests
      end

      post_opts[:conflict] = @options.conflict if @options.conflict
      post_opts[:browsers] = @options.browsers if @options.browsers
      post_opts[:gem_version] = Rainforest::Cli::VERSION

      puts "Issuing run"

      response = post(API_URL + '/runs', post_opts)

      if response['error']
        puts "Error starting your run: #{response['error']}"
        exit
      end

      run_id = response["id"]
      running = true

      return unless @options.foreground?

      while running 
        sleep 5
        response = get "#{API_URL}/runs/#{run_id}?gem_version=#{Rainforest::Cli::VERSION}"
        if %w(queued in_progress sending_webhook waiting_for_callback).include?(response["state"])
          puts "Run #{run_id} is #{response['state']} and is #{response['current_progress']['percent']}% complete"
          running = false if response["result"] == 'failed' && @options.failfast?
        else
          puts "Run #{run_id} is now #{response["state"]} and has #{response["result"]}"
          running = false
        end
      end

      if response["result"] != "passed"
        exit 1
      end
    end
    
    def self.list_generators
      get("#{API_URL}/generators")
    end
    
    def self.delete_generator(name)
      generator = list_generators.select {|g| g['type'] == 'custom' && g['key'] == name }.first
      delete("#{API_URL}/generators/#{generator['id']}") if generator
    end

    def self.delete(url, body = {})
      response = HTTParty.delete url, {
        body: body, 
        headers: {"CLIENT_TOKEN" => @options.token}
      }
      
      JSON.parse(response.body)
    end

    def self.post(url, body = {})
      response = HTTParty.post url, {
        body: body, 
        headers: {"CLIENT_TOKEN" => @options.token}
      }

      JSON.parse(response.body)
    end

    def self.get(url, body = {})
      response = HTTParty.get url, {
        body: body, 
        headers: {"CLIENT_TOKEN" => @options.token}
      }

      JSON.parse(response.body)
    end
  end
end
