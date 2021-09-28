require 'builder'
require 'httparty'
require 'json'
require 'fileutils'

res = HTTParty.get('https://api.github.com/repos/rainforestapp/rainforest-cli/releases')
raise StandardError, "Error #{res.code} while fetching releases:\n#{res.body}" unless res.code == 200

releases = JSON.parse(res.body)

latest_release = releases.find do |release|
  !release['draft'] && !release['prerelease']
end


class Release
  attr_reader :release

  def initialize(release)
    @release = release
  end

  def version
    @release['tag_name'][1..-1]
  end

  def windows_amd64_zip_name
    "rainforest-cli-#{version}-windows-amd64.zip"
  end

  def asset
    release['assets'].find do |asset|
      asset['name'] == windows_amd64_zip_name
    end
  end

  def download
    File.open(windows_amd64_zip_name, "w") do |file|
      response = HTTParty.get(asset['browser_download_url'], stream_body: true) do |fragment|
        if [301, 302].include?(fragment.code)
          puts "...skipped writing for redirect"
        elsif fragment.code == 200
          file.write(fragment)
        else
          raise StandardError, "Fail while streaming #{fragment.code}"
        end
      end
    end
  end
end

latest_release = Release.new(latest_release)

puts "Building Chocolatey package for #{latest_release.version}"
puts "\tFetching #{latest_release.windows_amd64_zip_name}"

# fetch
latest_release.download

# unzip, move
FileUtils.rm_rf('tmp')
FileUtils.rm_rf('rainforest-cli')
FileUtils.mkdir_p('tmp')
FileUtils.mkdir_p(File.join('rainforest-cli', 'tools'))
`unzip -n #{latest_release.windows_amd64_zip_name} -d tmp`
FileUtils.mv(File.join('tmp', 'rainforest-cli.exe'), File.join('rainforest-cli', 'tools'))
FileUtils.rm_rf('tmp')

# write the nuget
builder = Builder::XmlMarkup.new(indent: 2)
builder.instruct!(:xml, version: '1.0', encoding: 'UTF-8')

xml = builder.package(xmlns: 'http://schemas.microsoft.com/packaging/2015/06/nuspec.xsd') do |package|
  package.metadata do |metadata|
    metadata.title('rainforest-cli')
    metadata.id('rainforest-cli')
    metadata.version(latest_release.version)
    metadata.summary('A command line interface to interact with Rainforest QA - https://www.rainforestqa.com/.')
    metadata.tags('rainforest-cli rainforest')
    metadata.owners('@ukd1')

    metadata.packageSourceUrl('https://github.com/rainforestapp/rainforest-cli/tree/master/chocolatey')
    metadata.authors('https://github.com/rainforestapp/rainforest-cli/graphs/contributors')
    metadata.projectUrl('https://github.com/rainforestapp/rainforest-cli')
    metadata.iconUrl('https://assets.website-files.com/60da68c37e57671c365004bd/60da68c37e576749595005ae_favicon-large.svg')
    metadata.copyright('2021 Rainforest QA, Inc')

    metadata.licenseUrl('https://github.com/rainforestapp/rainforest-cli/blob/master/LICENSE.txt')
    metadata.requireLicenseAcceptance(true)
    metadata.projectSourceUrl('https://github.com/rainforestapp/rainforest-cli')
    metadata.docsUrl('https://github.com/rainforestapp/rainforest-cli/blob/master/README.md')
    metadata.bugTrackerUrl('https://github.com/rainforestapp/rainforest-cli/issues')
    metadata.description(File.read('../README.md').to_s)
    metadata.releaseNotes(File.read('../CHANGELOG.md').to_s)
  end

  package.files do |files|
    files.file(src: 'tools/rainforest-cli.exe', target: 'rainforest-cli.exe')
  end
end

File.write('rainforest-cli/rainforest-cli.nuspec', xml)