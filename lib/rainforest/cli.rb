require "rainforest/cli/version"
require "rainforest/cli/options"
require "rainforest"
require "json"

module RainforestCli
  def self.start(args)
    @options = OptionParser.new(args)
    Rainforest.api_key = @options.token

    post_opts = {}
    if !@options.tags.empty?
      post_opts[:tags] = @options.tags
    else
      post_opts[:tests] = @options.tests
    end

    post_opts[:conflict] = @options.conflict if @options.conflict
    post_opts[:browsers] = @options.browsers if @options.browsers
    post_opts[:gem_version] = RainforestCli::VERSION

    puts "Issuing run"

    response = Rainforest::Run.create(post_opts)

    if response.respond_to? :error
      puts "Error starting your run: #{response.error}"
      exit
    end

    run_id = response.id
    running = true

    return unless @options.foreground?

    while running
      sleep 5
      response = Rainforest::Run.retrieve(run_id)
      if %w(queued in_progress sending_webhook waiting_for_callback).include?(response.state)
        puts "Run #{run_id} is #{response.state} and is #{response.current_progress.percent}% complete"
        running = false if response.result == 'failed' && @options.failfast?
      else
        puts "Run #{run_id} is now #{response.state} and has #{response.result}"
        running = false
      end
    end

    if response.result != "passed"
      exit 1
    end
  end
end
