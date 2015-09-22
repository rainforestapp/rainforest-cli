require 'securerandom'

class RainforestCli::TestImporter
  attr_reader :options, :client
  SPEC_FOLDER = 'spec/rainforest'.freeze
  EXT = ".rfml".freeze

  SAMPLE_FILE = <<EOF
#! %s do not edit this line
#
# New test
#
# 1. steps:
#   a) pairs of lines are steps (first line = action, second = response)
#   b) second line must have a ?
#   c) second line must not be blank
# 2. embeds
#   a) lines starting with - are embedded tests
# 3. comments
#   a) lines starting # are comments
#

EOF

  def initialize(options)
    @options = options
    unless File.exists?(SPEC_FOLDER)
      logger.fatal "Rainforest folder not found (#{SPEC_FOLDER})" 
      exit 2
    end
  end

  def logger
    RainforestCli.logger
  end

  def upload
    ::Rainforest.api_key = @options.token
    test = Rainforest::Test.retrieve(123)
  end

  def validate
    tests = {}
    has_errors = []
    
    Dir.glob("#{SPEC_FOLDER}/**/*#{EXT}").each do |file_name|
      out = RainforestCli::TestParser::Parser.new(File.read(file_name)).process

      tests[file_name] = out
      has_errors << file_name if out.errors != {}
    end

    if !has_errors.empty?
      logger.error "Parsing errors:"
      logger.error ""
      has_errors.each do |file_name|
        logger.error " " + file_name
        tests[file_name].errors.each do |line, error|
          logger.error "\t#{error.to_s}"
        end
      end

      exit 2
    end

    if @options.debug
      tests.each do |file_name,test|
        logger.debug test.inspect
        logger.debug "#{file_name}"
        logger.debug test.description
        test.steps.each do |step|
          logger.debug "\t#{step}"
        end
      end
    else
      logger.info "[VALID]"
    end
  end

  def create_new
    name = @options.file_name if @options.file_name

    uuid = SecureRandom.uuid
    name = "#{uuid}#{EXT}" unless name
    name += EXT unless name[-EXT.length..-1] == EXT
    name = File.join([SPEC_FOLDER, name])

    File.open(name, "w") { |file| file.write(sprintf(SAMPLE_FILE, uuid)) }

    logger.info "Created #{name}"
  end
end
