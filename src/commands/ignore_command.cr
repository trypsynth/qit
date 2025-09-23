require "../command"

module Qit::Commands
  class IgnoreCommand < Command
    def name : String
      "ignore"
    end

    def description : String
      "download .gitignore template(s) from gitignore.io or list available templates"
    end

    def usage : String
      "ignore <template_name|list>"
    end

    def execute(args : Array(String)) : Nil
      Utils.require_args(args, "Missing template name(s). Use 'ignore list' to see available templates.")
      if args[0].downcase == "list"
        list_gitignore_templates
      else
        download_gitignore args[0]
      end
    end

    private def download_gitignore(templates : String)
      url = "https://www.toptal.com/developers/gitignore/api/#{templates}"
      Utils.http_request(url) do |response|
        if response.status_code == 200
          File.write(".gitignore", response.body)
          puts "Downloaded .gitignore for #{templates}"
        else
          STDERR.puts "Error: HTTP #{response.status_code} from gitignore.io"
          exit 1
        end
      end
    end

    private def list_gitignore_templates
      url = "https://www.toptal.com/developers/gitignore/api/list?format=lines"
      Utils.http_request(url) do |response|
        if response.status_code == 200
          puts "Available gitignore templates:"
          puts response.body
        else
          STDERR.puts "Error fetching list. HTTP #{response.status_code}"
        end
      end
    end
  end
end
