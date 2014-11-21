require "rainforest/cli/version"
require "rainforest/cli/options"
require "rainforest/cli/git_trigger"
require "rainforest/cli/csv_importer"
require "httparty"
require "json"
require "logger"

module Rainforest
  module Cli 
    API_URL = 'https://app.rainforestqa.com/api/1'.freeze
    
    def self.start(args)
      @options = OptionParser.new(args)

      unless @options.token
        logger.fatal "You must pass your API token using: --token TOKEN"
        exit 2
      end

      if @options.custom_url && @options.site_id.nil?
        logger.fatal "The site-id and custom-url options are both required."
        exit 2
      end

      if @options.import_file_name && @options.import_name
        unless File.exists?(@options.import_file_name)
          logger.fatal "Input file: #{@options.import_file_name} not found"
          exit 2
        end

        delete_generator(@options.import_name)
        CSVImporter.new(@options.import_name, @options.import_file_name, @options.token).import
      elsif @options.import_file_name || @options.import_name
        logger.fatal "You must pass both --import-variable-csv-file and --import-variable-name"
        exit 2
      end

      post_opts = {}

      if @options.git_trigger?
        logger.debug "Checking last git commit message:"
        commit_message = GitTrigger.last_commit_message
        logger.debug commit_message

        # Show some messages to users about tests/tags being overriden
        unless @options.tags.empty?
          logger.warn "Specified tags are ignored when using --git-trigger"
        else
          logger.warn "Specified tests are ignored when using --git-trigger"
        end

        if GitTrigger.git_trigger_should_run?(commit_message)
          tags = GitTrigger.extract_hashtags(commit_message)
          if tags.empty?
            logger.error "Triggered via git, but no hashtags detected. Please use commit message format:"
            logger.error "\t'some message. @rainforest #tag1 #tag2"
            exit 2
          else
            post_opts[:tags] = [tags.join(',')]
          end
        else
          logger.info "Not triggering as @rainforest was not mentioned in last commit message."
          exit 0
        end
      else
        # Not using git_trigger, so look for the 
        if !@options.tags.empty?
          post_opts[:tags] = @options.tags
        else
          post_opts[:tests] = @options.tests
        end
      end

      post_opts[:conflict] = @options.conflict if @options.conflict
      post_opts[:browsers] = @options.browsers if @options.browsers
      post_opts[:site_id] = @options.site_id if @options.site_id
      post_opts[:gem_version] = Rainforest::Cli::VERSION

      post_opts[:environment_id] = get_environment_id(@options.custom_url) if @options.custom_url

      logger.debug "POST options: #{post_opts.inspect}"
      logger.info "Issuing run"

      response = post(API_URL + '/runs', post_opts)

      if response['error']
        if !response['error'].index("There are no steps in this run").nil?
          logger.error "Error starting your run: #{response['error']}"
        else
          logger.fatal "Error starting your run: #{response['error']}"
          exit 1
        end
      end

      run_id = response["id"]
      running = true

      return unless @options.foreground?

      while running 
        Kernel.sleep 5
        response = get "#{API_URL}/runs/#{run_id}?gem_version=#{Rainforest::Cli::VERSION}"
        if response 
          if %w(queued in_progress sending_webhook waiting_for_callback).include?(response["state"])
            logger.info "Run #{run_id} is #{response['state']} and is #{response['current_progress']['percent']}% complete"
            running = false if response["result"] == 'failed' && @options.failfast?
          else
            logger.info "Run #{run_id} is now #{response["state"]} and has #{response["result"]}"
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
      generator = list_generators.find {|g| g['generator_type'] == 'tabular' && g['name'] == name }
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
        logger.fatal "The custom URL is invalid"
        exit 2
      end

      env_post_body = { name: 'temporary-env-for-custom-url-via-CLI', url: url }
      environment = post("#{API_URL}/environments", env_post_body)

      if environment['error']
        # I am talking about a URL here because the environments are pretty
        # much hidden from clients so far.
        logger.fatal "Error creating the ad-hoc URL: #{environment['error']}"
        exit 1
      end

      return environment['id']
    end

    def self.logger
      @logger ||= Logger.new(STDOUT)
    end
  end
end
