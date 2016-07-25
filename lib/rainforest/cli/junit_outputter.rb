require 'builder'
require 'time'
require 'json'

# frozen_string_literal: true
module RainforestCli
  class JunitOutputter
    attr_reader :builder, :client

    def initialize(token, run, tests)
      @client = HttpClient.new token: token
      @json_run = run # JSON containing the results of /1/runs/{run_id}.json
      @json_tests = tests # JSON containing the results of /1/runs/{run_id}/tests.json
      @builder = Builder::XmlMarkup.new( :indent => 2)
    end # end initialize

    def build_test_suite
      @json_tests.each do | test |
        build_test test
      end # end do
    end # end process_run_results

    def build_test(test)
      test_name = test['title']
      execution_time = Time.parse(test['updated_at']) - Time.parse(test['created_at'])
      test_status = test['result']
      @builder.testcase(:name => test_name, :time => execution_time) do
        case test_status
        when "failed"
          build_failed_test test
        end # end case
      end # end do
    end # end build_test

    def build_failed_test(test)
      response = client.get("/runs/#{@json_run['id']}/tests/#{test['id']}.json")
      response['steps'].each do | step |
        step['browsers'].each do | browser |
          browser_name = browser['name']
          browser['feedback'].each do | opinion |
            if opinion['answer_given'] == 'no' and opinion['job_state'] == 'approved'
              if opinion['note'] != ""
                @builder.failure(:type => browser_name, :message => opinion['note'])
              end
            end
          end
        end
      end
    end # end build_failed_test

    def parse
      @builder.instruct! :xml, :version => "1.0", :encoding => "UTF-8"
      @builder.testsuite(
        :name => @json_run['description'],
        :errors => @json_run['total_no_result_tests'],
        :failures => @json_run['total_failed_tests'],
        :tests => @json_run['total_tests'],
        :time => Time.parse(@json_run['timestamps']['complete']) - Time.parse(@json_run['timestamps']['created_at']),
        :timestamp => @json_run['created_at']) do
          build_test_suite
        end
    end # end parse

    def output(stream)
      stream.write(@builder.target!)
    end #end output

  end # end class
end # end module
