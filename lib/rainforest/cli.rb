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

      if @options.custom_url && @options.site_id.nil?
        puts "The site-id and custom-url options work together, you need both of them."
        exit 1
      end

      if @options.import_file_name && @options.import_name
        unless File.exists?(@options.import_file_name)
          puts "Input file: #{@options.import_file_name} not found"
          exit 2
        end

        delete_generator(@options.import_name)
        CSVImporter.new(@options.import_name, @options.import_file_name, @options.token).import
      elsif @options.import_file_name || @options.import_name
        puts "You must pass both --import-variable-csv-file and --import-variable-name"
        exit 2
      end

      post_opts = {}
      if !@options.tags.empty?
        post_opts[:tags] = @options.tags
      else
        post_opts[:tests] = @options.tests
      end

      post_opts[:conflict] = @options.conflict if @options.conflict
      post_opts[:browsers] = @options.browsers if @options.browsers
      post_opts[:site_id] = @options.site_id if @options.site_id
      post_opts[:gem_version] = Rainforest::Cli::VERSION

      post_opts[:environment_id] = get_environment_id(@options.custom_url) if @options.custom_url

      puts "Issuing run"

      response = post(API_URL + '/runs', post_opts)

      if response['error']
        puts "Error starting your run: #{response['error']}"
        exit 1
      end

      run_id = response["id"]
      running = true

      return unless @options.foreground?

      while running 
        Kernel.sleep 5
        response = get "#{API_URL}/runs/#{run_id}?gem_version=#{Rainforest::Cli::VERSION}"
        if response 
          if %w(queued in_progress sending_webhook waiting_for_callback).include?(response["state"])
            puts "Run #{run_id} is #{response['state']} and is #{response['current_progress']['percent']}% complete"
            running = false if response["result"] == 'failed' && @options.failfast?
          else
            puts "Run #{run_id} is now #{response["state"]} and has #{response["result"]}"
            running = false
          end
        end
      end

      if response["result"] != "passed"
        exit 1
      end
      true
    end

    def self.list_generators
      get("#{API_URL}/generators")
    end

    def self.delete_generator(name)
      generator = list_generators.find {|g| g['type'] == 'custom' && g['key'] == name }
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

      if response.code == 200
        JSON.parse(response.body)
      else
        nil
      end
    end

    def self.get_environment_id url
      begin
        URI.parse(url)
      rescue URI::InvalidURIError
        puts "The custom URL is invalid"
        exit 1
      end
      env_post_body = { name: 'temporary-env-for-custom-url-via-CLI', url: url }
      environment = post("#{API_URL}/environments", env_post_body)
      return environment['id']
    end
  end
end
