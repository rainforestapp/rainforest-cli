# frozen_string_literal: true
module RainforestCli
  module Notifications
    def parsing_error_notification!(rfml_tests)
      logger.error 'Parsing errors:'
      logger.error ''
      rfml_tests.each do |rfml_test|
        logger.error "\t#{rfml_test.file_name}"
        rfml_test.errors.each do |_line, error|
          logger.error "\t#{error}"
        end
      end

      exit 1
    end

    def nonexisting_embedded_id_notification!(rfml_tests)
      logger.error 'The following files contain unknown embedded test IDs:'
      logger.error ''
      rfml_tests.each do |rfml_test|
        logger.error "\t#{rfml_test.file_name}"
      end

      exit 2
    end

    def circular_dependencies_notification!(file_a, file_b)
      logger.error 'The following files are embedding one another:'
      logger.error ''
      logger.error "\t#{file_a}"
      logger.error "\t#{file_b}"

      exit 3
    end

    private

    def logger
      RainforestCli.logger
    end
  end
end
