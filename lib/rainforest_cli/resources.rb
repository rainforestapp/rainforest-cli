# frozen_string_literal: true
require 'terminal-table'

class RainforestCli::Resources
  HUMANIZE_SITE_CATEGORIES = {
    'device_farm' => 'Device Farm',
    'android' => 'Android',
    'ios' => 'iOS',
    'site' => 'Site',
  }.freeze

  def initialize(options)
    @client = RainforestCli::HttpClient.new(token: options.token)
  end

  def sites
    sites = @client.get('/sites').map do |s|
      category = s['category']
      [s['id'], s['name'], HUMANIZE_SITE_CATEGORIES[category] || category]
    end

    if sites.empty?
      logger.info('No sites found on your account.')
      logger.info('Please visit https://app.rainforestqa.com/settings/sites to create and edit your sites.')
    else
      print_table(['Site ID', 'Site Name', 'Category'], sites)
    end
  end

  def folders
    folders = @client.get('/folders?page_size=100').map { |f| f.values_at('id', 'title') }

    if folders.empty?
      logger.info('No folders found on your account.')
      logger.info('Please visit https://app.rainforestqa.com/folders to create and edit your sites.')
    else
      print_table(['Folder ID', 'Folder Name'], folders)
    end
  end

  def browsers
    account = @client.get('/clients')
    browsers = account['available_browsers'].map { |b| b.values_at('name', 'description') }
    print_table(['Browser ID', 'Browser Name'], browsers)
  end

  def print_table(headings, rows)
    table = Terminal::Table.new(headings: headings, rows: rows)
    table.align_column(0, :right)
    puts table
  end

  private

  def logger
    RainforestCli.logger
  end
end
