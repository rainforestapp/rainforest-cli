# frozen_string_literal: true
class RainforestCli::Resources
  Resource = Struct.new(:identifier, :name)

  def initialize(options)
    @client = RainforestCli::HttpClient.new(token: options.token)
  end

  def sites
    sites = @client.get('/sites').map { |s| Resource.new(s['id'], s['name']) }

    if sites.empty?
      logger.info('No sites found on your account.')
      logger.info('Please visit https://app.rainforestqa.com/settings/sites to create and edit your sites.')
    else
      print_table('Site', sites)
    end
  end

  def folders
    folders = @client.get('/folders?page_size=100').map { |f| Resource.new(f['id'], f['title']) }
    if folders.empty?
      logger.info('No folders found on your account.')
      logger.info('Please visit https://app.rainforestqa.com/folders to create and edit your sites.')
    else
      print_table('Folder', folders)
    end
  end

  def browsers
    account = @client.get('/clients')
    browsers = account['available_browsers'].map { |b| Resource.new(b['name'], b['description']) }
    print_table('Browser', browsers)
  end

  def print_table(resource_name, resources)
    id_col = "#{resource_name} ID"
    longest_id = resources.map { |r| r.identifier.to_s }.max_by(&:length)
    col_length = [id_col.length, longest_id.length].max

    table_heading = "#{id_col.rjust(col_length)} | #{resource_name} Name"
    puts table_heading
    puts '-' * table_heading.length
    resources.each do |resource|
      puts "#{resource.identifier.to_s.rjust(col_length)} | #{resource.name}"
    end
  end

  private

  def logger
    RainforestCli.logger
  end
end
