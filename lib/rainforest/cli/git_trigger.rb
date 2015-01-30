require 'optparse'

module Rainforest
  module Cli
    class GitTrigger
      def self.git_trigger_should_run?(commit_message)
        commit_message.include?('@rainforest')
      end

      def self.extract_hashtags(commit_message)
        commit_message.scan(/#([\w_-]+)/).flatten.map {|s| s.gsub('#','') }
      end

      def self.last_commit_message
        `git log -1 --pretty=%B`.strip
      end
    end
  end
end
