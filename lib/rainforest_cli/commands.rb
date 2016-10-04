# frozen_string_literal: true

module RainforestCli
  class Commands
    Command = Struct.new(:name, :description, :block)

    attr_reader :commands

    def initialize
      @commands = []
      yield(self) if block_given?
    end

    def add(command, description, &blk)
      @commands << Command.new(command, description, blk)
    end

    def call(command_name)
      command = @commands.find { |c| c.name == command_name }

      if command.nil?
        logger.fatal "Unknown command: #{command_name}"
        exit 2
      end

      command.block.call
    end

    def print_documentation
      command_col_width = @commands.map { |c| c.name.length }.max
      puts 'Rainforest CLI commands:'
      @commands.each do |command|
        puts "\t#{command.name.ljust(command_col_width)}\t\t#{command.description}"
      end
    end

    private

    def logger
      RainforestCli.logger
    end
  end
end
