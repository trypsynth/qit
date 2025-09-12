require "../command"
require "json"

module Qit::Commands
  class LicenseCommand < Command
    def name : String
      "license"
    end

    def description : String
      "download license template from GitHub or list available licenses"
    end

    def execute(args : Array(String)) : Nil
      Utils.require_args(args, "Missing license name or 'list'. Use 'license list' to see available licenses.")
      if args[0].downcase == "list"
        list_github_licenses
      else
        download_github_license args[0].downcase
      end
    end

    private def list_github_licenses
      url = "https://api.github.com/licenses"
      Utils.http_request(url) do |response|
        if response.status_code == 200
          licenses = Array(JSON::Any).from_json(response.body)
          puts "Available licenses:"
          licenses.each do |license|
            key = license["key"].as_s
            name = license["name"].as_s
            puts "#{key}: #{name}"
          end
        else
          STDERR.puts "Error fetching licenses. HTTP #{response.status_code}"
        end
      end
    end

    private def download_github_license(license_key : String)
      url = "https://api.github.com/licenses/#{license_key}"
      Utils.http_request(url) do |response|
        if response.status_code == 200
          license_data = JSON.parse(response.body)
          license_body = license_data["body"].as_s
          license_name = license_data["name"].as_s
          File.write("LICENSE", license_body)
          puts "Downloaded the #{license_name} license to LICENSE file."
        elsif response.status_code == 404
          STDERR.puts "License '#{license_key}' not found. Use 'qit license list' to see available licenses."
          exit 1
        else
          STDERR.puts "Error: HTTP #{response.status_code} from GitHub API"
          exit 1
        end
      end
    end
  end
end
