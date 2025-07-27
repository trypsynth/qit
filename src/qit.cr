require "http/client"

def usage
  puts <<-USAGE
  Qit - Quin's tiny Git helper.
  Usage: qit <command> [<args>...]
  Available commands:
    acp <message>: add all files, commit with the specified message, and push.
    amend <message>: amend the last commit with a new message.
    ignore <Template>: download a .gitignore template from GitHub.
    last [<number>]: show the last <number> commits (default: 1).
    log: show the commit log in a readable format.
    reset: hard reset to the last commit, discarding all local changes.
    undo: undo the last commit while keeping changes intact.
  USAGE
end

def git(*args : String)
  abort "Command 'git #{args.join(" ")}' failed." unless Process.run("git", args, output: STDOUT, error: STDERR).success?
end

def get_commit_message(prefix : String = "") : String
  if ARGV.size < 2
    desc = prefix.empty? ? "commit message" : "#{prefix} commit message"
    STDERR.puts "Missing#{desc}."
    usage
    exit 1
  end
  ARGV[1..].join(" ")
end

def download_gitignore_template(template : String)
  url = "https://raw.githubusercontent.com/github/gitignore/main/#{template}.gitignore"
  puts "Downloading .gitignore template from #{url}"
  begin
    response = HTTP::Client.get(url)
    File.write(".gitignore", response.body)
    puts "Downloaded and saved as .gitignore"
  rescue ex
    STDERR.puts "Failed to download template '#{template}'. Check if the name is correct and case-sensitive."
    STDERR.puts ex.message
    exit 1
  end
end

if ARGV.empty?
  usage
  exit 1
end
command = ARGV[0].downcase
case command
when "acp"
  message = get_commit_message
  git "add", "."
  git "commit", "-m", message
  git "push"
when "amend"
  message = get_commit_message "new"
  git "commit", "--amend", "--reset", "-m", message
when "ignore"
  template = ARGV[1]?
  unless template
    STDERR.puts "Missing template name."
    usage
    exit 1
  end
  download_gitignore_template template
when "last"
  count = ARGV[1]?.try(&.to_i?) || 1
  git "log", "-#{count}", "--pretty=format:%h %an: %s (%ad).", "--date=format:%Y-%m-%d %H:%M:%S"
when "log"
  git "log", "--pretty=format:%h %an: %s (%ad).", "--date=format:%Y-%m-%d %H:%M:%S"
when "reset"
  git "reset", "--hard"
when "undo"
  git "reset", "--soft", "HEAD~1"
else
  STDERR.puts "Unknown command: #{command}."
  usage
  exit 1
end
