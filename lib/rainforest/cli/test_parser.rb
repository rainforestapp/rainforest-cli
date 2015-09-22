module RainforestCli::TestParser
  class EmbeddedTest < Struct.new(:test_name)
    def to_s
      "--> embed: #{test_name}"
    end
  end

  class Step < Struct.new(:action, :response)
    def to_s
      "#{action} --> #{response}"
    end
  end

  class Error < Struct.new(:line, :reason)
    def to_s
      "Line #{line}: #{reason}"
    end
  end

  class Test < Struct.new(:id, :description, :steps, :errors)
  end

  class Parser
    attr_reader :steps, :errors, :text

    def initialize(text)
      @text = text.to_s

      @test = Test.new
      @test.description = ""
      @test.steps = []
      @test.errors = {}
    end

    def process
      scratch = []

      text.lines.map(&:chomp).each_with_index do |line, line_no|
        if line[0..1] == '#!'
          # special comment, don't ignore!
          @test.id = line[2..-1].strip.split(" ")[0]

        elsif line[0] == '#'
          # comment, store in description
          @test.description += line[1..-1] + "\n"

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

      if @test.id == nil
        @test.errors[0] = Error.new(0, "Missing test ID. Please start a line #! followed by a unique id.")
      end

      return @test
    end
  end
end