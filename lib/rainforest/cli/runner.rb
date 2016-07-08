# frozen_string_literal: true
module RainforestCli
  class Runner
    attr_reader :options, :client

    def initialize(options)
      @options = options
      @client = HttpClient.new token: options.token
    end

    def run
      if options.import_file_name && options.import_name
        delete_generator(options.import_name)
        CSVImporter.new(options.import_name, options.import_file_name, options.token).import
      end

      post_opts = make_create_run_options

      logger.debug "POST options: #{post_opts.inspect}"
      logger.info 'Issuing run'

      response = client.post('/runs', post_opts, retries_on_failures: true)

      if response['error']
        logger.fatal "Error starting your run: #{response['error']}"
        exit 1
      end

      if options.foreground?
        run_id = response.fetch('id')
        wait_for_run_completion(run_id)
        if options.junit?
          generate_junit_output(run_id, options.junit)
        end
      else
        true
      end
    end

    def wait_for_run_completion(run_id)
      running = true
      while running
        Kernel.sleep 5
        response = client.get("/runs/#{run_id}", {}, retries_on_failures: true)
        if response
          state_details = response.fetch('state_details')
          unless state_details.fetch('is_final_state')
            state, current_progress = response.values_at('state', 'current_progress')
            logger.info "Run #{run_id} is #{state} and is #{current_progress['percent']}% complete"
            running = false if response['result'] == 'failed' && options.failfast?
          else
            logger.info "Run #{run_id} is now #{response["state"]} and has #{response["result"]}"
            running = false
          end
        end
      end

      if response['frontend_url']
        logger.info "The detailed results are available at #{response['frontend_url']}"
      end

      if response['result'] != 'passed'
        exit 1
      end
    end

    def generate_junit_output(run_id, output_filename)
      run = client.get("/runs/#{run_id}.json")

      if run['error']
        logger.fatal "Error retrieving results for your run: #{run['error']}"
        exit 1
      end

      if run.has_key?('total_tests') and run['total_tests'] != 0
        tests = client.get("/runs/#{run_id}/tests.json?page_size=#{run['total_tests']}")

        if tests['error']
          logger.fatal "Error retrieving test details for your run: #{tests['error']}"
          exit 1
        end

        outputter = JunitOutputter.new(options.token, run, tests).parse
      end

      File.open(options.output_filename, 'w') { |file| outputter.output(file) }
    end

    def make_create_run_options
      post_opts = {}
      if options.git_trigger?
        logger.debug 'Checking last git commit message:'
        commit_message = GitTrigger.last_commit_message
        logger.debug commit_message

        # Show some messages to users about tests/tags being overriden
        unless options.tags.empty?
          logger.warn 'Specified tags are ignored when using --git-trigger'
        else
          logger.warn 'Specified tests are ignored when using --git-trigger'
        end

        if GitTrigger.git_trigger_should_run?(commit_message)
          tags = GitTrigger.extract_hashtags(commit_message)
          if tags.empty?
            logger.error 'Triggered via git, but no hashtags detected. Please use commit message format:'
            logger.error "\t'some message. @rainforest #tag1 #tag2"
            exit 2
          else
            post_opts[:tags] = [tags.join(',')]
          end
        else
          logger.info 'Not triggering as @rainforest was not mentioned in last commit message.'
          exit 0
        end
      else
        # Not using git_trigger, so look for the
        if !options.tags.empty?
          post_opts[:tags] = options.tags
        elsif !options.folder.nil?
          post_opts[:smart_folder_id] = @options.folder.to_i
        else
          post_opts[:tests] = options.tests
        end
      end

      app_source_url = options.app_source_url ? upload_app(options.app_source_url) : options.app_source_url

      post_opts[:app_source_url] = app_source_url if app_source_url
      post_opts[:crowd] = options.crowd if options.crowd
      post_opts[:conflict] = options.conflict if options.conflict
      post_opts[:browsers] = options.browsers if options.browsers
      post_opts[:site_id] = options.site_id if options.site_id
      post_opts[:description] = options.description if options.description

      if options.custom_url
        post_opts[:environment_id] = get_environment_id(options.custom_url)
      elsif options.environment_id
        post_opts[:environment_id] = options.environment_id
      end

      post_opts
    end

    def logger
      RainforestCli.logger
    end

    def list_generators
      client.get('/generators')
    end

    def delete_generator(name)
      generator = list_generators.find {|g| g['generator_type'] == 'tabular' && g['name'] == name }
      client.delete("generators/#{generator['id']}") if generator
    end

    def url_valid?(url)
      return false unless URI::regexp === url

      uri = URI.parse(url)
      %w(http https).include?(uri.scheme)
    end

    def get_environment_id url
      unless url_valid?(url)
        logger.fatal 'The custom URL is invalid'
        exit 2
      end

      env_post_body = { name: 'temporary-env-for-custom-url-via-CLI', url: url }
      environment = client.post('/environments', env_post_body)

      if environment['error']
        # I am talking about a URL here because the environments are pretty
        # much hidden from clients so far.
        logger.fatal "Error creating the ad-hoc URL: #{environment['error']}"
        exit 1
      end

      return environment['id']
    end

    def upload_app(app_source_url)
      return app_source_url if url_valid?(app_source_url)

      unless File.exist?(app_source_url)
        logger.fatal "App source file: #{app_source_url} not found"
        exit 1
      end

      unless File.extname(app_source_url) == '.ipa'
        logger.fatal "Invalid app source file: #{app_source_url}"
        exit 1
      end

      url = client.get('/uploads', {}, retries_on_failures: true)
      unless url
        logger.fatal "Failed to upload file #{app_source_url}. Please, check your API token."
        exit 1
      end
      data = File.read(app_source_url)
      logger.info 'Uploading app source file, this operation may take few minutes...'
      response_code = upload_file(url, data)

      if response_code != '200'
        logger.fatal "Failed to upload file #{app_source_url}"
        exit 1
      else
        logger.info 'Upload completed.'
        return url['path']
      end
    end

    def upload_file(url, body)
      begin
        # using Net::HTTP because HTTParty do not support large files
        http = Net::HTTP.new(url['host'], url['port'])
        http.use_ssl = true
        http.verify_mode = OpenSSL::SSL::VERIFY_NONE

        request = Net::HTTP::Put.new(url['uri'])
        request['Content-Type'] = 'application/zip'
        request.body = body

        response = http.request(request)
        return response.code
      rescue Http::Exceptions::HttpException => e
        logger.fatal "Error uploading the app source file: #{e.message}"
        exit 1
      end
    end
  end
end
