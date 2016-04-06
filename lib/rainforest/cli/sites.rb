# frozen_string_literal: true
class RainforestCli::Sites
  def initialize(options)
    @client = RainforestCli::HttpClient.new(token: options.token)
  end

  def list_sites
    sites = fetch_sites

    if sites.empty?
      logger.info('No configured sites found on your account.')
      logger.info('Please visit https://app.rainforestqa.com/settings/sites to create and edit your sites.')
    else
      print_site_table(sites)
    end
  end

  def print_site_table(sites)
    puts 'Site ID | Site Name'
    puts '-------------------'
    sites.each do |site|
      puts "#{site['id'].to_s.rjust(7)} | #{site['name']}"
    end
  end

  private

  def fetch_sites
    @sites ||= @client.get('/sites')
  end

  def logger
    RainforestCli.logger
  end
end
