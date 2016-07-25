require 'fileutils'

#frozen_string_literal: true
module RainforestCli
  class Reporter
    attr_reader :client
    attr_writer :run_id

    def initialize(options)
      @options = options
      @client = HttpClient.new token: options.token
      @run_id = options.run_id
      @output_filename = options.junit_file
    end

    def report
      if @run_id == nil
        logger.fatal "Reporter needs a valid run_id to report on"
      else
        logger.info "Generating JUNIT report for #{@run_id} : #{@output_filename}"
      end

      run = client.get("/runs/#{@run_id}.json")

      if run == nil
        logger.fatal "Non 200 code recieved"
        exit 1
      end

      if run['error']
        logger.fatal "Error retrieving results for your run: #{run['error']}"
        exit 1
      end

      if run.has_key?('total_tests') and run['total_tests'] != 0
        tests = client.get("/runs/#{@run_id}/tests.json?page_size=#{run['total_tests']}")

        if tests == nil
          logger.fatal "Non 200 code recieved"
          exit 1
        end

        if tests.kind_of?(Hash) and tests['error'] # if this had worked tests would be an array
          logger.fatal "Error retrieving test details for your run: #{tests['error']}"
          exit 1
        end

        outputter = JunitOutputter.new(@options.token, run, tests)
        outputter.parse
      end

      unless File.directory?(File.dirname(@output_filename))
        FileUtils.mkdir_p(File.dirname(@output_filename))
      end

      File.open(@output_filename, 'w') { |file| outputter.output(file) }
    end

    def logger
      RainforestCli.logger
    end

  end
end
