module RainforestCli::TestParser
  class EmbeddedTest < Struct.new(:rfml_id)
    def type
      :test
    end

    def to_s
      "--> embed: #{rfml_id}"
    end
  end

  class Step < Struct.new(:action, :response)
    def type
      :step
    end

    def to_s
      "#{action} --> #{response}"
    end
  end

  class Error < Struct.new(:line, :reason)
    def to_s
      "Line #{line}: #{reason}"
    end
  end

  class Test < Struct.new(:rfml_id, :description, :title, :start_uri, :steps, :errors, :tags, :browsers)
    def embedded_ids
      steps.inject([]) { |embeds, step| step.type == :test ? embeds + [step.rfml_id] : embeds }
    end
  end

  class Parser
    attr_reader :steps, :errors, :text

    def initialize(text)
      @text = text.to_s

      @test = Test.new
      @test.description = ""
      @test.steps = []
      @test.errors = {}
      @test.tags = []
      @test.browsers = []
    end

    TEXT_FIELDS = [:start_uri, :title, :tags].freeze
    CSV_FIELDS = [:tags, :browsers].freeze

    def process
      scratch = []

      text.lines.map(&:chomp).each_with_index do |line, line_no|
        if line[0..1] == '#!'
          # special comment, don't ignore!
          @test.rfml_id = line[2..-1].strip.split(" ")[0]
          @test.description += line[1..-1] + "\n"

        elsif line[0] == '#'
          # comment, store in description
          @test.description += line[1..-1] + "\n"

          (CSV_FIELDS + TEXT_FIELDS).each do |field|
            next unless line[1..-1].strip[0..(field.length)] == "#{field}:"

            # extract just the text of the field
            @test[field] = line[1..-1].split(" ")[1..-1].join(" ").strip

            # if it's supposed to be a CSV field, split and trim it
            if CSV_FIELDS.include?(field)
              @test[field] = @test[field].split(',').map(&:strip).map(&:downcase)
            end
          end

        elsif scratch.count == 0 && line.strip != ''
          if line[0] == '-'
            @test.steps << EmbeddedTest.new(line[1..-1].strip)
          else
            scratch << line.strip
          end

        elsif scratch.count == 1
          if line.strip == ''
            @test.errors[line_no] = Error.new(line_no, "Missing question")
          elsif !line.include?('?')
            @test.errors[line_no] = Error.new(line_no, "Missing ?")
          else
            scratch << line.strip
          end

        end

        if @test.errors.has_key?(line_no)
          scratch = []
        end

        if scratch.count == 2
          @test.steps << Step.new(scratch[0], scratch[1])
          scratch = []
        end
      end

      if @test.rfml_id == nil
        @test.errors[0] = Error.new(0, "Missing test ID. Please start a line #! followed by a unique id.")
      end

      return @test
    end
  end
end
