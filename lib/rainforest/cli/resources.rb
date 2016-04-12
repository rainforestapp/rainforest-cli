# frozen_string_literal: true
class RainforestCli::Resources
  def initialize(options)
    @client = RainforestCli::HttpClient.new(token: options.token)
  end

  def sites
    sites = @client.get('/sites')

    if sites.empty?
      logger.info('No sites found on your account.')
      logger.info('Please visit https://app.rainforestqa.com/settings/sites to create and edit your sites.')
    else
      print_table('Site', sites) do |site|
        { id: site['id'], name: site['name'] }
      end
    end
  end

  def folders
    folders = @client.get('/folders')

    if folders.empty?
      logger.info('No folders found on your account.')
      logger.info('Please visit https://app.rainforestqa.com/folders to create and edit your sites.')
    else
      print_table('Folder', folders) do |folder|
        { id: folder['id'], name: folder['title'] }
      end
    end
  end

  def print_table(resource_name, resources)
    table_heading = "#{resource_name} ID | #{resource_name} Name"
    puts table_heading
    puts '-' * table_heading.length
    resources.each do |resource|
      resource_data = yield(resource)
      puts "#{resource_data[:id].to_s.rjust(7)} | #{resource_data[:name]}"
    end
  end

  private

  def logger
    RainforestCli.logger
  end
end
