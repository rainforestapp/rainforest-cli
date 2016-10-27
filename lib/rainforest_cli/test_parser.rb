# frozen_string_literal: true

module RainforestCli::TestParser
  require 'rainforest_cli/test_parser/test'
  require 'rainforest_cli/test_parser/step'
  require 'rainforest_cli/test_parser/embedded_test'

  class Error < Struct.new(:line, :reason)
    def to_s
      "Line #{line}: #{reason}"
    end
  end

  class Parser
    attr_reader :steps, :errors, :text

    def initialize(file_name)
      @text = File.read(file_name).to_s

      @test = Test.new
      @test.file_name = File.expand_path(file_name)
      @test.description = ''
      @test.steps = []
      @test.errors = {}
      @test.tags = []
      @test.browsers = []
    end

    TEST_DATA_FIELDS = [:start_uri, :title, :site_id].freeze
    STEP_DATA_FIELDS = [:redirect].freeze
    CSV_FIELDS = [:tags, :browsers].freeze

    def process
      step_scratch = []
      step_settings_scratch = {}

      text.lines.each_with_index do |line, line_no|
        line = line.chomp
        if line[0..1] == '#!'
          @test.rfml_id = line[2..-1].strip.split(' ')[0]

        elsif line[0] == '#'
          comment = line[1..-1].strip

          if comment.start_with?('redirect:')
            value = comment.split(' ')[1..-1].join(' ').strip
            if %(true false).include?(value)
              step_settings_scratch[:redirect] = value
            else
              @test.errors[line_no] = Error.new(line_no, 'Redirection value must be true or false')
            end
          else
            special_fields = (CSV_FIELDS + TEST_DATA_FIELDS)
            matched_field = special_fields.find { |f| comment.start_with?("#{f}:") }
            if matched_field.nil?
              # comment, store in description
              @test.description += comment + "\n"
            else
              # extract just the text of the field
              @test[matched_field] = comment.split(' ')[1..-1].join(' ').strip

              # if it's supposed to be a CSV field, split and trim it
              if CSV_FIELDS.include?(matched_field)
                @test[matched_field] = @test[matched_field].split(',').map(&:strip).map(&:downcase)
              end
            end
          end

        elsif step_scratch.count == 0 && line.strip != ''
          if line[0] == '-'
            @test.steps << EmbeddedTest.new(line[1..-1].strip, step_settings_scratch[:redirect])
            step_settings_scratch = {}
          else
            step_scratch << line.strip
          end

        elsif step_scratch.count == 1
          if line.strip == ''
            @test.errors[line_no] = Error.new(line_no, 'Missing question')
          elsif !line.include?('?')
            @test.errors[line_no] = Error.new(line_no, 'Missing ?')
          else
            step_scratch << line.strip
          end

        elsif line.strip.empty? && step_settings_scratch.any?
          @test.errors[line_no] = Error.new(line_no, 'Extra space between step attributes and step content.')
        end

        if @test.errors.has_key?(line_no)
          step_scratch = []
        end

        if step_scratch.count == 2
          @test.steps << Step.new(step_scratch[0], step_scratch[1], step_settings_scratch[:redirect])
          step_scratch = []
          step_settings_scratch = {}
        end
      end

      if @test.rfml_id == nil
        @test.errors[0] = Error.new(0, 'Missing RFML ID. Please start a line #! followed by a unique id.')
      end

      return @test
    end
  end
end
