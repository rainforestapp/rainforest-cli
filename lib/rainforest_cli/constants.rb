# frozen_string_literal: true
module RainforestCli
  THREADS = ENV.fetch('RAINFOREST_THREADS', 16).to_i
end
