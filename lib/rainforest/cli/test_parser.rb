# frozen_string_literal: true
module RainforestCli::TestParser
  class EmbeddedTest < Struct.new(:rfml_id, :redirect)
    def type
      :test
    end

    def to_s
      "--> embed: #{rfml_id}"
    end

    def redirection
      redirect || 'true'
    end

    def to_element(primary_key_id)
      {
        type: 'test',
        redirection: redirection,
        element: {
          id: primary_key_id
        }
      }
    end
  end

  class Step < Struct.new(:action, :response, :redirect)
    def type
      :step
    end

    def redirection
      redirect || 'true'
    end

    def to_s
      "#{action} --> #{response}"
    end

    def to_element
      {
        type: 'step',
        redirection: redirection,
        element: {
          action: action,
          response: response
        }
      }
    end
  end

  class Test < Struct.new(
    :file_name,
    :rfml_id,
    :description,
    :title,
    :start_uri,
    :site_id,
    :steps,
    :errors,
    :tags,
    :browsers
  )
    def embedded_ids
      steps.inject([]) { |embeds, step| step.type == :test ? embeds + [step.rfml_id] : embeds }
    end
  end

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
      @test.file_name = file_name
      @test.description = ''
      @test.steps = []
      @test.errors = {}
      @test.tags = []
      @test.browsers = []
    end

    TEST_DATA_FIELDS = [:start_uri, :title, :site_id, :browsers].freeze
    STEP_DATA_FIELDS = [:redirect].freeze
    CSV_FIELDS = [:tags, :browsers].freeze

    def process
      step_scratch = []
      step_settings_scratch = {}

      text.lines.each_with_index do |line, line_no|
        line = line.chomp
        if line[0..1] == '#!'
          # special comment, don't ignore!
          @test.rfml_id = line[2..-1].strip.split(' ')[0]
          @test.description += line[1..-1] + "\n"

        elsif line[0] == '#'
          # comment, store in description
          @test.description += line[1..-1] + "\n"

          if line[1..-1].strip[0..8] == 'redirect:'
            value = line[1..-1].split(' ')[1..-1].join(' ').strip
            if %(true false).include?(value)
              step_settings_scratch[:redirect] = value
            else
              @test.errors[line_no] = Error.new(line_no, 'Redirection value must be true or false')
            end
          end

          (CSV_FIELDS + TEST_DATA_FIELDS).each do |field|
            next unless line[1..-1].strip[0..(field.length)] == "#{field}:"

            # extract just the text of the field
            @test[field] = line[1..-1].split(' ')[1..-1].join(' ').strip

            # if it's supposed to be a CSV field, split and trim it
            if CSV_FIELDS.include?(field)
              @test[field] = @test[field].split(',').map(&:strip).map(&:downcase)
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
